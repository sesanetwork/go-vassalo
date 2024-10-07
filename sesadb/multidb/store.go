package multidb

import "github.com/sesanetwork/go-vassalo/sesadb"

type closableTable struct {
	sesadb.Store
	underlying sesadb.Store
	noDrop     bool
}

// Close leaves underlying database.
func (s *closableTable) Close() error {
	return s.underlying.Close()
}

// Drop whole database.
func (s *closableTable) Drop() {
	if s.noDrop {
		return
	}
	s.underlying.Drop()
}
