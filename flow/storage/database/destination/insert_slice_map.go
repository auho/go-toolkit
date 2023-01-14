package destination

import (
	"github.com/auho/go-simple-db/simple"
	"github.com/auho/go-toolkit/flow/storage"
)

var _ destinationer[storage.MapEntry] = (*InsertSliceMap)(nil)

type InsertSliceMap struct {
}

func (i *InsertSliceMap) desFunc(sd simple.Driver, tableName string, items storage.MapEntries) error {
	_, err := sd.BulkInsertFromSliceMap(tableName, items)

	return err
}

func NewInsertSliceMap(config Config) (*Destination[storage.MapEntry], error) {
	return newDestination[storage.MapEntry](config, &InsertSliceMap{})
}
