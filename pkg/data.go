package BookQuest

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/blevesearch/bleve/v2"
	bolt "go.etcd.io/bbolt"
)

type Store struct {
	bleve bleve.Index
	*bolt.DB
	Bucket []byte
}

func exists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	} else {
		return false
	}
}

func (store *Store) createIndex(path string) error {
	var index bleve.Index
	var err error
	blevePath := filepath.Join(path, "data.bleve")
	if exists(blevePath) {
		index, err = bleve.Open(blevePath)
		if err != nil {
			return err
		}
	} else {
		mapping := bleve.NewIndexMapping()
		index, err = bleve.New(blevePath, mapping)
		if err != nil {
			return err
		}
	}
	store.bleve = index
	return nil
}

func (store *Store) createData(path string) error {
	boltPath := filepath.Join(path, "data.db")
	db, err := bolt.Open(boltPath, 0600, nil)
	if err != nil {
		return err
	}
	store.DB = db
	store.Bucket = []byte("pages")
	return store.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(store.Bucket)
		return err
	})
}

func NewStore(path string) (*Store, error) {
	os.MkdirAll(path, os.ModePerm)
	store := Store{}
	if err := store.createData(path); err != nil {
		return nil, err
	}
	if err := store.createIndex(path); err != nil {
		return nil, err
	}
	return &store, nil
}

func (store *Store) Save(value Page, user Profile) error {
	value.Meta.LastModified = time.Now()
	value.Meta.LastModifiedBy = user.Id
	return store.update(value)
}

func (store *Store) update(value Page) error {
	rawData, err := json.Marshal(value)
	if err != nil {
		return err
	}
	store.bleve.Index(value.ID, value)
	return store.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(store.Bucket).Put([]byte(value.ID), rawData)
	})
}

func (store *Store) Get(key string) (Page, error) {
	var result Page
	err := store.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(store.Bucket).Get([]byte(key))
		return json.Unmarshal(b, &result)
	})
	return result, err
}

func (store *Store) UpdateStats(profile *Profile) error {
	profile.LinksAdded = 0
	profile.LinksClicked = 0
	profile.LinksFavoirted = 0
	err := store.View(func(tx *bolt.Tx) error {
		tx.Bucket(store.Bucket).ForEach(func(k, v []byte) error {
			var result Page
			err := json.Unmarshal(v, &result)
			if err != nil {
				return nil
			}
			if result.Meta.CreatedBy == profile.Id || result.Meta.LastModifiedBy == profile.Id {
				profile.LinksAdded++
			}
			for _, fav := range result.Favourite {
				if fav == profile.Id {
					profile.LinksFavoirted++
				}
			}
			return nil
		})
		return nil
	})
	return err
}

func (store *Store) Delete(key string) (Page, error) {
	var result Page
	err := store.View(func(tx *bolt.Tx) error {
		return tx.Bucket(store.Bucket).Delete([]byte(key))
	})
	store.bleve.Delete(key)
	return result, err
}

func (store *Store) Search(term string) ([]Page, error) {
	results := []Page{}
	query := bleve.NewMatchQuery(term)
	search := bleve.NewSearchRequest(query)
	searchResults, err := store.bleve.Search(search)
	if err != nil {
		return results, err
	}
	for _, hit := range searchResults.Hits {
		if page, err := store.Get(hit.ID); err == nil {
			results = append(results, page)
		}
	}
	return results, err
}
