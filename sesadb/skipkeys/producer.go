package skipkeys

import "github.com/sesanetwork/go-vassalo/sesadb"

func openDB(p sesadb.DBProducer, skipPrefix []byte, name string) (sesadb.Store, error) {
	store, err := p.OpenDB(name)
	if err != nil {
		return nil, err
	}
	return &Store{store, skipPrefix}, nil
}

type AllDBProducer struct {
	sesadb.FullDBProducer
	skipPrefix []byte
}

func WrapAllProducer(p sesadb.FullDBProducer, skipPrefix []byte) *AllDBProducer {
	return &AllDBProducer{
		FullDBProducer: p,
		skipPrefix:     skipPrefix,
	}
}

func (p *AllDBProducer) OpenDB(name string) (sesadb.Store, error) {
	return openDB(p.FullDBProducer, p.skipPrefix, name)
}

type DBProducer struct {
	sesadb.DBProducer
	skipPrefix []byte
}

func WrapProducer(p sesadb.DBProducer, skipPrefix []byte) *DBProducer {
	return &DBProducer{
		DBProducer: p,
		skipPrefix: skipPrefix,
	}
}

func (p *DBProducer) OpenDB(name string) (sesadb.Store, error) {
	return openDB(p.DBProducer, p.skipPrefix, name)
}
