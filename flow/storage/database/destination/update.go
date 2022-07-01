package destination

import "github.com/auho/go-simple-db/simple"

func usm(idName string) desFunc[MapEntries] {
	return func(sd simple.Driver, tableName string, items MapEntries) error {
		return sd.BulkUpdateFromSliceMapById(tableName, idName, items)
	}
}

func NewUpdateSliceMap(config Config, idName string) (*Destination[MapEntries], error) {
	d := &Destination[MapEntries]{}
	d.desFunc = usm(idName)
	err := d.config(config)

	if err != nil {
		return nil, err
	}

	return d, nil
}
