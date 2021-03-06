/*
 * Copyright (C) 2017-2018 GIG Technology NV and Contributors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package zerodb

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"

	"github.com/threefoldtech/0-stor/client/datastor"
)

var (
	readTimeout    = 60 * time.Second
	writeTimeout   = 60 * time.Second
	connectTimeout = 3 * time.Second
)

// Client defines a data client,
// to connect to a 0-db server
type Client struct {
	pool          *redis.Pool
	namespace     string
	utilization   int64
	muUtilization sync.Mutex
}

// NewClient creates a new data client,
// with given server address & password,
// and use the given namespace.
func NewClient(addr, passwd, namespace string) (*Client, error) {
	var opts = []redis.DialOption{
		redis.DialReadTimeout(readTimeout),
		redis.DialWriteTimeout(writeTimeout),
		redis.DialConnectTimeout(connectTimeout),
	}

	if len(addr) == 0 {
		return nil, fmt.Errorf("no address given")
	}
	if len(namespace) == 0 {
		return nil, fmt.Errorf("no namespace given")
	}

	selectArgs := []interface{}{namespace}
	if passwd != "" {
		selectArgs = append(selectArgs, passwd)
	}

	// creates pool
	pool := &redis.Pool{
		Wait:      true,
		MaxActive: 5,
		MaxIdle:   5,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", addr, opts...)
			if err != nil {
				return nil, fmt.Errorf("failed to dial 0-db: %v", err)
			}

			_, err = conn.Do("SELECT", selectArgs...)
			if err != nil {
				return nil, fmt.Errorf("failed to select %s: %v", selectArgs[0], err)
			}
			return conn, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}

	client := &Client{
		pool:      pool,
		namespace: namespace,
	}
	// save the current utilization of the namespace in the client object
	// this is then used during the lifetime of the client
	// to allow different sorting algorithms in ShardIterator
	ns, err := client.GetNamespace()
	if err != nil {
		return nil, err
	}
	client.utilization = ns.Used
	return client, nil
}

// CreateObject implements datastor.Client.CreateObject
func (c *Client) CreateObject(data []byte) (key []byte, err error) {
	conn := c.pool.Get()
	defer conn.Close()

	key, err = redis.Bytes(conn.Do("SET", dummyKey, data))
	if err != nil {
		if err.Error() == "No space left on this namespace" {
			err = datastor.ErrNamespaceFull
		}
		return key, err
	}
	c.muUtilization.Lock()
	c.utilization += int64(len(data))
	c.muUtilization.Unlock()
	return key, nil
}

// GetObject implements datastor.Client.GetObject
func (c *Client) GetObject(key []byte) (*datastor.Object, error) {
	conn := c.pool.Get()
	defer conn.Close()

	data, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		if err == redis.ErrNil {
			return nil, datastor.ErrKeyNotFound
		}
		return nil, err
	}

	return &datastor.Object{
		Key:  key,
		Data: data,
	}, nil
}

// DeleteObject implements datastor.Client.DeleteObject
func (c *Client) DeleteObject(key []byte) error {
	conn := c.pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	return err
}

// GetObjectStatus implements datastor.Client.GetObjectStatus
func (c *Client) GetObjectStatus(key []byte) (datastor.ObjectStatus, error) {
	conn := c.pool.Get()
	defer conn.Close()

	status, err := redis.Int(conn.Do("CHECK", key))
	if err != nil {
		return 0, err
	}
	switch status {
	case 1:
		return datastor.ObjectStatusOK, nil
	case 0:
		return datastor.ObjectStatusCorrupted, nil
	default:
		return datastor.ObjectStatusMissing, nil
	}
}

// ExistObject implements datastor.Client.ExistObject
func (c *Client) ExistObject(key []byte) (bool, error) {
	conn := c.pool.Get()
	defer conn.Close()

	exist, err := redis.Int(conn.Do("EXIST", key))
	if err != nil {
		return false, err
	}
	return exist == 1, nil
}

// ListObjectKeyIterator implements datastor.Client.ListObjectKeyIterator
// TODO: depend on 	https://github.com/threefoldtech/0-db/issues/6
func (c *Client) ListObjectKeyIterator(ctx context.Context) (<-chan datastor.ObjectKeyResult, error) {
	return nil, fmt.Errorf("not supported")
}

// GetNamespace implements datastor.Client.GetNamespace
func (c *Client) GetNamespace() (*datastor.Namespace, error) {
	conn := c.pool.Get()
	defer conn.Close()

	infoStr, err := redis.String(conn.Do("NSINFO", c.namespace))
	if err != nil {
		return nil, err
	}

	var ns datastor.Namespace
	var health datastor.Health
	ns.Health = &health

	for _, info := range strings.Split(infoStr, "\n") {
		elems := strings.Split(info, ":")
		if len(elems) != 2 {
			continue
		}
		val := strings.TrimSpace(elems[1])
		switch strings.TrimSpace(elems[0]) {
		case "name":
			ns.Label = val
		case "entries":
			entries, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return nil, err
			}
			ns.NrObjects = entries
		case "data_size_bytes":
			used, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return nil, err
			}
			ns.Used = used
		case "data_disk_freespace_bytes":
			free, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return nil, err
			}
			ns.Free = free
		case "stats_index_io_errors":
			count, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return nil, err
			}
			health.IndexIOErrors = count
		case "stats_index_faults":
			count, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return nil, err
			}
			health.IndexFaults = count

		case "stats_data_io_errors":
			count, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return nil, err
			}
			health.DataIOErrors = count
		case "stats_data_faults":
			count, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return nil, err
			}
			health.DataFaults = count

		case "stats_index_io_error_last":
			timestamp, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return nil, err
			}
			health.IndexIOErrorLast = timestamp
		case "stats_data_io_error_last":
			timestamp, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return nil, err
			}
			health.DataIOErrorLast = timestamp
		}
	}

	return &ns, nil
}

// Utilization returns the amount of storage used by the namespace
func (c *Client) Utilization() int64 {
	return c.utilization
}

// Close implements datastor.Client.Close
func (c *Client) Close() error {
	return c.pool.Close()
}

const (
	dummyKey = ""
)

var (
	_ datastor.Client = (*Client)(nil)
)
