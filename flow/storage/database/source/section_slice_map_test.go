package source

import (
	"testing"

	"github.com/auho/go-simple-db/v2"
	"github.com/auho/go-toolkit/flow/storage"
	"github.com/auho/go-toolkit/flow/storage/database"
)

func TestSectionSliceMapFromTable(t *testing.T) {
	s, err := NewSectionSliceMap(
		&QueryConfig{
			Config: Config{
				Concurrency: 4,
				Maximum:     100000,
				StartId:     0,
				EndId:       100000,
				PageSize:    337,
				TableName:   tableName,
				IdName:      idName,
			},
			Fields: []string{"name", "value"},
		}, func() (*database.DB, error) {
			return database.NewDB(func() (*go_simple_db.SimpleDB, error) {
				return go_simple_db.NewMysql(mysqlDsn)
			})
		})

	if err != nil {
		t.Error(err)
	}

	_testSection[storage.MapEntry](t, s)
}
