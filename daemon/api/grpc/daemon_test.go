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

package grpc

import (
	"testing"

	"github.com/threefoldtech/0-stor/client/datastor/pipeline"
	"github.com/threefoldtech/0-stor/client/metastor"

	"github.com/stretchr/testify/require"
)

func TestConfig_ValidateAndSanitize(t *testing.T) {
	var cfg Config

	err := cfg.validateAndSanitize()
	require.Error(t, err)

	cfg.Pipeline = new(pipeline.SingleObjectPipeline)
	err = cfg.validateAndSanitize()
	require.NoError(t, err)
	require.Equal(t, DefaultMaxMsgSize, cfg.MaxMsgSize)

	cfg.MetaClient = new(metastor.Client)
	cfg.MaxMsgSize = 0

	err = cfg.validateAndSanitize()
	require.NoError(t, err)

	require.Equal(t, DefaultMaxMsgSize, cfg.MaxMsgSize)
}
