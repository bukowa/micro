/*
Copyright © 2020 Mateusz Kurowski

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package bolt

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"reflect"
	"sync"

	"github.com/boltdb/bolt"
)

type Storage interface {
	Bolt() *bolt.DB
	Init(types ...Model) error

	Create(Model) error
	Delete(Model) error
	Pop(Model) error
	Get(Model) error
	Exists(Model) (bool, error)
	ForEach(m Model, f func(m Model)) error

	BucketFor(m Model, tx *bolt.Tx) (*bolt.Bucket, error)
	NextID(*bolt.Bucket) ([]byte, error)
	Stats(Model) (bolt.BucketStats, error)
}

// NewStorage creates or opens exiting Storage.
func NewStorage(opts *bolt.Options, path string, types ...Model) (Storage, error) {
	var err error
	var d = &dbx{}
	d.bolt, err = bolt.Open(path, 0600, opts)
	if err != nil {
		return nil, err
	}
	err = d.Init(types...)
	return d, err
}

type dbx struct {
	mu   sync.Mutex
	bolt *bolt.DB
}

func (db *dbx) Bolt() *bolt.DB {
	return db.bolt
}

func (db *dbx) Lock() {
	db.mu.Lock()
}

func (db *dbx) Unlock() {
	db.mu.Unlock()
}

func (db *dbx) Init(types ...Model) error {
	return db.bolt.Update(func(tx *bolt.Tx) error {
		for _, each := range types {
			name := getType(each)
			if _, err := tx.CreateBucketIfNotExists([]byte(name)); err != nil {
				return err
			}
		}
		return nil
	})
}

func (db *dbx) Get(m Model) (err error) {
	if err = checkKey(m); err != nil {
		return
	}
	err = db.bolt.View(func(tx *bolt.Tx) error {
		bucket, err := db.BucketFor(m, tx)
		if err != nil {
			return err
		}
		b := bucket.Get(m.Key())
		if len(b) == 0 {
			return ErrorNotFound
		}
		if err := json.Unmarshal(b, m); err != nil {
			return err
		}
		return nil
	})
	return
}

func (db *dbx) Create(m Model) (err error) {
	if err = checkKey(m); err != nil {
		return
	}
	var b []byte
	b, err = json.Marshal(m)
	if err != nil {
		return
	}
	err = db.bolt.Update(func(tx *bolt.Tx) error {
		bucket, err := db.BucketFor(m, tx)
		if err != nil {
			return err
		}
		return bucket.Put(m.Key(), b)
	})
	return
}

func (db *dbx) Delete(m Model) (err error) {
	if err = checkKey(m); err != nil {
		return
	}
	err = db.bolt.Update(func(tx *bolt.Tx) error {
		bucket, err := db.BucketFor(m, tx)
		if err != nil {
			return err
		}
		return bucket.Delete(m.Key())
	})
	return
}

func (db *dbx) Pop(m Model) error {
	var found bool
	err := db.bolt.Update(func(tx *bolt.Tx) error {
		b, err := db.BucketFor(m, tx)
		if err != nil {
			return err
		}

		c := b.Cursor()
		if k, v := c.First(); k != nil && v != nil {
			found = true
			if err := json.Unmarshal(v, m); err != nil {
				return err
			}
			m.SetKey(k)
			if err := b.Delete(k); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	if !found {
		return ErrorNotFound
	}
	return nil
}

func (db *dbx) Exists(m Model) (t bool, err error) {
	if err = checkKey(m); err != nil {
		return
	}
	err = db.bolt.View(func(tx *bolt.Tx) (err error) {
		bucket, err := db.BucketFor(m, tx)
		if bucket == nil {
			panic(fmt.Sprintf("bucket for %v is nil", m))
		}
		if err != nil {
			return err
		}
		b := bucket.Get(m.Key())
		if b != nil {
			t = true
		}
		return
	})
	return
}

func (db *dbx) Stats(m Model) (bs bolt.BucketStats, err error) {
	err = db.bolt.View(func(tx *bolt.Tx) error {
		bucket, err := db.BucketFor(m, tx)
		if err != nil {
			return err
		}
		bs = bucket.Stats()
		return nil
	})
	return
}

func (db *dbx) NextID(bucket *bolt.Bucket) (b []byte, err error) {
	seq, err := bucket.NextSequence()
	if err != nil {
		return
	}
	return BigEndian(seq), nil
}

func (db *dbx) ForEach(m Model, f func(m Model)) error {
	err := db.Bolt().View(func(tx *bolt.Tx) error {
		b, err := db.BucketFor(m, tx)
		if err != nil {
			return err
		}
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var vv = m
			if err := json.Unmarshal(v, vv); err != nil {
				return err
			}
			vv.SetKey(k)
			f(vv)
		}
		return nil
	})
	return err
}

func (db *dbx) BucketFor(m Model, tx *bolt.Tx) (*bolt.Bucket, error) {
	name := []byte(getType(m))
	bucket := tx.Bucket(name)
	if bucket == nil {
		return nil, ErrorBucketDoesNotExists(name)
	}
	return bucket, nil
}

func BigEndian(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

func getType(x interface{}) string {
	var t = reflect.TypeOf(x)
	if t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	}
	return t.Name()
}

func checkKey(m Model) error {
	if len(m.Key()) < 1 {
		return ErrorEmptyKey(getType(m))
	}
	return nil
}
//
//func (db *dbx) GetAll(m Model) (data [][]byte, err error) {
//	return data, db.bolt.View(func(tx *bolt.Tx) error {
//		b, err := db.BucketFor(m, tx)
//		if err != nil {
//			return err
//		}
//		c := b.Cursor()
//		for k, v := c.First(); k != nil; k, v = c.Next() {
//			data = append(data, v)
//		}
//		return nil
//	})
//}
