package provider

import (
	"github.com/dgraph-io/badger/v2"
	"projectionist/models"
	"testing"
)

func NewTestDB(t *testing.T, logEnable bool, readOnly bool) *badger.DB {
	opts := badger.DefaultOptions("").WithInMemory(true).
		WithEventLogging(logEnable).
		WithReadOnly(readOnly)

	db, err := badger.Open(opts)
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func prepareData(db *badger.DB, entrys []*badger.Entry) error {
	return db.Update(func(txn *badger.Txn) error {
		for _, entry := range entrys {
			err := txn.SetEntry(entry)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func marshalModel(t *testing.T, m models.Model) []byte {
	d, err := json.Marshal(m)
	if err != nil {
		t.Errorf("marshalModel error: %v", err)
	}

	return d
}
