package destination

import (
	"github.com/auho/go-toolkit/flow/storage"
)

type InsertSliceSlice struct {
	Destination[storage.SliceEntry]
}

func (i *InsertSliceSlice) Receive(items storage.SliceEntries) {
	i.itemsChan <- items
}

func NewInsertSliceSlice() (*InsertSliceSlice, error) {
	return &InsertSliceSlice{}, nil
}
