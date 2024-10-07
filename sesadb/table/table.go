package table

import (
	"math/big"

	"github.com/sesanetwork/go-sesa/common"

	"github.com/sesanetwork/go-vassalo/sesadb"
)

// Table wraper the underling DB, so all the table's data is stored with a prefix in underling DB
type Table struct {
	IteratedReader
	underlying sesadb.Store
}

var (
	// NOTE: key collisions are possible
	separator = []byte{}
)

// prefixed key (prefix + separator + key)
func prefixed(key, prefix []byte) []byte {
	prefixedKey := make([]byte, 0, len(prefix)+len(separator)+len(key))
	prefixedKey = append(prefixedKey, prefix...)
	prefixedKey = append(prefixedKey, separator...)
	prefixedKey = append(prefixedKey, key...)
	return prefixedKey
}

func noPrefix(key, prefix []byte) []byte {
	if len(key) < len(prefix)+len(separator) {
		return key
	}
	return key[len(prefix)+len(separator):]
}

/*
 * Database
 */

func New(db sesadb.Store, prefix []byte) *Table {
	return &Table{
		IteratedReader: IteratedReader{
			prefix:     prefix,
			underlying: db,
		},
		underlying: db,
	}
}

func (t *Table) NewTable(prefix []byte) *Table {
	return New(t, prefix)
}

func (t *Table) Close() error {
	return sesadb.ErrUnsupportedOp
}

func (t *Table) Drop() {}

func (t *Table) Put(key []byte, value []byte) error {
	return t.underlying.Put(prefixed(key, t.prefix), value)
}

func (t *Table) Delete(key []byte) error {
	return t.underlying.Delete(prefixed(key, t.prefix))
}

func (t *Table) NewBatch() sesadb.Batch {
	return &batch{t.underlying.NewBatch(), t.prefix}
}

func incPrefix(prefix []byte) []byte {
	if len(prefix) == 0 {
		return nil
	}
	endBn := new(big.Int).SetBytes(prefix)
	endBn.Add(endBn, common.Big1)
	if len(endBn.Bytes()) > len(prefix) {
		// overflow
		return nil
	}
	res := make([]byte, len(prefix)-len(endBn.Bytes()), len(prefix))
	return append(res, endBn.Bytes()...)
}

func (t *Table) Compact(start []byte, limit []byte) error {
	end := prefixed(limit, t.prefix)
	if limit == nil {
		end = incPrefix(t.prefix)
	}
	return t.underlying.Compact(prefixed(start, t.prefix), end)
}

/*
 * Batch
 */

type batch struct {
	batch  sesadb.Batch
	prefix []byte
}

func (b *batch) Put(key, value []byte) error {
	return b.batch.Put(prefixed(key, b.prefix), value)
}

func (b *batch) Delete(key []byte) error {
	return b.batch.Delete(prefixed(key, b.prefix))
}

func (b *batch) ValueSize() int {
	return b.batch.ValueSize()
}

func (b *batch) Write() error {
	return b.batch.Write()
}

func (b *batch) Reset() {
	b.batch.Reset()
}

func (b *batch) Replay(w sesadb.Writer) error {
	return b.batch.Replay(&replayer{w, b.prefix})
}

/*
 * Replayer
 */

type replayer struct {
	writer sesadb.Writer
	prefix []byte
}

func (r *replayer) Put(key, value []byte) error {
	return r.writer.Put(noPrefix(key, r.prefix), value)
}

func (r *replayer) Delete(key []byte) error {
	return r.writer.Delete(noPrefix(key, r.prefix))
}
