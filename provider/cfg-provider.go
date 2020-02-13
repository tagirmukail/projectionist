package provider

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dgraph-io/badger"
	jsoniter "github.com/json-iterator/go"
	"google.golang.org/grpc/grpclog"

	"projectionist/models"
	"projectionist/utils/errors"
)

const (
	sep          = "|"
	MaxID string = "MaxID"

	indexID        = 0
	indexName      = 1
	indexIsDeleted = 2
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type CfgProvider struct {
	db    *badger.DB
	maxID int
}

func NewCfgProvider(db *badger.DB) (*CfgProvider, error) {
	var cfgProvider = &CfgProvider{
		db: db,
	}

	err := db.View(func(txn *badger.Txn) error {
		id, err := getMaxID(txn)
		if err != nil && err != badger.ErrKeyNotFound {
			return err
		}

		cfgProvider.maxID = id

		return nil
	})

	if err != nil {
		return nil, err
	}

	return cfgProvider, nil
}

func getMaxID(txn *badger.Txn) (int, error) {
	var valCopy []byte
	item, err := txn.Get([]byte(MaxID))
	if err != nil {
		if err == badger.ErrKeyNotFound {
			grpclog.Infof("max id not exist")
		}
		return 0, err
	}

	err = item.Value(func(val []byte) error {
		valCopy = append([]byte{}, val...)
		return nil
	})
	if err != nil {
		return 0, err
	}

	idStr := string(valCopy)
	return strconv.Atoi(idStr)
}

func (c *CfgProvider) GetDB() interface{} {
	return c.db
}

func find(txn *badger.Txn, keyPrefStr string) *badger.Item {
	iter := txn.NewIterator(badger.DefaultIteratorOptions)
	defer iter.Close()

	var item *badger.Item
	keyPref := []byte(keyPrefStr)
	for iter.Seek(keyPref); iter.ValidForPrefix(keyPref); iter.Next() {
		item = iter.Item()
		if item != nil {
			break
		}
	}

	return item
}

func (c *CfgProvider) Save(m models.Model) error {
	var err error
	c.maxID++
	defer func(err error) {
		if err != nil {
			c.maxID--
		}
	}(err)

	m.SetID(c.maxID)

	txn := c.db.NewTransaction(true)
	defer txn.Discard()

	data, err := json.Marshal(m)
	if err != nil {
		return err
	}

	item := find(txn, buildKeyPref(m))
	if item != nil {
		err = errors.Newf(
			1,
			500,
			"config with this name already exist",
			"config model with key %v already exist in kv db", m.GetName())
		return err
	}

	entryNameToData := badger.NewEntry([]byte(buildKey(
		m,
	)), data)
	err = txn.SetEntry(entryNameToData)
	if err != nil {
		return err
	}

	err = txn.Set([]byte(MaxID), []byte(strconv.Itoa(c.maxID)))
	if err != nil {
		return err
	}

	return txn.Commit()
}

func (c *CfgProvider) GetByName(_ models.Model, name string) (models.Model, error) {
	var valCopy []byte
	err := c.db.View(func(txn *badger.Txn) error {
		var err error
		iter := txn.NewIterator(badger.DefaultIteratorOptions)
		defer iter.Close()
		for iter.Rewind(); iter.Valid(); iter.Next() {
			item := iter.Item()
			key := string(item.Key())
			if !strings.Contains(key, name) {
				continue
			}

			valCopy, err = item.ValueCopy(nil)
			if err != nil {
				return err
			}

			return nil
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	if len(valCopy) == 0 {
		return nil, errors.Newf(
			2,
			404,
			"configuration not exist",
			"configuration by name %s not exist",
			name,
		)
	}

	conf := &models.Configuration{}
	return conf, json.Unmarshal(valCopy, conf)
}

func (c *CfgProvider) GetByID(_ models.Model, id int64) (models.Model, error) {
	var valCopy []byte
	err := c.db.View(func(txn *badger.Txn) error {
		var err error

		item := find(txn, strconv.FormatInt(id, 10))
		if item == nil {
			return errors.ErrNotExist
		}

		valCopy, err = item.ValueCopy(nil)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	if len(valCopy) == 0 {
		return nil, errors.ErrNotExist
	}

	conf := &models.Configuration{}
	return conf, json.Unmarshal(valCopy, conf)
}

func (c *CfgProvider) IsExistByName(m models.Model) (error, bool) {
	err := c.db.View(func(txn *badger.Txn) error {
		iter := txn.NewIterator(badger.DefaultIteratorOptions)
		defer iter.Close()
		for iter.Rewind(); iter.Valid(); iter.Next() {
			item := iter.Item()
			key := string(item.Key())
			if !strings.Contains(key, m.GetName()) {
				continue
			}

			return nil
		}

		return errors.ErrNotExist
	})

	return err, err == nil
}

func (c *CfgProvider) Count(_ models.Model) (int, error) {
	var count int
	return count, c.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			count++
		}

		return nil
	})
}

func (c *CfgProvider) Pagination(m models.Model, start, stop int) ([]models.Model, error) {
	var result []models.Model

	if start < 0 {
		start = 0
	}

	if start > stop {
		return result, errors.Newf(
			3,
			400,
			"pagination invalid",
			"invalid pagination: start:%d and stop:%d",
			strconv.Itoa(start), strconv.Itoa(stop),
		)
	}

	return result, c.db.View(func(txn *badger.Txn) error {
		var err error
		iter := txn.NewIterator(badger.DefaultIteratorOptions)
		defer iter.Close()
		var count int
		var valCopy []byte
		for iter.Rewind(); iter.Valid(); iter.Next() {
			item := iter.Item()
			key := string(item.Key())
			keyParts := strings.Split(key, sep)
			if key == MaxID || len(keyParts) != 3 {
				continue
			}

			count++
			if count > start && count <= stop {
				valCopy, err = item.ValueCopy(nil)
				if err != nil {
					return err
				}
				conf := &models.Configuration{}
				err = json.Unmarshal(valCopy, conf)
				if err != nil {
					return err
				}
				result = append(result, conf)
			}
		}

		return nil
	})
}

func (c *CfgProvider) Update(m models.Model, id int) error {
	txn := c.db.NewTransaction(true)
	defer txn.Discard()

	item := find(txn, strconv.Itoa(id))
	if item == nil {
		return errors.ErrNotExist
	}

	m.SetID(id)
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}

	err = txn.Delete(item.Key())
	if err != nil {
		return err
	}

	err = txn.Set([]byte(buildKey(m)), data)
	if err != nil {
		return err
	}

	return txn.Commit()
}

func (c *CfgProvider) Delete(m models.Model, id int) error {
	txn := c.db.NewTransaction(true)
	defer txn.Discard()

	item := find(txn, strconv.Itoa(id))
	if item == nil {
		return errors.ErrNotExist
	}

	err := txn.Delete(item.Key())
	if err != nil {
		return err
	}

	return txn.Commit()
}

// buildKey build key (id|name|0 if not deleted, 1 if deleted), example: 23|app1|0
func buildKey(m models.Model) string {
	var b = strings.Builder{}
	b.WriteString(strconv.Itoa(m.GetID()))
	b.WriteString(sep)
	b.WriteString(m.GetName())
	b.WriteString(sep)
	if m.IsDeleted() {
		b.WriteString(strconv.Itoa(1))
	} else {
		b.WriteString(strconv.Itoa(0))
	}
	return b.String()
}

// buildKeyPref build key pref by id and name, example: 23|app1
func buildKeyPref(m models.Model) string {
	var b = strings.Builder{}
	b.WriteString(strconv.Itoa(m.GetID()))
	b.WriteString(sep)
	b.WriteString(m.GetName())
	b.WriteString(sep)
	return b.String()
}

func decodeKey(key string) (models.Model, error) {
	pairs := strings.Split(key, sep)
	if len(pairs) != 3 {
		return nil, fmt.Errorf("invalid key: %s", key)
	}

	idStr, name, delStr := pairs[indexID], pairs[indexName], pairs[indexIsDeleted]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, err
	}

	delN, err := strconv.Atoi(delStr)
	if err != nil {
		return nil, err
	}

	return &models.Configuration{
		ID:      id,
		Name:    name,
		Deleted: delN,
	}, nil
}
