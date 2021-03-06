package gofifo_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/dwburke/gofifo"
)

type Fuu struct {
	Name string
	Time time.Time
}

func TestSetRoom(t *testing.T) {
	fifo, err := gofifo.NewFifo("test")
	expect(t, err, nil, "")
	defer fifo.Close()

	fifo.Register(Fuu{})

	rec := &Fuu{
		Name: "foo",
		Time: time.Now(),
	}
	err = fifo.Enqueue(rec)
	expect(t, err, nil, "")

	var obj = Fuu{}
	err = fifo.Dequeue(&obj)
	expect(t, err, nil, "")

	err = fifo.Dequeue(&obj)
	expect(t, err, gofifo.Empty, "")

	var str string
	err = fifo.Enqueue("this is a string")
	expect(t, err, nil, "")

	err = fifo.Dequeue(&str)
	expect(t, err, nil, "")
	expect(t, str, "this is a string", "")
}

func expect(t *testing.T, a interface{}, b interface{}, body string) {
	if a != b {
		t.Errorf("Expected [%v] (type %v) - Got [%v] (type %v) : %s", b, reflect.TypeOf(b), a, reflect.TypeOf(a), body)
	}
}
