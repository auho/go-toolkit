package destination

import (
	"github.com/auho/go-simple-db/simple"
	"github.com/auho/go-toolkit/flow/storage"
)

func usm(idName string) desFunc[storage.MapEntries] {
	return func(sd simple.Driver, tableName string, items storage.MapEntries) error {
		return sd.BulkUpdateFromSliceMapById(tableName, idName, items)
	}
}

func NewUpdateSliceMap(config Config, idName string) (*Destination[storage.MapEntries], error) {
	d := &Destination[storage.MapEntries]{}
	d.desFunc = usm(idName)
	err := d.config(config)

	if err != nil {
		return nil, err
	}

	return d, nil
}
