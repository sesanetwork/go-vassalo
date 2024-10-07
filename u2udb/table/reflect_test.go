package table

//go:generate go run github.com/golang/mock/mockgen -package=table -destination=mock_test.go github.com/Fantom-foundation/lachesis-base/sesadb DBProducer,DropableStore

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	sesadb "github.com/sesanetwork/go-vassalo/sesadb"
)

type testTables struct {
	NoTable interface{}
	Manual  sesadb.Store `table:"-"`
	Nil     sesadb.Store `table:"-"`
	Auto1   sesadb.Store `table:"A"`
	Auto2   sesadb.Store `table:"B"`
	Auto3   sesadb.Store `table:"C"`
}

func TestOpenTables(t *testing.T) {
	require := require.New(t)
	ctrl := gomock.NewController(t)

	prefix := "prefix"

	mockStore := func() sesadb.Store {
		store := NewMockDropableStore(ctrl)
		store.EXPECT().Close().
			Times(1).
			Return(nil)
		return store
	}

	dbs := NewMockDBProducer(ctrl)
	dbs.EXPECT().OpenDB(gomock.Any()).
		Times(3).
		DoAndReturn(func(name string) (sesadb.Store, error) {
			require.Contains(name, prefix)
			return mockStore(), nil
		})

	tt := &testTables{}

	// open auto
	err := OpenTables(tt, dbs, prefix)
	require.NoError(err)
	require.NotNil(tt.Auto1)
	require.NotNil(tt.Auto2)
	require.NotNil(tt.Auto3)
	require.Nil(tt.NoTable)
	require.Nil(tt.Nil)

	// open manual
	require.Nil(tt.Manual)
	tt.Manual = mockStore()
	require.NotNil(tt.Manual)

	// close all
	err = CloseTables(tt)
	require.NoError(err)
	require.NotNil(tt.Auto1)
	require.NotNil(tt.Auto2)
	require.NotNil(tt.Auto3)
	require.Nil(tt.NoTable)
	require.Nil(tt.Nil)
	require.NotNil(tt.Manual)
}
