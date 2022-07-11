package destination

import (
	"github.com/auho/go-simple-db/simple"
	"github.com/auho/go-toolkit/flow/storage"
)

var _ storage.Destination[storage.MapEntry] = (*InsertSliceMap)(nil)

type InsertSliceMap struct {
	Destination[storage.MapEntry]
}

func (i *InsertSliceMap) withDesFunc() desFunc[storage.MapEntry] {
	return func(sd simple.Driver, tableName string, items storage.MapEntries) error {
		_, err := sd.BulkInsertFromSliceMap(tableName, items)

		return err
	}
}

func NewInsertSliceMap(config Config) (*InsertSliceMap, error) {
	i := &InsertSliceMap{}
	i.desFunc = i.withDesFunc()

	err := i.config(config)
	if err != nil {
		return nil, err
	}

	return i, nil
}
