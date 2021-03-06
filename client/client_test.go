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

package client

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/threefoldtech/0-stor/client/datastor"
	"github.com/threefoldtech/0-stor/client/datastor/pipeline"
	"github.com/threefoldtech/0-stor/client/datastor/pipeline/storage"
	"github.com/threefoldtech/0-stor/client/metastor"
	"github.com/threefoldtech/0-stor/client/metastor/metatypes"
	"github.com/threefoldtech/0-stor/client/processing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClientFromConfigErrors(t *testing.T) {
	require := require.New(t)

	_, err := NewClientFromConfig(Config{}, nil, -1)
	require.Error(err, "missing: data shards and namespace")

	_, err = NewClientFromConfig(Config{Namespace: "foo"}, nil, -1)
	require.Error(err, "missing: data shards")

	// hard to test metastor creation, as it would require an etcd connection for now
	// TODO: once we have alternatives meta clients (e.g. badger), complete this test
	//       see: https://github.com/threefoldtech/0-stor/issues/419
}

func TestNewClientPanics(t *testing.T) {
	require := require.New(t)

	require.Panics(func() {
		NewClient(nil, nil)
	}, "nothing given")
	require.Panics(func() {
		NewClient(new(metastor.Client), nil)
	}, "no data pipeline given")
}

func TestRoundTrip(t *testing.T) {
	servers, serverClean := testZdbServer(t, 4)
	defer serverClean()

	shards := make([]datastor.ShardConfig, len(servers))
	for i, server := range servers {
		shards[i] = datastor.ShardConfig{Address: server.Address()}
	}

	config := Config{
		Namespace: "namespace1",
		DataStor:  DataStorConfig{Shards: shards},
	}

	const blockSize = 64

	tt := []struct {
		name string

		BlockSize int

		DataShards   int
		ParityShards int

		Compress   bool
		EncryptKey string
	}{
		{
			name: "simple-write",
		},
		{
			name:      "chunks",
			BlockSize: blockSize,
		},
		{
			name:     "compress",
			Compress: true,
		},
		{
			name:       "encrypt",
			EncryptKey: "cF0BFpIsljOS8UmaP8YRHRX0nBPVRVPw",
		},
		{
			name:       "replication",
			DataShards: len(shards),
		},
		{
			name:         "distribution",
			DataShards:   2,
			ParityShards: 1,
		},
		{
			name:         "chunks-distribution",
			BlockSize:    blockSize,
			DataShards:   2,
			ParityShards: 1,
		},
		{
			name:       "compress-encrypt",
			Compress:   true,
			EncryptKey: "cF0BFpIsljOS8UmaP8YRHRX0nBPVRVPw",
		},
		{
			name:       "chunk-compress-encrypt",
			BlockSize:  blockSize,
			Compress:   true,
			EncryptKey: "cF0BFpIsljOS8UmaP8YRHRX0nBPVRVPw",
		},
		{
			name:       "compress-encrypt-replication",
			Compress:   true,
			EncryptKey: "cF0BFpIsljOS8UmaP8YRHRX0nBPVRVPw",
			DataShards: len(shards),
		},
		{
			name:         "compress-encrypt-distribution",
			Compress:     true,
			EncryptKey:   "cF0BFpIsljOS8UmaP8YRHRX0nBPVRVPw",
			DataShards:   2,
			ParityShards: 1,
		},
		{
			name:       "chunks-compress-encrypt-replication",
			BlockSize:  blockSize,
			Compress:   true,
			EncryptKey: "cF0BFpIsljOS8UmaP8YRHRX0nBPVRVPw",
			DataShards: len(shards),
		},
		{
			name:         "chunks-compress-encrypt-distribution",
			BlockSize:    blockSize,
			Compress:     true,
			EncryptKey:   "cF0BFpIsljOS8UmaP8YRHRX0nBPVRVPw",
			DataShards:   2,
			ParityShards: 1,
		},
	}

	for _, tc := range tt {

		t.Run(tc.name, func(t *testing.T) {
			config.DataStor.Pipeline.BlockSize = tc.BlockSize
			if tc.Compress {
				config.DataStor.Pipeline.Compression.Mode = processing.CompressionModeDefault
			} else {
				config.DataStor.Pipeline.Compression.Mode = processing.CompressionModeDisabled
			}
			config.DataStor.Pipeline.Encryption.PrivateKey = tc.EncryptKey
			config.DataStor.Pipeline.Distribution.DataShardCount = tc.DataShards
			config.DataStor.Pipeline.Distribution.ParityShardCount = tc.ParityShards

			c, _, err := getTestClient(config)
			require.NoError(t, err, "fail to create client")

			data := make([]byte, blockSize*4)
			_, err = rand.Read(data)
			require.NoError(t, err, "fail to read random data")

			// write data to the store
			key := []byte("testkey")
			md, err := c.Write(key, bytes.NewReader(data))
			require.NoError(t, err, "fail to write data to the store")
			require.EqualValues(t, len(data), md.Size, "metadata.Size should be egals to the total size of the written file")

			// b, err := json.Marshal(meta)
			// require.NoError(t, err)
			// fmt.Println(string(b))

			// read data back
			dataReadBuf := bytes.NewBuffer(nil)
			meta, err := c.metastorClient.GetMetadata(key)
			require.NoError(t, err)

			err = c.Read(*meta, dataReadBuf)
			require.NoError(t, err, "fail to read data from the store")
			dataRead := dataReadBuf.Bytes()
			if bytes.Compare(data, dataRead) != 0 {
				t.Errorf("data read from store is not the same as original data")
				t.Errorf("len original: %d len actual %d", len(data), len(dataRead))
			}

			//delete data
			err = c.Delete(*meta)
			require.NoError(t, err, "failed to delete from the store")

			// makes sure metadata does not exist anymore
			_, err = c.metastorClient.GetMetadata(key)
			require.Error(t, err)
		})
	}
}

