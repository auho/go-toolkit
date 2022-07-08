package source

import (
	"fmt"
	"github.com/auho/go-simple-db/mysql"
	"github.com/auho/go-simple-db/simple"
	"log"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"
)

var driverName = mysql.DriverName
var _mysqlDsn = "test:Test123$@tcp(127.0.0.1:3306)/"
var mainMysqlDsn = _mysqlDsn + "mysql"
var dbName = "_test_flow"
var mysqlDsn = _mysqlDsn + dbName
var tableName = "source"

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

	rand.Seed(time.Now().UnixNano())
	page := int64(rand.Intn(10)) + 10
	pageSize := int64((rand.Intn(4) + 1) * 1000)

	for i := int64(0); i < page; i++ {
		data := make([][]interface{}, pageSize, pageSize)
		for j := int64(0); j < pageSize; j++ {
			data[j] = []interface{}{
				fmt.Sprintf("name-%d-%d", i, j),
				i * j,
			}
		}

		_, err = driver.BulkInsertFromSliceSlice(tableName, []string{"name", "value"}, data)
		if err != nil {
			log.Fatal("bulk insert ", err, data)
		}
	}

	countRes, err := driver.QueryFieldInterface("_count", fmt.Sprintf("SELECT COUNT(*) AS `_count` FROM `%s`", tableName))
	if err != nil {
		log.Fatal("build data count ", err)
	}

	count, err := strconv.ParseInt(string(countRes.([]uint8)), 10, 64)
	if err != nil {
		log.Fatal(fmt.Sprintf("build data count %v", countRes))
	}

	if count != page*pageSize {
		log.Fatal(fmt.Sprintf("build data bulk insert actual != expected [%d] != [%d]", count, pageSize*page))
	}
}

func tearDown() {
	//driver, err := simple.NewDriver(driverName, mysqlDsn)
	//if err != nil {
	//	log.Fatal("new driver tearDown ", err)
	//}
	//
	//err = driver.Truncate(tableName)
	//if err != nil {
	//	log.Fatal("truncate table ", err)
	//}

	//_, err = driver.Exec(fmt.Sprintf("DROP DATABASE %s;", dbName))
	//if err != nil {
	//	log.Fatal("drop database ", err)
	//}
	//
	log.Println("tear down")
}
