package consensus

import (
	"errors"

	"github.com/sesanetwork/go-sesa/rlp"

	"github.com/sesanetwork/go-vassalo/native/idx"
	"github.com/sesanetwork/go-vassalo/sesadb"
	"github.com/sesanetwork/go-vassalo/sesadb/memorydb"
	"github.com/sesanetwork/go-vassalo/sesadb/table"
	"github.com/sesanetwork/go-vassalo/utils/simplewlru"
)

// Store is a consensus persistent storage working over parent key-value database.
type Store struct {
	getEpochDB EpochDBProducer
	cfg        StoreConfig
	crit       func(error)

	mainDB sesadb.Store
	table  struct {
		LastDecidedState sesadb.Store `table:"c"`
		EpochState       sesadb.Store `table:"e"`
	}

	cache struct {
		LastDecidedState *LastDecidedState
		EpochState       *EpochState
		FrameRoots       *simplewlru.Cache `cache:"-"` // store by pointer
	}

	epochDB    sesadb.Store
	epochTable struct {
		Roots          sesadb.Store `table:"r"`
		VectorIndex    sesadb.Store `table:"v"`
		ConfirmedEvent sesadb.Store `table:"C"`
	}
}

var (
	ErrNoGenesis = errors.New("genesis not applied")
)

type EpochDBProducer func(epoch idx.Epoch) sesadb.Store

// NewStore creates store over key-value db.
func NewStore(mainDB sesadb.Store, getDB EpochDBProducer, crit func(error), cfg StoreConfig) *Store {
	s := &Store{
		getEpochDB: getDB,
		cfg:        cfg,
		crit:       crit,
		mainDB:     mainDB,
	}

	table.MigrateTables(&s.table, s.mainDB)

	s.initCache()

	return s
}

func (s *Store) initCache() {
	s.cache.FrameRoots = s.makeCache(s.cfg.Cache.RootsNum, s.cfg.Cache.RootsFrames)
}

// NewMemStore creates store over memory map.
// Store is always blank.
func NewMemStore() *Store {
	getDb := func(epoch idx.Epoch) sesadb.Store {
		return memorydb.New()
	}
	cfg := LiteStoreConfig()
	crit := func(err error) {
		panic(err)
	}
	return NewStore(memorydb.New(), getDb, crit, cfg)
}

// Close leaves underlying database.
func (s *Store) Close() error {
	setnil := func() interface{} {
		return nil
	}

	table.MigrateTables(&s.table, nil)
	table.MigrateCaches(&s.cache, setnil)
	table.MigrateTables(&s.epochTable, nil)
	err := s.mainDB.Close()
	if err != nil {
		return err
	}

	if s.epochDB != nil {
		err = s.epochDB.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// dropEpochDB drops existing epoch DB
func (s *Store) dropEpochDB() error {
	prevDb := s.epochDB
	if prevDb != nil {
		err := prevDb.Close()
		if err != nil {
			return err
		}
		prevDb.Drop()
	}
	return nil
}

// openEpochDB makes new epoch DB
func (s *Store) openEpochDB(n idx.Epoch) error {
	// Clear full LRU cache.
	s.cache.FrameRoots.Purge()

	s.epochDB = s.getEpochDB(n)
	table.MigrateTables(&s.epochTable, s.epochDB)
	return nil
}

/*
 * Utils:
 */

// set RLP value
func (s *Store) set(table sesadb.Store, key []byte, val interface{}) {
	buf, err := rlp.EncodeToBytes(val)
	if err != nil {
		s.crit(err)
	}

	if err := table.Put(key, buf); err != nil {
		s.crit(err)
	}
}

// get RLP value
func (s *Store) get(table sesadb.Store, key []byte, to interface{}) interface{} {
	buf, err := table.Get(key)
	if err != nil {
		s.crit(err)
	}
	if buf == nil {
		return nil
	}

	err = rlp.DecodeBytes(buf, to)
	if err != nil {
		s.crit(err)
	}
	return to
}

func (s *Store) makeCache(weight uint, size int) *simplewlru.Cache {
	cache, err := simplewlru.New(weight, size)
	if err != nil {
		s.crit(err)
	}
	return cache
}
