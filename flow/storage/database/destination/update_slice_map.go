package destination

import (
	"github.com/auho/go-simple-db/simple"
	"github.com/auho/go-toolkit/flow/storage"
)

var _ storage.Destination[storage.MapEntry] = (*UpdateSliceMap)(nil)

type UpdateSliceMap struct {
	Destination[storage.MapEntry]
	idName string
}

func (u *UpdateSliceMap) Receive(items storage.MapEntries) {
	u.itemsChan <- items
}

func (u *UpdateSliceMap) withDesFunc() desFunc[storage.MapEntry] {
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
