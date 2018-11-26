package gofifo

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/streadway/simpleuuid"
	"github.com/syndtr/goleveldb/leveldb"
	//leveldb_errors "github.com/syndtr/goleveldb/leveldb/errors"
	//"container/list"
)

type Fifo struct {
	ldb *leveldb.DB
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

func (fifo *Fifo) GobRegister(t interface{}) {
	gob.Register(t)
}

func (fifo *Fifo) nextId() (string, error) {
	uuid, err := simpleuuid.NewTime(time.Now())
	if err != nil {
		return "", err
	}

	return uuid.String(), nil
}

func (fifo *Fifo) Push(obj interface{}) (err error) {
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

// https://play.golang.org/p/tfjGBp48-ZV
func (fifo *Fifo) Pop() interface{} {
	//data, err := ldb.Get([]byte(key), nil)

	iter := fifo.ldb.NewIterator(nil, nil)
	defer iter.Release()

	var obj interface{}

	for iter.Next() {
		value := iter.Value()
		pCache := bytes.NewBuffer(value)
		decCache := gob.NewDecoder(pCache)
		decCache.Decode(&obj)
		return obj
		break
	}

	return nil
}
