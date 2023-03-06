package destination

import (
	"github.com/auho/go-toolkit/flow/storage"
	"github.com/auho/go-toolkit/flow/storage/database"
)

var _ destinationer[storage.MapEntry] = (*InsertSliceMap)(nil)

type InsertSliceMap struct {
}

func (i *InsertSliceMap) exec(d *Destination[storage.MapEntry], items storage.MapEntries) error {
	return d.db.BulkInsertFromSliceMap(d.table, items, int(d.pageSize))
}

func NewInsertSliceMap(config *Config, b database.BuildDb) (*Destination[storage.MapEntry], error) {
	return newDestination[storage.MapEntry](config, &InsertSliceMap{}, b)
}
