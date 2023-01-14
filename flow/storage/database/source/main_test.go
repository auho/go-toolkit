package source

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/auho/go-simple-db/mysql"
	"github.com/auho/go-simple-db/simple"
	"github.com/auho/go-toolkit/flow/storage"
)

var driverName = mysql.DriverName
var _mysqlDsn = "test:Test123$@tcp(127.0.0.1:3306)/"
var mainMysqlDsn = _mysqlDsn + "mysql"
var dbName = "_test_flow"
var mysqlDsn = _mysqlDsn + dbName
var tableName = "source"
var idName = "id"

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

func _testSection[E storage.Entry](
	t *testing.T,
	s *Section[E],
) {
	err := s.Scan()
	if err != nil {
		t.Error("scan ", err)
	}

	amount := 0
	for items := range s.ReceiveChan() {
		l := len(items)
		amount = amount + l
	}

	fmt.Println(s.Summary())
	fmt.Println(s.State())

	if s.total != s.state.Amount() && s.state.Amount() != int64(amount) {
		t.Error(fmt.Sprintf("total != amount != actual %d != %d != %d", s.total, s.state.Amount(), amount))
	}

	dbAmountRes, err := s.GetDriver().QueryFieldInterface("_count", fmt.Sprintf("SELECT COUNT(*) AS `_count` FROM `%s`", tableName))
	if err != nil {
		t.Error("db amount ", err)
	}

	dbAmount, err := strconv.ParseInt(string(dbAmountRes.([]uint8)), 10, 64)
	if err != nil {
		t.Error(fmt.Sprintf("db amount error %v", dbAmountRes))
	}

	if s.total != dbAmount {
		t.Error(fmt.Sprintf("total != db amount %d != %d", s.total, dbAmount))
	}
}

func tearDown() {
	log.Println("tear down")
}