func TestWriteWithUserMeta(t *testing.T) {
	servers, serverClean := testZdbServer(t, 4)
	defer serverClean()

	shards := make([]datastor.ShardConfig, len(servers))
	for i, server := range servers {
		shards[i] = datastor.ShardConfig{Address: server.Address()}
	}

	const (
		blockSize = 256
	)

	config := newDefaultConfig(shards, blockSize)

	c, _, err := getTestClient(config)
	require.NoError(t, err, "fail to create client")
	defer c.Close()

	// initialize the data
	data := make([]byte, blockSize*10)
	_, err = rand.Read(data)
	require.NoError(t, err, "fail to read random data")
	key := []byte("testkey")

	testCases := []struct {
		name     string
		userMeta map[string]string
	}{
		{
			name:     "nil user meta",
			userMeta: nil,
		},
		{
			name: "non nil usermeta",
			userMeta: map[string]string{
				"key1": "val1",
				"key2": "val2",
				"key3": "val3",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err = c.WriteWithUserMeta(key, bytes.NewReader(data), tc.userMeta)
			require.NoError(t, err, "fail write data")

			// check the metadata
			md, err := c.metastorClient.GetMetadata(key)
			require.NoError(t, err)
			require.Equal(t, tc.userMeta, md.UserDefined)
		})
	}

}

