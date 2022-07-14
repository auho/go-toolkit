package destination

import (
	"github.com/auho/go-toolkit/flow/storage"
)

var _ storage.Destinationer[storage.MapEntry] = (*InsertSliceMap)(nil)

type InsertSliceMap struct {
	Destination[storage.MapEntry]
}

func NewInsertSliceMap() (*InsertSliceMap, error) {
	return &InsertSliceMap{}, nil
}
