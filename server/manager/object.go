package manager

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/zero-os/0-stor/server/db"
)

type ObjectManager struct {
	db        db.DB
	namespace string
}

func NewObjectManager(namespace string, db db.DB) *ObjectManager {
	return &ObjectManager{
		namespace: namespace,
		db:        db,
	}
}

func objKey(namespace string, key []byte) []byte {
	return []byte(fmt.Sprintf("%s:%s", namespace, key))
}

// Set saved an object into the key-value store as a blob of bytes
func (mgr *ObjectManager) Set(key []byte, data []byte, referenceList []string) error {
	obj := db.NewObjet()
	obj.Data = data
	for i := range referenceList {
		copy(obj.ReferenceList[i][:], []byte(referenceList[i]))
	}

	b, err := obj.Encode()
	if err != nil {
		return err
	}

	k := objKey(mgr.namespace, key)
	log.Debugf("set objet %s into namespace %s", string(k), mgr.namespace)
	return mgr.db.Set(k, b)
}

func (mgr *ObjectManager) List(start, count int) ([][]byte, error) {
	prefix := fmt.Sprintf("%s:", mgr.namespace)
	keys, err := mgr.db.List([]byte(prefix))
	if err != nil {
		return nil, err
	}

	// remove namespace prefix
	for i := range keys {
		keys[i] = keys[i][len(mgr.namespace)+1:]
	}
	return keys, nil
}

func (mgr *ObjectManager) Get(key []byte) (*db.Object, error) {
	b, err := mgr.db.Get(objKey(mgr.namespace, key))
	if err != nil {
		return nil, err
	}

	obj := db.NewObjet()
	err = obj.Decode(b)
	return obj, err
}

func (mgr *ObjectManager) Delete(key []byte) error {
	return mgr.db.Delete(objKey(mgr.namespace, key))
}

func (mgr *ObjectManager) Exists(key []byte) (bool, error) {
	return mgr.db.Exists(objKey(mgr.namespace, key))
}
