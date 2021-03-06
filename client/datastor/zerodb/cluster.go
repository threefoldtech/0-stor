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
	"crypto/tls"
	"errors"
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/threefoldtech/0-stor/client/datastor"
)

// Cluster implements datastor.Cluster for
// clients which interface with 0-db server using redis protocol
type Cluster struct {
	namespace     string
	listedShards  map[string]*Shard // shards listed in config
	listedSlice   []*Shard
	unlistedMux   sync.Mutex
	passwd        string
	spreadingType datastor.SpreadingType
}

// NewCluster creates a new cluster,
// and pre-loading it with a client for each of the listed (and thus known) shards.
// Unlisted shards's clients are also stored, bu those are loaded on the fly, only when needed.
func NewCluster(addresses []datastor.ShardConfig, passwd, namespace string, tlsConfig *tls.Config, spreadingType datastor.SpreadingType) (*Cluster, error) {
	var (
		listedShards = make(map[string]*Shard, len(addresses))
		listedSlice  []*Shard
	)

	or := func(a, b string) string {
		if len(a) == 0 {
			return b
		}

		return a
	}

	for _, cfg := range addresses {

		client, err := NewClient(cfg.Address, or(cfg.Password, passwd), or(cfg.Namespace, namespace))
		if err != nil {
			return nil, err
		}
		shard := &Shard{
			Client:    client,
			address:   cfg.Address,
			namespace: or(cfg.Namespace, namespace),
			password:  or(cfg.Password, passwd),
		}
		listedShards[shard.Identifier()] = shard
		listedSlice = append(listedSlice, shard)
	}
	return &Cluster{
		namespace:     namespace,
		listedShards:  listedShards,
		listedSlice:   listedSlice,
		passwd:        passwd,
		spreadingType: spreadingType,
	}, nil
}

// GetShard implements datastor.Cluster.GetShard
func (c *Cluster) GetShard(address string) (datastor.Shard, error) {
	shard, ok := c.listedShards[address]
	if ok {
		return shard, nil
	}

	return nil, fmt.Errorf("shard %s not found", address)
}

// GetRandomShard implements datastor.Cluster.GetRandomShard
func (c *Cluster) GetRandomShard() (datastor.Shard, error) {
	index := datastor.RandShardIndex(int64(len(c.listedSlice)))
	return c.listedSlice[index], nil
}

// GetShardIterator implements datastor.Cluster.GetShardIterator
func (c *Cluster) GetShardIterator(exceptShards []string) datastor.ShardIterator {
	filtered := c.filteredSlice(exceptShards)

	switch c.spreadingType {
	case datastor.SpreadingTypeRandom:
		return datastor.NewRandomShardIterator(filtered)
	case datastor.SpreadingTypeLeastUsed:
		return datastor.NewLeastUsedShardIterator(filtered)
	default:
		panic("unsupported spreading algorithm")
	}
}

// ListedShardCount implements datastor.Cluster.ListedShardCount
func (c *Cluster) ListedShardCount() int {
	return len(c.listedSlice)
}

// Close implements datastor.Cluster.Close
func (c *Cluster) Close() error {
	c.unlistedMux.Lock()
	defer c.unlistedMux.Unlock()

	var (
		err      error
		errCount int
	)

	// close all listed shards next
	for address, shard := range c.listedShards {
		err = shard.Close()
		if err != nil {
			errCount++
			log.Errorf(
				"error while closing listed shard (%s): %v", address, err)
		}
	}

	// if at least one shard returned an error while closing,
	// return it as a generic error for now
	if errCount > 0 {
		return errors.New("one or multiple shards returned an error while closing")
	}
	return nil

}

func (c *Cluster) filteredSlice(exceptShards []string) []datastor.Shard {
	var (
		exceptMap = make(map[string]struct{}, len(exceptShards))
		filtered  = make([]datastor.Shard, 0, len(c.listedSlice))
	)

	for _, shard := range exceptShards {
		exceptMap[shard] = struct{}{}
	}

	for _, shard := range c.listedSlice {
		if _, ok := exceptMap[shard.Identifier()]; !ok {
			filtered = append(filtered, shard)
		}
	}
	return filtered
}

var (
	_ datastor.Cluster = (*Cluster)(nil)
)
