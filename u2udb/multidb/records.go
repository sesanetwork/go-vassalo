package multidb

import (
	"github.com/sesanetwork/go-sesa/rlp"

	"github.com/sesanetwork/go-vassalo/sesadb"
)

type TableRecord struct {
	Req   string
	Table string
}

func ReadTablesList(store sesadb.Store, tableRecordsKey []byte) (res []TableRecord, err error) {
	b, err := store.Get(tableRecordsKey)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return []TableRecord{}, nil
	}
	err = rlp.DecodeBytes(b, &res)
	return
}

func WriteTablesList(store sesadb.Store, tableRecordsKey []byte, records []TableRecord) error {
	b, err := rlp.EncodeToBytes(records)
	if err != nil {
		return err
	}
	return store.Put(tableRecordsKey, b)
}
