package destination

import (
	"github.com/auho/go-toolkit/flow/storage"
	"github.com/auho/go-toolkit/flow/storage/database"
)

var _ destinationer[storage.MapEntry] = (*UpdateSliceMap)(nil)

type UpdateSliceMap struct {
	idName string
}

func (u *UpdateSliceMap) exec(d *Destination[storage.MapEntry], items storage.MapEntries) error {
	return d.db.BulkUpdateFromSliceMapById(d.table, u.idName, items)
}

func NewUpdateSliceMap(config Config, idName string, b database.BuildDb) (*Destination[storage.MapEntry], error) {
	usm := &UpdateSliceMap{}
	usm.idName = idName

	return newDestination[storage.MapEntry](config, usm, b)
}
