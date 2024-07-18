package storage

import (
	"context"

	"github.com/dgraph-io/badger/v2"
	ds "github.com/ipfs/go-datastore"
	badgerds "github.com/ipfs/go-ds-badger2"
)

var _ StorageDB = (*BadgerDBWrapper)(nil)

// BadgerDBWrapper wraps the BadgerDB to implement the StorageDB interface.
type BadgerDBWrapper struct {
	ds *badgerds.Datastore
}

func NewBadgerDBWrapper(datastorePath string, options *badgerds.Options) (*BadgerDBWrapper, error) {
	ds, err := badgerds.NewDatastore(datastorePath, options)
	if err != nil {
		return nil, err
	}

	return &BadgerDBWrapper{ds}, nil
}

func (b *BadgerDBWrapper) Datastore() ds.Batching {
	return b.ds
}

func (b *BadgerDBWrapper) Keys(prefix []byte) ([][]byte, error) {
	var keys [][]byte

	err := b.ds.DB.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.IteratorOptions{
			PrefetchValues: false,
			Prefix:         prefix,
		})
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			keys = append(keys, it.Item().KeyCopy(nil))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func (b *BadgerDBWrapper) CollectGarbage(ctx context.Context) error {
	return b.ds.CollectGarbage(ctx)
}

func (b *BadgerDBWrapper) Get(key []byte) (StorageItem, error) {
	var item *badger.Item
	var err error
	err = b.ds.DB.View(func(txn *badger.Txn) error {
		item, err = txn.Get(key)
		if err != nil {
			return err
		}
		return nil
	})
	return item, err
}

func (b *BadgerDBWrapper) Set(key, val []byte) error {
	return b.ds.DB.Update(func(txn *badger.Txn) error {
		return txn.Set(key, val)
	})
}

func (b *BadgerDBWrapper) Delete(key []byte) error {
	return b.ds.DB.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
}

func (b *BadgerDBWrapper) Close() error {
	return b.ds.Close()
}
