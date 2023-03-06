package destination

import (
	"fmt"
	"log"
	"os"
	"testing"

	goSimpleDb "github.com/auho/go-simple-db/v2"
	"github.com/auho/go-toolkit/flow/storage/database"
)

var _mysqlDsn = "test:Test123$@tcp(127.0.0.1:3306)/"
var dbName = "_test_flow"
var mysqlDsn = _mysqlDsn + dbName
var tableName = "destination"
var idName = "id"
var nameName = "name"
var valueName = "value"

func TestMain(m *testing.M) {
	setUp()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func setUp() {
	createTable()
	buildData()
}

func createTable() {
	db, err := database.NewDB(func() (*goSimpleDb.SimpleDB, error) {
		return goSimpleDb.NewMysql(mysqlDsn)
	})

	if err != nil {
		log.Fatal("new db create table ", err)
	}

	err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET `utf8mb4` COLLATE `utf8mb4_general_ci`;", dbName)).Error
	if err != nil {
		log.Fatal("create database ", err)
	}

	query := "CREATE TABLE IF NOT EXISTS `" + dbName + "`.`" + tableName + "` (" +
		"`id` int(11) unsigned NOT NULL AUTO_INCREMENT," +
		"`name` varchar(32) NOT NULL DEFAULT ''," +
		"`value` int(11) NOT NULL DEFAULT '0'," +
		"`created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP," +
		"PRIMARY KEY (`id`)" +
		") ENGINE=MyISAM DEFAULT CHARSET=utf8mb4;"
	err = db.Exec(query).Error
	if err != nil {
		log.Fatal("create table ", err)
	}
}

func buildData() {
	db, err := database.NewDB(func() (*goSimpleDb.SimpleDB, error) {
		return goSimpleDb.NewMysql(mysqlDsn)
	})
	if err != nil {
		log.Fatal("new driver build data ", err)
	}

	err = db.Truncate(tableName)
	if err != nil {
		log.Fatal("build data truncate table ", err)
	}
}

func tearDown() {
	log.Println("tear down")
}
