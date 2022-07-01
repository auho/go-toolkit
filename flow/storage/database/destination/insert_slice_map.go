package destination

import (
	"github.com/auho/go-simple-db/simple"
	"github.com/auho/go-toolkit/flow/storage"
)

type InsertSliceMap struct {
	fields []string
	Destination[storage.MapEntries]
}

func (i *InsertSliceMap) Receive(items storage.MapEntries) {
	i.itemsChan <- items
}

func (i *InsertSliceMap) withDesFunc() desFunc[storage.MapEntries] {
	return func(sd simple.Driver, tableName string, items storage.MapEntries) error {
		_, err := sd.BulkInsertFromSliceMap(tableName, items)

		return err
	}
}

func NewInsertSliceMap(config Config) (*InsertSliceMap, error) {
	i := &InsertSliceMap{}
	err := i.config(config)
	if err != nil {
		return nil, err
	}

	i.desFunc = i.withDesFunc()

	return i, nil
}
