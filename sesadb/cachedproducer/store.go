package cachedproducer

import "github.com/sesanetwork/go-vassalo/sesadb"

type StoreWithFn struct {
	sesadb.Store
	CloseFn func() error
	DropFn  func()
}

func (s *StoreWithFn) Close() error {
	return s.CloseFn()
}

func (s *StoreWithFn) Drop() {
	s.DropFn()
}
