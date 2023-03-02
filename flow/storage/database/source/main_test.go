package source

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/auho/go-toolkit/flow/storage/database"
	"gorm.io/driver/mysql"

	"github.com/auho/go-toolkit/flow/storage"
)

var _mysqlDsn = "test:Test123$@tcp(127.0.0.1:3306)/"
var dbName = "_test_flow"
var tableName = "source"
var idName = "id"
var mysqlDsn = _mysqlDsn + dbName

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
	db, err := database.NewDB(mysql.Open(mysqlDsn), nil)
	if err != nil {
		log.Fatal("new db create table ", err)
	}

	err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET `utf8mb4` COLLATE `utf8mb4_general_ci`;", dbName)).Error
	if err != nil {
		log.Fatal("create database ", err)
	}

	query := "CREATE TABLE IF NOT EXISTS `" + dbName + "`.`" + tableName + "` (" +
		"	`id` int(11) unsigned NOT NULL AUTO_INCREMENT," +
		"	`name` varchar(32) NOT NULL DEFAULT ''," +
		"	`value` int(11) NOT NULL DEFAULT '0'," +
		"	PRIMARY KEY (`id`)" +
		") ENGINE=MyISAM DEFAULT CHARSET=utf8mb4;"
	err = db.Exec(query).Error
	if err != nil {
		log.Fatal("create table ", err)
	}
}

func buildData() {
	db, err := database.NewDB(mysql.Open(mysqlDsn), nil)
	if err != nil {
		log.Fatal("new db build data ", err)
	}

	err = db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", tableName)).Error
	if err != nil {
		log.Fatal("build data truncate table ", err)
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

		err = db.Table(tableName).Create(data).Error
		if err != nil {
			log.Fatal("bulk insert ", err, data)
		}
	}

	var count int64
	err = db.Table(tableName).Count(&count).Error
	if err != nil {
		log.Fatal("build data count ", err)
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
	var dbAmount int64
	err = s.db.Table(tableName).Count(&dbAmount).Error
	if err != nil {
		t.Error("db amount ", err)
	}

	if s.total != dbAmount {
		t.Error(fmt.Sprintf("total != db amount %d != %d", s.total, dbAmount))
	}
}

func tearDown() {
	log.Println("tear down")
}