func TestBlocksizes(t *testing.T) {
	servers, serverClean := testZdbServer(t, 4)
	defer serverClean()

	shards := make([]datastor.ShardConfig, len(servers))
	for i, server := range servers {
		shards[i] = datastor.ShardConfig{Address: server.Address()}
	}

	const baseBlockSize = 8

	config := newDefaultConfig(shards, 0)

	for i := 0; i <= 5; i++ {
		var blockSize int
		if i == 0 {
			blockSize = baseBlockSize * 10
		} else {
			blockSize = baseBlockSize * 10 * (i * 10)
		}

		t.Run(fmt.Sprint(blockSize), func(t *testing.T) {
			config.DataStor.Pipeline.BlockSize = blockSize
			c, _, err := getTestClient(config)
			require.NoError(t, err, "fail to create client")

			data := make([]byte, blockSize)
			_, err = rand.Read(data)
			require.NoError(t, err, "fail to read random data")

			// write data to the store
			key := []byte(fmt.Sprintf("testkey-%d", i))
			md, err := c.Write(key, bytes.NewReader(data))
			require.NoError(t, err, "fail to write data to the store")

			// read data back
			dataReadBuf := bytes.NewBuffer(nil)
			err = c.Read(*md, dataReadBuf)
			require.NoError(t, err, "fail to read data from the store")
			dataRead := dataReadBuf.Bytes()
			if bytes.Compare(data, dataRead) != 0 {
				t.Errorf("data read from store is not the same as original data")
				t.Errorf("len original: %d len actual %d", len(data), len(dataRead))
			}
		})
	}
}

func TestReadRangeChunked(t *testing.T) {
	testReadRange(t, 256)
}

func TestReadRangeNotChunked(t *testing.T) {
	testReadRange(t, 0)
}

func testReadRange(t *testing.T, chunkSize int) {
	servers, serverClean := testZdbServer(t, 4)
	defer serverClean()

	shards := make([]datastor.ShardConfig, len(servers))
	for i, server := range servers {
		shards[i] = datastor.ShardConfig{Address: server.Address()}
	}

	const (
		blockSize = 256
		numBlocks = 10
	)

	config := newDefaultConfig(shards, chunkSize)

	c, _, err := getTestClient(config)
	require.NoError(t, err, "fail to create client")
	defer c.Close()

	require.Equal(t, chunkSize, c.dataPipeline.ChunkSize())

	data := make([]byte, blockSize*numBlocks)

	_, err = rand.Read(data)
	require.NoError(t, err, "fail to read random data")
	key := []byte("testkey")

	_, err = c.Write(key, bytes.NewReader(data))
	require.NoError(t, err, "fail write data")

	testCases := []struct {
		name   string
		offset int
		length int
	}{
		{
			name:   "min offset max length",
			offset: 0,
			length: len(data),
		},
		{
			name:   "multiply of blocksize",
			offset: blockSize,
			length: blockSize * 3,
		},
		{
			name:   "length not multiply of blocksize",
			offset: blockSize,
			length: (blockSize * 3) - 1,
		},
		{
			name:   "offset not multiply of blocksize",
			offset: blockSize + 1,
			length: blockSize * 3,
		},
		{
			name:   "offset & length not multiply of blocksize",
			offset: blockSize + 1,
			length: (blockSize * 3) - 1,
		},
		{
			name:   "offset & length not multiply of blocksize, in same chunk",
			offset: blockSize + 1,
			length: blockSize - 10,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)

			meta, err := c.metastorClient.GetMetadata(key)
			require.NoError(t, err)

			err = c.ReadRange(*meta, buf, int64(tc.offset), int64(tc.length))

			require.NoError(t, err)
			require.Equal(t, tc.length, buf.Len())
			require.Equal(t, data[tc.offset:tc.offset+tc.length], buf.Bytes())
		})
	}
}

func TestMultipleDownload_Issue208(t *testing.T) {
	// #test for https://github.com/threefoldtech/0-stor/issues/208

	servers, serverClean := testZdbServer(t, 4)
	defer serverClean()

	shards := make([]datastor.ShardConfig, len(servers))
	for i, server := range servers {
		shards[i] = datastor.ShardConfig{Address: server.Address()}
	}

	const blockSize = 256

	config := newDefaultConfig(shards, blockSize)

	c, _, err := getTestClient(config)
	require.NoError(t, err, "fail to create client")
	defer c.Close()

	data := make([]byte, blockSize/16)

	_, err = rand.Read(data)
	require.NoError(t, err, "fail to read random data")
	key := []byte("testkey")

	md, err := c.Write(key, bytes.NewReader(data))
	require.NoError(t, err, "fail write data")

	buf := bytes.NewBuffer(nil)
	for i := 0; i < 100; i++ {
		err = c.Read(*md, buf)
		require.NoError(t, err, "fail read data")
		result := buf.Bytes()
		require.Equal(t, data, result)
		buf.Reset()
	}
}

