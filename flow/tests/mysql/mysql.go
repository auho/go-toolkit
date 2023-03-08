package mysql

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	goSimpleDb "github.com/auho/go-simple-db/v2"
	"github.com/auho/go-toolkit/flow/storage/database"
)

var _dsn = "test:Test123$@tcp(127.0.0.1:3306)/"
var dbName = "_test_flow"

var SourceTable = "source"
var DestinationTable = "destination"
var IdName = "id"
var NameName = "name"
var ValueName = "value"
var Dsn = _dsn + dbName
var DB *database.DB

func init() {
	var err error
	DB, err = database.NewDB(func() (*goSimpleDb.SimpleDB, error) {
		return goSimpleDb.NewMysql(Dsn)
	})
	if err != nil {
		log.Fatal("new DB create table ", err)
	}
}

func CreateTable(table string) {
	err := DB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET `utf8mb4` COLLATE `utf8mb4_general_ci`;", dbName)).Error
	if err != nil {
		log.Fatal("create database ", err)
	}

	query := "CREATE TABLE IF NOT EXISTS `" + dbName + "`.`" + table + "` (" +
		"`id` int(11) unsigned NOT NULL AUTO_INCREMENT," +
		"`name` varchar(32) NOT NULL DEFAULT ''," +
		"`value` int(11) NOT NULL DEFAULT '0'," +
		"`created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP," +
		"PRIMARY KEY (`id`)" +
		") ENGINE=MyISAM DEFAULT CHARSET=utf8mb4;"
	err = DB.Exec(query).Error
	if err != nil {
		log.Fatal("create table ", err)
	}
}

func BuildData(table string) {
	err := DB.Exec(fmt.Sprintf("TRUNCATE TABLE %s", table)).Error
	if err != nil {
		log.Fatal("build data", err)
	}

	rand.Seed(time.Now().UnixNano())
	page := int64(rand.Intn(10)) + 10
	pageSize := int64((rand.Intn(4) + 1) * 1000)

	for i := int64(0); i < page; i++ {
		data := make([]map[string]any, pageSize, pageSize)
		for j := int64(0); j < pageSize; j++ {
			data[j] = map[string]any{
				"name":  fmt.Sprintf("name-%d-%d", i, j),
				"value": i * j,
			}
		}

		err = DB.Table(table).Create(data).Error
		if err != nil {
			log.Fatal("bulk insert ", err, data)
		}
	}

	var count int64
	err = DB.Table(table).Count(&count).Error
	if err != nil {
		log.Fatal("build data count ", err)
	}

	if count != page*pageSize {
		log.Fatal(fmt.Sprintf("build data bulk insert actual != expected [%d] != [%d]", count, pageSize*page))
	}
}

func CleanData(table string) {
	err := DB.Truncate(table)
	if err != nil {
		log.Fatal("clean data", err)
	}
}
