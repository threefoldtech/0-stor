package client

import (
	"github.com/zero-os/0-stor/client/config"
	"github.com/zero-os/0-stor/client/lib/block"
	"github.com/zero-os/0-stor/client/meta"
	"github.com/zero-os/0-stor/client/pipe"
)

// Client defines 0-stor client
type Client struct {
	conf       *config.Config
	metaCli    *meta.Client
	storWriter block.Writer
	storReader block.Reader
}

// New creates new client
func New(conf *config.Config) (*Client, error) {
	// stor writer
	storWriter, err := pipe.NewWritePipe(conf, nil)
	if err != nil {
		return nil, err
	}

	storReader, err := pipe.NewReadPipe(conf)
	if err != nil {
		return nil, err
	}

	// meta client
	metaCli, err := meta.NewClient(conf.MetaShards)
	if err != nil {
		return nil, err
	}

	return &Client{
		conf:       conf,
		metaCli:    metaCli,
		storWriter: storWriter,
		storReader: storReader,
	}, nil

}

// Store stores payload with key=key
func (c *Client) Store(key, payload []byte) error {
	resp := c.storWriter.WriteBlock(payload)
	if resp.Err != nil {
		return resp.Err
	}
	if resp.Meta == nil {
		return nil
	}
	return c.metaCli.Put(string(key), resp.Meta)
}

// Get fetch data for given key
func (c *Client) Get(key []byte) ([]byte, error) {
	// get the meta
	meta, err := c.metaCli.Get(string(key))
	if err != nil {
		return nil, err
	}

	// decode the meta
	metaBytes, err := meta.Bytes()
	if err != nil {
		return nil, err
	}

	return c.storReader.ReadBlock(metaBytes)
}