func TestConcurrentWriteRead(t *testing.T) {
	servers, serverClean := testZdbServer(t, 4)
	defer serverClean()

	shards := make([]datastor.ShardConfig, len(servers))
	for i, server := range servers {
		shards[i] = datastor.ShardConfig{Address: server.Address()}
	}
	const blockSize = 128

	config := newDefaultConfig(shards, blockSize)

	doReadWrite := func(i, size int) {
		c, _, err := getTestClient(config)
		require.NoError(t, err, "fail to create client")
		defer c.Close()

		data := make([]byte, size)
		_, err = rand.Read(data)
		require.NoError(t, err, "fail to read random data")
		key := []byte(fmt.Sprintf("testkey-%d", i))

		md, err := c.Write(key, bytes.NewReader(data))
		require.NoError(t, err, "fail write data")
		require.EqualValues(t, len(data), md.Size, "metadata.Size should be egals to total size of the written file")

		buf := bytes.NewBuffer(nil)
		err = c.Read(*md, buf)
		require.NoError(t, err, "fail read data")
		result := buf.Bytes()
		require.Equal(t, data, result, "data read is not same as data written")
	}

	// Seems we can't increased the number of concurrent write more then around 32
	for concurrent := 1; concurrent <= 64; concurrent *= 2 {
		for size := blockSize; size < blockSize*10; size *= 4 {
			name := fmt.Sprintf("Concurrent client: %d - Size of the data: %d", concurrent, size)
			t.Log(name)

			wg := &sync.WaitGroup{}
			wg.Add(10)
			start := time.Now()
			for i := 0; i < 10; i++ {
				go func(i int) {
					defer wg.Done()
					doReadWrite(i, size)
				}(i)
			}
			wg.Wait()
			end := time.Now()
			t.Logf("duration %d ms\n\n", (end.Sub(start).Nanoseconds() / 1000000))
		}
	}
}

func BenchmarkWriteFilesSizes(b *testing.B) {
	servers, serverClean := testZdbServer(b, 4)
	defer serverClean()

	shards := make([]datastor.ShardConfig, len(servers))
	for i, server := range servers {
		shards[i] = datastor.ShardConfig{Address: server.Address()}
	}

	config := newDefaultConfig(shards, 1024*1024)

	c, _, err := getTestClient(config)
	require.NoError(b, err, "fail to create client")
	defer c.Close()

	tt := []struct {
		Size int
	}{
		{1024},              // 1k
		{1024 * 4},          // 4k
		{1024 * 10},         // 10k
		{1024 * 1024},       // 1M
		{1024 * 1024 * 10},  // 10M
		{1024 * 1024 * 100}, // 100M
		{1024 * 1024 * 500}, // 500M
	}

	for _, tc := range tt {

		b.Run(fmt.Sprintf("%d", tc.Size), func(b *testing.B) {

			data := make([]byte, tc.Size)
			_, err = rand.Read(data)
			require.NoError(b, err, "fail to read random data")

			key := []byte("testkey")

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				// write data to the store
				_, err = c.Write(key, bytes.NewReader(data))
				require.NoError(b, err, "fail to write data to the store")
			}
		})
	}
}

