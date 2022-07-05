package source

import (
	"fmt"
	"github.com/auho/go-simple-db/mysql"
	"github.com/auho/go-simple-db/simple"
	"log"
	"os"
	"testing"
)

var mysqlDsn = "test:Test123$@tcp(127.0.0.1:3306)/mysql"
var dbName = "_test_flow"
var tableName = "source"
var driver simple.Driver

func TestMain(m *testing.M) {
	setUp()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func setUp() {
	var err error
	driver, err = simple.NewDriver(mysql.DriverName, mysqlDsn)
	if err != nil {
		log.Fatal(err)
	}

	_, err = driver.Exec(fmt.Sprintf("CREATE DATABASE `%s` DEFAULT CHARACTER SET `utf8mb4` COLLATE `utf8mb4_general_ci`;", dbName))
	if err != nil {
		log.Fatal(err)
	}

	query := "CREATE TABLE IF NOT EXISTS `" + tableName + "` (" +
		"	`id` int(11) unsigned NOT NULL AUTO_INCREMENT," +
		"	`name` varchar(32) NOT NULL DEFAULT ''," +
		"	`value` int(11) NOT NULL DEFAULT '0'," +
		"	PRIMARY KEY (`id`)" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;"
	_, err = driver.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func tearDown() {
	err := driver.Truncate(tableName)
	if err != nil {
		log.Fatal(err)
	}

	_, err = driver.Exec(fmt.Sprintf("DROP DATABASE %s;", dbName))
	if err != nil {
		log.Fatal(err)
	}
}
