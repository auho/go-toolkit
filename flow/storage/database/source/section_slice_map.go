package source

import (
	"github.com/auho/go-simple-db/simple"
	"github.com/auho/go-toolkit/flow/storage"
	"github.com/auho/go-toolkit/flow/tool"
)

var _ sectioner[storage.MapEntry] = (*SectionSliceMap)(nil)

type SectionSliceMap struct {
}

func (ssm *SectionSliceMap) sourceFunc(driver simple.Driver, query string, startId, size int64) ([]storage.MapEntry, error) {
	return driver.QueryInterface(query, startId, size)
}

func (ssm *SectionSliceMap) duplicate(items storage.MapEntries) storage.MapEntries {
	return tool.DuplicateSliceMap[tool.InterfaceEntry](items)
}

func NewSectionSliceMapFromQuery(config FromQueryConfig) (*Section[storage.MapEntry], error) {
	s, err := newSection[storage.MapEntry](withConfigFromQuery[storage.MapEntry](config), withSectioner[storage.MapEntry](&SectionSliceMap{}))
	if err != nil {
		return nil, err
	}

	return s, nil
}

func NewSectionSliceMapFromTable(config FromTableConfig) (*Section[storage.MapEntry], error) {
	s, err := newSection[storage.MapEntry](withConfigFromTable[storage.MapEntry](config), withSectioner[storage.MapEntry](&SectionSliceMap{}))
	if err != nil {
		return nil, err
	}

	return s, nil
}
