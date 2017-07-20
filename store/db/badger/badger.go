package badger

import (
	"os"

	log "github.com/Sirupsen/logrus"
	badgerkv "github.com/dgraph-io/badger"
	"github.com/zero-os/0-stor/store/db"
)

var _ db.DB = (*BadgerDB)(nil)

// BadgerDB implements the db.DB interace
type BadgerDB struct {
	KV *badgerkv.KV
	// Config *config.Settings
}

// Constructor
func New(data, meta string) (*BadgerDB, error) {
	// log.Println("Initializing db directories")

	if err := os.MkdirAll(meta, 0774); err != nil {
		log.Errorf("\t\tMeta dir: %v [ERROR]", meta)
		return nil, err
	}

	// log.Printf("\t\tMeta dir: %v [SUCCESS]", meta)

	if err := os.MkdirAll(data, 0774); err != nil {
		log.Errorf("\t\tData dir: %v [ERROR]", data)
		return nil, err
	}

	// log.Printf("\t\tData dir: %v [SUCCESS]", data)

	opts := badgerkv.DefaultOptions
	opts.Dir = meta
	opts.ValueDir = data
	opts.SyncWrites = true

	kv, err := badgerkv.NewKV(&opts)

	// if err == nil {
	// 	log.Println("Loading db [SUCCESS]")
	// } else {
	// 	log.Println("Loading db [ERROR]")
	// }

	return &BadgerDB{
		KV: kv,
	}, err
}

func (b BadgerDB) Close() error {
	err := b.KV.Close()
	if err != nil {
		log.Errorln(err.Error())
	}
	return err
}

func (b BadgerDB) Delete(key string) error {
	err := b.KV.Delete([]byte(key))
	if err != nil {
		log.Errorln(err.Error())
	}
	return err
}

func (b BadgerDB) Set(key string, val []byte) error {
	err := b.KV.Set([]byte(key), val)
	if err != nil {
		log.Errorln(err.Error())
	}
	return err
}

func (b BadgerDB) Get(key string) ([]byte, error) {
	var item badgerkv.KVItem

	err := b.KV.Get([]byte(key), &item)

	if err != nil {
		log.Errorln(err.Error())
		return nil, err
	}

	v := item.Value()

	if len(v) == 0 {
		err = db.ErrNotFound
	}

	return v, err
}

func (b BadgerDB) Exists(key string) (bool, error) {
	exists, err := b.KV.Exists([]byte(key))
	if err != nil {
		log.Errorln(err.Error())
	}
	return exists, err
}

// Pass count = -1 to get all elements starting from the provided index
func (b BadgerDB) Filter(prefix string, start int, count int) ([][]byte, error) {
	opt := badgerkv.DefaultIteratorOptions

	it := b.KV.NewIterator(opt)
	defer it.Close()

	result := [][]byte{}

	counter := 0 // Number of namespaces encountered

	prefixBytes := []byte(prefix)

	for it.Seek(prefixBytes); it.ValidForPrefix(prefixBytes); it.Next() {
		counter++

		// Skip until starting index
		if counter < start {
			continue
		}

		item := it.Item()
		value := item.Value()
		result = append(result, value)

		if count > 0 && len(result) == count {
			break
		}
	}

	return result, nil
}

func (b BadgerDB) List(prefix string) ([]string, error) {
	opt := badgerkv.DefaultIteratorOptions
	opt.FetchValues = false

	it := b.KV.NewIterator(opt)
	defer it.Close()

	result := []string{}

	prefixBytes := []byte(prefix)

	for it.Seek(prefixBytes); it.ValidForPrefix(prefixBytes); it.Next() {
		item := it.Item()
		key := string(item.Key()[:])
		result = append(result, key)
	}

	return result, nil
}