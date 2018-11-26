package gofifo_test

import (
	"reflect"
	"testing"

	"github.com/dwburke/gofifo"
)

func TestSetRoom(t *testing.T) {
	fifo, err := gofifo.NewFifo("test")
	expect(t, err, nil, "")
	defer fifo.Close()

	err = fifo.Push("test string")
	expect(t, err, nil, "")

	obj := fifo.Pop()
	t.Log(obj)

	if obj == nil {
		t.Errorf("Expected object - Got [nil]")
	}
}

func expect(t *testing.T, a interface{}, b interface{}, body string) {
	if a != b {
		t.Errorf("Expected [%v] (type %v) - Got [%v] (type %v) : %s", b, reflect.TypeOf(b), a, reflect.TypeOf(a), body)
	}
}
