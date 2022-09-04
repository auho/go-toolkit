package destination

import (
	"errors"

	"github.com/auho/go-simple-db/simple"
	"github.com/auho/go-toolkit/flow/storage"
)

var _ destinationer[storage.SliceEntry] = (*InsertSliceSlice)(nil)

type InsertSliceSlice struct {
	fields []string
}

func (i *InsertSliceSlice) desFunc(sd simple.Driver, tableName string, items storage.SliceEntries) error {
	_, err := sd.BulkInsertFromSliceSlice(tableName, i.fields, items)

	return err
}

func NewInsertSliceSlice(config Config, fields []string) (*Destination[storage.SliceEntry], error) {
	if len(fields) <= 0 {
		return nil, errors.New("fields is error")
	}

	iss := &InsertSliceSlice{}
	iss.fields = fields

	return newDestination[storage.SliceEntry](config, iss)
}
