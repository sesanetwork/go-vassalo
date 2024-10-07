package readonlystore

import "github.com/sesanetwork/go-vassalo/sesadb"

type Store struct {
	sesadb.Store
}

func Wrap(s sesadb.Store) *Store {
	return &Store{s}
}

// Put inserts the given value into the key-value data store.
func (s *Store) Put(key []byte, value []byte) error {
	return sesadb.ErrUnsupportedOp
}

// Delete removes the key from the key-value data store.
func (s *Store) Delete(key []byte) error {
	return sesadb.ErrUnsupportedOp
}

type Batch struct {
	sesadb.Batch
}

func (s *Store) NewBatch() sesadb.Batch {
	return &Batch{s.Store.NewBatch()}
}

// Put inserts the given value into the key-value data store.
func (s *Batch) Put(key []byte, value []byte) error {
	return sesadb.ErrUnsupportedOp
}

// Delete removes the key from the key-value data store.
func (s *Batch) Delete(key []byte) error {
	return sesadb.ErrUnsupportedOp
}
