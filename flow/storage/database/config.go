package database

type Config struct {
	Concurrency int
	Maximum     int64
	StartId     int64
	EndId       int64
	PageSize    int64
	TableName   string
	IdName      string
	Driver      string
	Dsn         string
}

type FromQueryConfig struct {
	Config
	Query string
}

type FromTableConfig struct {
	Config
	Fields []string
}
