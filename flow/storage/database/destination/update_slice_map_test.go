package destination

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/auho/go-simple-db/simple"
	"github.com/auho/go-toolkit/flow/storage"
)

var ussItemsChan = make(chan storage.MapEntries)
var uss *Destination[storage.MapEntry]

func TestUpdateSliceMap(t *testing.T) {
	var err error
	uss, err = NewUpdateSliceMap(Config{
		IsTruncate:  true,
		Concurrency: 4,
		PageSize:    7,
		TableName:   tableName,
		Driver:      driverName,
		Dsn:         mysqlDsn,
	}, idName)

	if err != nil {
		log.Fatal(err)
	}

	rand.Seed(time.Now().UnixNano())
	page := int64(rand.Intn(10)) + 10
	pageSize := int64((rand.Intn(4) + 1) * 10)

	go _buildDataForUpdateSliceMap(t, page, pageSize)

	err = uss.Accept()
	if err != nil {
		log.Fatal(err)
	}

	for items := range ussItemsChan {
		uss.Receive(items)
	}

	uss.Done()

	uss.Finish()

	fmt.Println(uss.Summary())
	fmt.Println(uss.State())

	if uss.state.Amount != page*pageSize {
		t.Error(fmt.Sprintf("actual != expected %d != %d", uss.state.Amount, page*pageSize))
	}

	driver, err := simple.NewDriver(driverName, mysqlDsn)
	if err != nil {
		t.Error(err)
	}

	dbAmountRes, err := driver.QueryFieldInterface("_count", fmt.Sprintf("SELECT COUNT(*) AS `_count` FROM `%s` WHERE `%s` = 2", tableName, valueName))
	if err != nil {
		t.Error("db amount ", err)
	}

	dbAmount, err := strconv.ParseInt(string(dbAmountRes.([]uint8)), 10, 64)
	if err != nil {
		t.Error(fmt.Sprintf("db amount error %v", dbAmountRes))
	}

	if uss.state.Amount != dbAmount {
		t.Error(fmt.Sprintf("total != db amount %d != %d", uss.state.Amount, dbAmount))
	}
}

func _buildDataForUpdateSliceMap(t *testing.T, page, pageSize int64) {
	d, err := simple.NewDriver(driverName, mysqlDsn)
	if err != nil {
		t.Error(err)
	}

	for i := int64(0); i < page; i++ {
		rows := make([][]interface{}, pageSize, pageSize)
		for j := int64(0); j < pageSize; j++ {
			rows[j] = []interface{}{
				fmt.Sprintf("name-%d-%d", i, j),
				1,
			}
		}

		_, err = d.BulkInsertFromSliceSlice(tableName, []string{"name", "value"}, rows)
		if err != nil {
			t.Error(err)
		}
	}

	query := fmt.Sprintf("SELECT `id`, `name`, `value` FROM `%s` WHERE `%s` > ? ORDER BY %s ASC limit ?", tableName, idName, idName)
	for k := int64(0); k < page*pageSize; k += pageSize {
		rows, err := d.QueryInterface(query, k, pageSize)
		if err != nil {
			t.Error(err)
		}

		for index, v := range rows {
			v[valueName] = 2
			rows[index] = v
		}

		ussItemsChan <- rows
	}

	close(ussItemsChan)
}
