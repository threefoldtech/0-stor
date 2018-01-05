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

package test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRoundTrip(t *testing.T) {
	c := NewClient()
	require.NotNil(t, c)
	defer c.Close()

	RoundTrip(t, c)
}

func TestClientNilKeys(t *testing.T) {
	c := NewClient()
	require.NotNil(t, c)
	defer c.Close()

	ClientNilKeys(t, c)
}

func TestClientUpdate(t *testing.T) {
	c := NewClient()
	require.NotNil(t, c)
	defer c.Close()

	ClientUpdate(t, c)
}

func TestClientUpdateAsync(t *testing.T) {
	c := NewClient()
	require.NotNil(t, c)
	defer c.Close()

	ClientUpdateAsync(t, c)
}