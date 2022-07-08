package source

import (
	"fmt"
	"log"
	"strconv"
	"testing"
)

func TestSection(t *testing.T) {
	s, err := NewSectionFromTable(
		FromTableConfig{
			Config: Config{
				Concurrency: 4,
				Maximum:     100000,
				StartId:     0,
				EndId:       100000,
				PageSize:    1000,
				TableName:   tableName,
				IdName:      "id",
				Driver:      driverName,
				Dsn:         mysqlDsn,
			},
			Fields: []string{"name", "value"},
		})

	if err != nil {
		t.Error(err)
	}

	err = s.Scan()
	if err != nil {
		log.Fatal("scan ", err)
	}

	amount := 0
	for items := range s.ReceiveChan() {
		l := len(items)
		amount = amount + l
	}

	fmt.Println(s.Summary())
	fmt.Println(s.State())

	dbAmountRes, err := s.Driver.QueryFieldInterface("_count", fmt.Sprintf("SELECT COUNT(*) AS `_count` FROM `%s`", tableName))
	if err != nil {
		log.Fatal("db amount ", err)
	}

	if s.total != s.state.Amount && s.state.Amount != int64(amount) {
		t.Error(fmt.Sprintf("total != amount != actual %d != %d != %d", s.total, s.state.Amount, amount))
	}

	dbAmount, err := strconv.ParseInt(string(dbAmountRes.([]uint8)), 10, 64)
	if err != nil {
		log.Fatal(fmt.Sprintf("db amount error %v", dbAmountRes))
	}

	if s.total != dbAmount {
		log.Fatal(fmt.Sprintf("total != db amount %d != %d", s.total, dbAmount))
	}
}