func TestIssue225(t *testing.T) {
	servers, serverClean := testZdbServer(t, 4)
	defer serverClean()

	shards := make([]datastor.ShardConfig, len(servers))
	for i, server := range servers {
		shards[i] = datastor.ShardConfig{Address: server.Address()}
	}
	const blockSize = 256

	config := newDefaultConfig(shards, blockSize)

	c, _, err := getTestClient(config)
	require.NoError(t, err, "fail to create client")
	defer c.Close()

	data := make([]byte, blockSize*11)

	_, err = rand.Read(data)
	require.NoError(t, err, "fail to read random data")
	key := []byte("testkey")

	meta, err := c.Write(key, bytes.NewReader(data))
	require.NoError(t, err, "fail write data")

	buf := bytes.NewBuffer(nil)
	err = c.Read(*meta, buf)
	require.NoError(t, err, "fail read data")
	result := buf.Bytes()
	assert.Equal(t, data, result)
}

func newDefaultConfig(dataShards []datastor.ShardConfig, blockSize int) Config {
	return Config{
		Namespace: "namespace1",
		DataStor: DataStorConfig{
			Shards: dataShards,
			Pipeline: pipeline.Config{
				BlockSize: blockSize,
				Compression: pipeline.CompressionConfig{
					Mode: processing.CompressionModeDefault,
				},
				Encryption: pipeline.EncryptionConfig{
					PrivateKey: "cF0BFpIsljOS8UmaP8YRHRX0nBPVRVPw",
				},
				Distribution: pipeline.ObjectDistributionConfig{
					DataShardCount:   2,
					ParityShardCount: 1,
				},
			},
		},
	}
}
func TestClientCheck(t *testing.T) {
	servers, serverClean := testZdbServer(t, 4)
	defer serverClean()

	shards := make([]datastor.ShardConfig, len(servers))
	for i, server := range servers {
		shards[i] = datastor.ShardConfig{Address: server.Address()}
	}

	config := newDefaultConfig(shards, 1024)

	c, datastorCluster, err := getTestClient(config)
	require.NoError(t, err, "fail to create client")
	defer c.Close()

	data := make([]byte, 602*10)

	_, err = rand.Read(data)
	require.NoError(t, err, "fail to read random data")
	key := []byte("testkey")

	meta, err := c.Write(key, bytes.NewReader(data))
	require.NoError(t, err, "fail write data")

	// Check status is ok after a write
	status, err := c.Check(*meta, false)
	require.NoError(t, err, "fail to check object")
	require.Equal(t, storage.CheckStatusOptimal, status)
	status, err = c.Check(*meta, true)
	require.NoError(t, err, "fail to check object")
	require.True(t, status == storage.CheckStatusValid || status == storage.CheckStatusOptimal)

	meta, err = c.metastorClient.GetMetadata(key)
	require.NoError(t, err)
	// corrupt file by removing blocks
	for i := 0; i < len(meta.Chunks); i += 4 {
		if i%4 == 0 {
			chunk := &meta.Chunks[i]
			store, err := datastorCluster.GetShard(chunk.Objects[0].ShardID)
			require.NoError(t, err)
			err = store.DeleteObject(chunk.Objects[0].Key)
			require.NoError(t, err)
		}
	}

	// Check status is corrupted
	status, err = c.Check(*meta, false)
	require.NoError(t, err, "fail to check object")
	require.Truef(t, status == storage.CheckStatusValid || status == storage.CheckStatusInvalid,
		"status=%v, %d", status, status)
}

func TestClientRepair(t *testing.T) {
	servers, serverClean := testZdbServer(t, 4)
	defer serverClean()

	shards := make([]datastor.ShardConfig, len(servers))
	for i, server := range servers {
		shards[i] = datastor.ShardConfig{Address: server.Address()}
	}

	config := newDefaultConfig(shards, 1024)

	tt := []struct {
		name string

		DataShardCount   int
		ParityShardCount int

		repairErr error
	}{
		{
			name:           "replication",
			DataShardCount: 4,
			repairErr:      nil,
		},
		{
			name:             "distribution",
			DataShardCount:   3,
			ParityShardCount: 1,
			repairErr:        nil,
		},
		{
			name:             "no-repair-suport",
			DataShardCount:   0,
			ParityShardCount: 0,
			repairErr:        ErrRepairSupport,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			config.DataStor.Pipeline.Distribution.DataShardCount = tc.DataShardCount
			config.DataStor.Pipeline.Distribution.ParityShardCount = tc.ParityShardCount
			testRepair(t, config, tc.repairErr)
		})
	}
}

