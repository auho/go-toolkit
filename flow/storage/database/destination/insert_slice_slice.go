package destination

import (
	"errors"

	"github.com/auho/go-simple-db/simple"
	"github.com/auho/go-toolkit/flow/storage"
)

var _ storage.Destination[storage.SliceEntry] = (*InsertSliceSlice)(nil)

type InsertSliceSlice struct {
	Destination[storage.SliceEntry]
	fields []string
}

func (i *InsertSliceSlice) withDesFunc() desFunc[storage.SliceEntry] {
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
	i.fields = fields
	i.desFunc = i.withDesFunc()

	err := i.config(config)
	if err != nil {
		return nil, err
	}

	return i, nil
}
