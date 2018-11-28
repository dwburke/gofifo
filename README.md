# gofifo

FiFo queue backed by leveldb, using Go's gob format for data encoding.

Examples to follow...


```
import (
	"fmt"
	ldb_errors "github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/dwburke/gofifo"
)

type Foo struct {
	Name string
	Time time.Time
}


func main() {
	fifo, err := gofifo.NewFifo("test")
	if err != nil {
		panic(err)
	}
	defer fifo.Close()

	fifo.Register(Foo{})

	rec := &Foo{
		Name: "foo",
		Time: time.Now(),
	}
	err = fifo.Enqueue(rec)
	if err != nil {
		panic(err)
	}

	while (1) {
		var obj = Fuu{}
		err = fifo.Dequeue(&obj)
		if err != nil {
			if err == gofifo.Empty {
				break
			}
			panic(err)
		}
	}

}
```
