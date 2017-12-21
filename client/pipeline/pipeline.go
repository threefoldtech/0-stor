package pipeline

import (
	"io"

	"github.com/zero-os/0-stor/client/metastor"
	"github.com/zero-os/0-stor/client/pipeline/crypto"
	"github.com/zero-os/0-stor/client/pipeline/processing"
	"github.com/zero-os/0-stor/client/pipeline/storage"
)

// Pipeline defines the interface to write and read content
// to/from a zstordb cluster.
//
// Prior to storage content can be
// processed (compressed and/or encrypted), as well as split
// (into smaller chunks) and distributed in terms of replication
// or erasure coding.
//
// Content written in one way,
// has to be read in a way that is compatible. Meaning that if
// content was compressed and encrypted using a certain configuration,
// it will have to be decrypted and decompressed using that same configuration,
// or else the content will not be able to be read.
type Pipeline interface {
	// Write content to a zstordb cluster,
	// the details depend upon the specific implementation.
	Write(r io.Reader) ([]metastor.Chunk, error)
	// Read content from a zstordb cluster,
	// the details depend upon the specific implementation.
	Read(chunks []metastor.Chunk, w io.Writer) error

	// GetChunkStorage returns the underlying and used ChunkStorage
	GetChunkStorage() storage.ChunkStorage
}

// Constructor types which are used to create unique instances of the types involved,
// for each branch (goroutine) of a pipeline.
type (
	// HasherConstructor is a constructor type which is used to create a unique
	// Hasher for each goroutine where the Hasher is needed within a pipeline.
	// This is required as a (crypto) Hasher is not thread-safe.
	HasherConstructor func() (crypto.Hasher, error)
	// ProcessorConstructor is a constructor type which is used to create a unique
	// Processor for each goroutine where the Processor is needed within a pipeline.
	// This is required as a Processor is not thread-safe.
	ProcessorConstructor func() (processing.Processor, error)
)

// DefaultHasherConstructor is an implementation of a HasherConstructor,
// which can be used as a safe default HasherConstructor, by pipeline implementations,
// should such a constructor not be given by the user.
func DefaultHasherConstructor() (crypto.Hasher, error) {
	return crypto.NewDefaultHasher256(nil)
}

// DefaultProcessorConstructor is an implementation of a ProcessorConstructor,
// which can be used as a safe default ProcessorConstructor, by pipeline implementations,
// should such a constructor not be given by the user.
func DefaultProcessorConstructor() (processing.Processor, error) {
	return processing.NopProcessor{}, nil
}
