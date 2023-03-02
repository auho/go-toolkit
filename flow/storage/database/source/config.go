package source

import "github.com/auho/go-toolkit/flow/storage/database"

type Config struct {
	Concurrency int
	Maximum     int64
	StartId     int64
	EndId       int64
	PageSize    int64
	TableName   string
	IdName      string
}

type QueryConfig struct {
	Config
	Fields []string
	Where  string // "field1 = ? and field2 = ?"
	Order  string // "field1 desc"
}

func (q *QueryConfig) querior(db *database.DB) *database.DB {
	tx := db.Table(q.TableName)
	if len(q.Fields) > 0 {
		tx = tx.Select(q.Fields)
	}

	if q.Where != "" {
		tx = tx.Where(q.Where)
	}

	if q.Order != "" {
		tx = tx.Order(q.Order)
	}

	return &database.DB{DB: tx}
}
