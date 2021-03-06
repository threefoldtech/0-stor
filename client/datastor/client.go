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

package datastor

import (
	"context"
)

// Client defines the API for any client,
// used to interface with a zstordb server.
// It allows you to manage objects,
// as well as get information about them and their namespaces.
//
// All operations work within a namespace,
// which is defined by the label given when creating
// this  client.
//
// If the server requires authentication,
// this will have to be configured when creating the client as well, otherwise the methods of this interface will fail.
//
// Errors that can be returned for all methods:
//
type Client interface {
	// Creates an object, using the given data.
	// The key of the object is generated by the server, and returned.
	CreateObject(data []byte) (key []byte, err error)

	// Get an existing object, linked to a given key.
	//
	// ErrKeyNotFound is returned in case the requested key couldn't be found.
	// ErrObjectCorrupted is returned in case the stored object is corrupted.
	GetObject(key []byte) (*Object, error)

	// DeleteObject deletes an object, using a given key.
	// Deleting an non-existing object is considered valid.
	DeleteObject(key []byte) error

	// GetObjectStatus returns the status of an object,
	// indicating whether it's OK, missing or corrupt.
	GetObjectStatus(key []byte) (ObjectStatus, error)

	// ExistObject returns whether or not an object exists.
	//
	// ErrObjectCorrupted is returned in case the object key exists,
	// but the object is corrupted.
	ExistObject(key []byte) (bool, error)

	// ListObjectKeyIterator returns an iterator,
	// from which the keys of all stored objects within the namespace
	// (identified by the given label), an be retrieved.
	//
	// In case an error while the iterator is active,
	// it will be returned as part of the last returned result,
	// which is then considered to be invalid.
	// When an error is returned, as part of a result,
	// the iterator channel will be automatically closed as soon
	// as that item is received.
	ListObjectKeyIterator(ctx context.Context) (<-chan ObjectKeyResult, error)

	// GetNamespace returns the available information of a namespace.
	//
	// ErrKeyNotFound is returned in case no
	// stored namespace exist for the used label.
	GetNamespace() (*Namespace, error)

	// Utilization return the amount of bytes stored in
	// the namespace the client is connected to
	Utilization() int64

	// Close any open resources.
	Close() error
}
