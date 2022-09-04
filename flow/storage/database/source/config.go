package source

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
	// SELECT `id` FROM `table` WHERE `id` > ? ORDER BY `id` ASC limit ?
	Query string
}

type FromTableConfig struct {
	Config
	Fields []string
}
