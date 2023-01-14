package destination

import (
	"github.com/auho/go-simple-db/simple"
	"github.com/auho/go-toolkit/flow/storage"
)

var _ destinationer[storage.MapEntry] = (*UpdateSliceMap)(nil)

type UpdateSliceMap struct {
	idName string
}

func (u *UpdateSliceMap) desFunc(sd simple.Driver, tableName string, items storage.MapEntries) error {
	return sd.BulkUpdateFromSliceMapById(tableName, u.idName, items)
}

func NewUpdateSliceMap(config Config, idName string) (*Destination[storage.MapEntry], error) {
	usm := &UpdateSliceMap{}
	usm.idName = idName

	return newDestination[storage.MapEntry](config, usm)
}
