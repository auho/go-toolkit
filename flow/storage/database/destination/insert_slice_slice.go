package destination

import (
	"errors"
	"github.com/auho/go-simple-db/simple"
	"github.com/auho/go-toolkit/flow/storage"
)

type InsertSliceSlice struct {
	fields []string
	Destination[storage.SliceEntries]
}

func (i *InsertSliceSlice) Receive(items storage.SliceEntries) {
	i.itemsChan <- items
}

func (i *InsertSliceSlice) withDesFunc() desFunc[storage.SliceEntries] {
	return func(sd simple.Driver, tableName string, items storage.SliceEntries) error {
		_, err := sd.BulkInsertFromSliceSlice(tableName, i.fields, items)

		return err
	}
}

func NewInsertSliceSlice(config Config, fields []string) (*InsertSliceSlice, error) {
	if len(fields) <= 0 {
		return nil, errors.New("fields is error")
	}

	i := &InsertSliceSlice{}
	err := i.config(config)
	if err != nil {
		return nil, err
	}

	i.desFunc = i.withDesFunc()

	return i, nil
}
