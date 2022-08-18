package source

import (
	"fmt"
	"strings"

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
	return tool.DuplicateSliceMap(items)
}

func NewSectionSliceMapFromQuery(config FromQueryConfig) (*Section[storage.MapEntry], error) {
	s, err := newSectionSliceMap(config.Config)
	if err != nil {
		return nil, err
	}

	s.query = config.Query

	return s, nil
}

func NewSectionSliceMapFromTable(config FromTableConfig) (*Section[storage.MapEntry], error) {
	s, err := newSectionSliceMap(config.Config)
	if err != nil {
		return nil, err
	}

	s.fields = config.Fields

	fieldsSting := fmt.Sprintf("`%s`", strings.Join(s.fields, "`,`"))
	s.query = fmt.Sprintf("SELECT %s FROM `%s` WHERE `%s` > ? ORDER BY `%s` ASC limit ?", fieldsSting, s.tableName, s.idName, s.idName)

	return s, nil
}

func newSectionSliceMap(config Config) (*Section[storage.MapEntry], error) {
	s := &Section[storage.MapEntry]{}
	err := s.config(config)
	if err != nil {
		return nil, err
	}

	s.sectioner = &SectionSliceMap{}

	return s, nil
}
