package source

import (
	"fmt"
	"testing"

	"github.com/auho/go-toolkit/flow/storage"
)

func TestSectionSliceMapFromTable(t *testing.T) {
	s, err := NewSectionSliceMapFromTable(
		FromTableConfig{
			Config: Config{
				Concurrency: 4,
				Maximum:     100000,
				StartId:     0,
				EndId:       100000,
				PageSize:    337,
				TableName:   tableName,
				IdName:      idName,
				Driver:      driverName,
				Dsn:         mysqlDsn,
			},
			Fields: []string{"name", "value"},
		})

	if err != nil {
		t.Error(err)
	}

	_testSection[storage.MapEntry](t, s)
}

func TestSectionSliceMapFromQuery(t *testing.T) {
	s, err := NewSectionSliceMapFromQuery(
		FromQueryConfig{
			Config: Config{
				Concurrency: 4,
				Maximum:     100000,
				StartId:     0,
				EndId:       100000,
				PageSize:    223,
				TableName:   tableName,
				IdName:      idName,
				Driver:      driverName,
				Dsn:         mysqlDsn,
			},
			Query: fmt.Sprintf("SELECT %s FROM `%s` WHERE `%s` > ? ORDER BY `%s` ASC limit ?",
				"`id`, `name`, `value`", tableName, idName, idName),
		})

	if err != nil {
		t.Error(err)
	}

	_testSection[storage.MapEntry](t, s)
}
