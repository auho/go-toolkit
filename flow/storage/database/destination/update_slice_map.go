package destination

import (
	"github.com/auho/go-simple-db/simple"
	"github.com/auho/go-toolkit/flow/storage"
)

type UpdateSliceMap struct {
	Destination[storage.MapEntries]
	idName string
}

func (u *UpdateSliceMap) Receive(items storage.MapEntries) {
	u.itemsChan <- items
}

func (u *UpdateSliceMap) withDesFunc() desFunc[storage.MapEntries] {
	return func(sd simple.Driver, tableName string, items storage.MapEntries) error {
		return sd.BulkUpdateFromSliceMapById(tableName, u.idName, items)
	}
}

func NewUpdateSliceMap(config Config, idName string) (*UpdateSliceMap, error) {
	u := &UpdateSliceMap{}
	u.idName = idName
	u.desFunc = u.withDesFunc()

	err := u.config(config)
	if err != nil {
		return nil, err
	}

	return u, nil
}
