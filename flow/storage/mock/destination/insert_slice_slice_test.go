package destination

import (
	"fmt"
	"github.com/auho/go-toolkit/flow/storage"
	"math/rand"
	"testing"
	"time"
)

func TestNewInsertSliceSlice(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	page = rand.Intn(49) + 1
	pageSize = (rand.Intn(9) + 1) * pageSize

	var err error
	var dd storage.Destination[storage.SliceEntry]

	dd, err = NewInsertSliceSlice()
	if err != nil {
		t.Error(err)
	}

	d, ok := dd.(*InsertSliceSlice)
	if !ok {
		t.Error("InsertSliceSlice not interface of storage.Destination[storage.SliceEntry]")
	}

	err = d.Accept()
	if err != nil {
		t.Error(err)
	}

	go func() {
		for i := 0; i < page; i++ {
			var sliceSlice [][]interface{}
			for j := 0; j < pageSize; j++ {
				m := make([]interface{}, 0)
				m = append(m, i*page+j)
				sliceSlice = append(sliceSlice, m)
			}

			d.Receive(sliceSlice)
		}

		d.Done()
	}()

	d.Finish()

	fmt.Printf("page: %d, pageSize: %d, amount: %d \n", page, pageSize, page*pageSize)
	fmt.Println(d.Summary())
	fmt.Println(d.State())

	if d.amount != int64(page*pageSize) {
		t.Error(" amount ")
	}
}
