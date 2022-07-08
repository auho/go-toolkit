package destination

import (
	"fmt"
	"github.com/auho/go-simple-db/mysql"
	"github.com/auho/go-simple-db/simple"
	"log"
	"os"
	"testing"
)

var driverName = mysql.DriverName
var _mysqlDsn = "test:Test123$@tcp(127.0.0.1:3306)/"
var mainMysqlDsn = _mysqlDsn + "mysql"
var dbName = "_test_flow"
var mysqlDsn = _mysqlDsn + dbName
var tableName = "destination"

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
	driver, err := simple.NewDriver(driverName, mainMysqlDsn)
	if err != nil {
		log.Fatal("new driver create table ", err)
	}

	_, err = driver.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET `utf8mb4` COLLATE `utf8mb4_general_ci`;", dbName))
	if err != nil {
		log.Fatal("create database ", err)
	}

	query := "CREATE TABLE IF NOT EXISTS `" + dbName + "`.`" + tableName + "` (" +
		"	`id` int(11) unsigned NOT NULL AUTO_INCREMENT," +
		"	`name` varchar(32) NOT NULL DEFAULT ''," +
		"	`value` int(11) NOT NULL DEFAULT '0'," +
		"	PRIMARY KEY (`id`)" +
		") ENGINE=MyISAM DEFAULT CHARSET=utf8mb4;"
	_, err = driver.Exec(query)
	if err != nil {
		log.Fatal("create table ", err)
	}
}

func buildData() {
	driver, err := simple.NewDriver(driverName, mysqlDsn)
	if err != nil {
		log.Fatal("new driver build data ", err)
	}

	err = driver.Truncate(tableName)
	if err != nil {
		log.Fatal("build data truncate table ", err)
	}
}

func tearDown() {
	log.Println("tear down")
}
