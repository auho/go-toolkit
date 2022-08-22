package source

import (
	"sync/atomic"
	"time"

	"github.com/auho/go-toolkit/flow/storage"
	"github.com/auho/go-toolkit/flow/tool"
)

var _ mocker[storage.MapEntry] = (*SliceMap)(nil)

type SliceMap struct {
}

func NewSliceMap(config Config) *mock[storage.MapEntry] {
	return newMock[storage.MapEntry](config, &SliceMap{})
}

func (sm SliceMap) scan(idName string, id *int64, amount int64) (*int64, storage.MapEntries) {
	items := make([]map[string]interface{}, amount, amount)
	for i := int64(0); i < amount; i++ {
		item := make(map[string]interface{})
		item[idName] = time.Now().Unix()*1e8 + atomic.AddInt64(id, 1)
		items[i] = item
	}

	return id, items
}

func (sm SliceMap) duplicate(items []storage.MapEntry) []storage.MapEntry {
	return tool.DuplicateSliceMap(items)
}
