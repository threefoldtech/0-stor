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
	"testing"

	"github.com/zero-os/0-stor/client/datastor"
	"github.com/zero-os/0-stor/client/datastor/pipeline"
	zdbtest "github.com/zero-os/0-stor/client/datastor/zerodb/test"
	dbp "github.com/zero-os/0-stor/client/metastor/db"
	"github.com/zero-os/0-stor/client/metastor/db/etcd"
	"github.com/zero-os/0-stor/client/metastor/db/test"

	"github.com/stretchr/testify/require"
)

type testServer struct {
	addr string
}

func (ts *testServer) Address() string {
	return ts.addr
}

func testZdbServer(t testing.TB, n int) (servers []*testServer, cleanups func()) {
	require := require.New(t)

	var (
		namespace    = "ns"
		cleanupFuncs []func()
	)

	for i := 0; i < n; i++ {
		addr, cleanup, err := zdbtest.NewInMem0DBServer(namespace)
		require.NoError(err)
		cleanupFuncs = append(cleanupFuncs, cleanup)
		servers = append(servers, &testServer{
			addr: addr,
		})
	}

	cleanups = func() {
		for _, cleanup := range cleanupFuncs {
			cleanup()
		}
	}
	return
}

func getTestClient(cfg Config) (*Client, datastor.Cluster, error) {
	// create datastor cluster
	datastorCluster, err := createDataClusterFromConfig(&cfg)
	if err != nil {
		return nil, nil, err
	}

	// create data pipeline, using our datastor cluster
	dataPipeline, err := pipeline.NewPipeline(cfg.DataStor.Pipeline, datastorCluster, -1)
	if err != nil {
		return nil, nil, err
	}

	// create ETCD or in-memory db
	var db dbp.DB
	if len(cfg.MetaStor.Database.Endpoints) != 0 {
		db, err = etcd.New(cfg.MetaStor.Database.Endpoints)
		if err != nil {
			return nil, nil, err
		}
	} else {
		db = test.New()
	}

	// create metastor client first,
	// and than create our master 0-stor client with all features.
	metastorClient, err := createMetastorClientFromConfigAndDatabase(cfg.Namespace, &cfg.MetaStor, db)
	if err != nil {
		return nil, nil, err
	}
	return NewClient(metastorClient, dataPipeline), datastorCluster, nil
}