func testRepair(t *testing.T, config Config, repairErr error) {
	c, datastorCluster, err := getTestClient(config)
	require.NoError(t, err, "fail to create client")
	defer c.Close()

	data := make([]byte, 1024*10)
	key := make([]byte, 64)

	_, err = rand.Read(data)
	require.NoError(t, err, "fail to read random data")
	_, err = rand.Read(key)
	require.NoError(t, err, "fail to read random key")

	_, err = c.Write(key, bytes.NewReader(data))
	require.NoError(t, err, "fail write data")

	meta, err := c.metastorClient.GetMetadata(key)
	require.NoError(t, err)
	require.EqualValues(t, len(data), meta.Size, "metadata.Size should be egals to the total size of the written file")
	// store last-write epoch, so we can compare it later after repair
	lastWriteEpoch := meta.LastWriteEpoch
	// store data size, so we can compare it later after repair
	size := meta.Size

	// Check status is ok after a write
	status, err := c.Check(*meta, false)
	require.NoError(t, err, "fail to check object")
	require.True(t, status == storage.CheckStatusValid || status == storage.CheckStatusOptimal)
	status, err = c.Check(*meta, true)
	require.NoError(t, err, "fail to check object")
	require.True(t, status == storage.CheckStatusValid || status == storage.CheckStatusOptimal)

	// corrupt file by removing a block
	store, err := datastorCluster.GetShard(meta.Chunks[0].Objects[0].ShardID)
	require.NoError(t, err)
	err = store.DeleteObject(meta.Chunks[0].Objects[0].Key)
	require.NoError(t, err)

	// Check status is corrupted
	status, err = c.Check(*meta, false)
	require.NoError(t, err, "fail to check object")
	require.True(t, status == storage.CheckStatusValid || status == storage.CheckStatusInvalid)

	// try to repair
	repairMeta, err := c.Repair(*meta)
	if repairErr != nil {
		require.Error(t, repairErr, err)
		return
	}
	require.NoError(t, err)

	// ensure the last-write epoch is updated
	require.NoError(t, err)
	require.NotNil(t, repairMeta)
	require.True(t, repairMeta.LastWriteEpoch != 0 && repairMeta.LastWriteEpoch != lastWriteEpoch)
	require.EqualValues(t, size, repairMeta.Size)
	require.Equal(t, meta.Key, repairMeta.Key)

	// make sure we can read the data again
	buf := bytes.NewBuffer(nil)
	err = c.Read(*repairMeta, buf)
	require.NoError(t, err)
	readData := buf.Bytes()
	require.Equal(t, data, readData, "restored data is not the same as initial data")
}

func TestClient_ExplicitErrors(t *testing.T) {
	require := require.New(t)

	servers, serverClean := testZdbServer(t, 1)
	defer serverClean()

	dataShards := []datastor.ShardConfig{{Address: servers[0].Address()}}
	config := newDefaultConfig(dataShards, 0)
	config.DataStor.Pipeline.Distribution = pipeline.ObjectDistributionConfig{}

	cli, _, err := getTestClient(config)
	require.NoError(err)
	defer cli.Close()

	_, err = cli.Write(nil, nil)
	require.Error(err, "no key or reader given")
	_, err = cli.Write([]byte("foo"), nil)
	require.Error(err, "no reader given")
	_, err = cli.Write(nil, bytes.NewReader(nil))
	require.Error(err, "no key given")

	err = cli.Read(metatypes.Metadata{}, nil)
	require.Error(err, "no key or writer given")

	require.NoError(cli.Close())
}
