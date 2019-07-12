package storage

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/boltdb/bolt"

	log "github.com/cihub/seelog"
)

// BoltDbStorage bolt storage struct
type BoltDbStorage struct {
	Db         *bolt.DB
	bucketName string
	contents   sync.Map
	count      int32
}

var (
	fileName   = "proxy.db"
	bucketName = "IpList"

	//Storage db
	Storage *BoltDbStorage
)

func init() {
	db, err := bolt.Open(fileName, 0600, nil)
	if err != nil {
		log.Critical("Blot DB init failed,", err)
	}

	Storage = &BoltDbStorage{
		Db:         db,
		bucketName: bucketName,
	}

	log.Info("storage init")
	// Sync data from database to memory.
	Storage.sync()
}

func (s *BoltDbStorage) sync() {
	s.Db.View(func(tx *bolt.Tx) error {
		tx.Bucket([]byte(s.bucketName)).ForEach(func(k, v []byte) error {
			key, value := make([]byte, len(k)), make([]byte, len(v))
			copy(key, k)
			copy(value, v)
			s.contents.Store(string(key), value)
			atomic.AddInt32(&s.count, 1)
			return nil
		})

		return nil
	})

	//log.Debugf("content:%v,count:%d", s.contents, s.count)
}

// Close will close the DB.
func (s *BoltDbStorage) Close() {
	s.Db.Close()
}

// GetRandomOne Get one random record.
func (s *BoltDbStorage) GetRandomOne() []byte {
	if s.count == 0 {
		return nil
	}

	var randomKey string
	var defaultKey string
	index := rand.New(rand.NewSource(time.Now().Unix())).Intn(int(atomic.LoadInt32(&s.count)))

	s.contents.Range(func(key, value interface{}) bool {
		// Set a default key to avoid that other goroutine is deleting content at the same time.
		if defaultKey == "" {
			defaultKey, _ = key.(string)
		}

		if index == 0 {
			randomKey, _ = key.(string)
			return false
		}

		index--
		return true
	})

	if randomKey == "" {
		randomKey = defaultKey
	}

	return s.Get(randomKey)
}

// Get will get the json byte value of key.
func (s *BoltDbStorage) Get(key string) []byte {
	var value []byte

	if temp, ok := s.contents.Load(key); ok {
		if content, ok := temp.([]byte); ok {
			value = append(value, content...)
		}
	}

	return value
}
