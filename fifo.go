package gofifo

import (
	"bytes"
	"encoding/gob"
	"errors"
	"sync"
	"time"

	"github.com/streadway/simpleuuid"
	"github.com/syndtr/goleveldb/leveldb"
)

var Empty = errors.New("fifo: empty")

type Fifo struct {
	ldb *leveldb.DB
	mux sync.Mutex
}

func NewFifo(datadir string) (*Fifo, error) {
	fifo := &Fifo{}

	var err error
	fifo.ldb, err = leveldb.OpenFile(datadir, nil)
	if err != nil {
		return nil, err
	}

	return fifo, nil
}

func (fifo *Fifo) Close() {
	if fifo.ldb != nil {
		fifo.ldb.Close()
		fifo.ldb = nil
	}
}

func (fifo *Fifo) Register(t interface{}) {
	gob.Register(t)
}

func (fifo *Fifo) nextId() (string, error) {
	uuid, err := simpleuuid.NewTime(time.Now())
	if err != nil {
		return "", err
	}

	return uuid.String(), nil
}

func (fifo *Fifo) Enqueue(obj interface{}) (err error) {
	uuid, err := fifo.nextId()
	if err != nil {
		return
	}

	mCache := new(bytes.Buffer)
	encCache := gob.NewEncoder(mCache)
	encCache.Encode(obj)

	err = fifo.ldb.Put([]byte(uuid), mCache.Bytes(), nil)
	return
}

func (fifo *Fifo) Dequeue(obj interface{}) (err error) {
	fifo.mux.Lock()
	defer fifo.mux.Unlock()

	err = Empty

	iter := fifo.ldb.NewIterator(nil, nil)
	defer iter.Release()

	for iter.Next() {
		var key = string(iter.Key())
		value := iter.Value()
		pCache := bytes.NewBuffer(value)
		decCache := gob.NewDecoder(pCache)

		err = decCache.Decode(obj)
		if err != nil {
			return
		}

		fifo.ldb.Delete([]byte(key), nil)
		return
	}

	return
}
