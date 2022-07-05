package destination

import (
	"github.com/auho/go-toolkit/flow/storage"
)

type InsertSliceMap struct {
	Destination[storage.MapEntry]
}

func (i *InsertSliceMap) Receive(items storage.MapEntries) {
	i.itemsChan <- items
}

func NewInsertSliceMap() (*InsertSliceMap, error) {
	return &InsertSliceMap{}, nil
}
