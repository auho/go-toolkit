package destination

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/auho/go-toolkit/flow/storage/database"
	"github.com/auho/go-toolkit/flow/tests/mysql"
)

func TestInsertSliceSlice(t *testing.T) {
	iss, err := NewInsertSliceSlice(
		&Config{
			IsTruncate:  true,
			Concurrency: 4,
			PageSize:    337,
			TableName:   tableName,
		},
		[]string{nameName, valueName},
		func() (*database.DB, error) {
			return mysql.DB, nil
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	err = iss.Accept()
	if err != nil {
		log.Fatal(err)
	}

	rand.Seed(time.Now().UnixNano())
	page := int64(rand.Intn(10)) + 10
	pageSize := int64((rand.Intn(4) + 1) * 1000)

	go func() {
		for i := int64(0); i < page; i++ {
			data := make([][]any, pageSize, pageSize)
			for j := int64(0); j < pageSize; j++ {
				data[j] = []any{
					fmt.Sprintf("name-%d-%d", i, j),
					i * j,
				}
			}

			iss.Receive(data)
		}

		iss.Done()
	}()

	iss.Finish()

	fmt.Println(iss.Summary())
	fmt.Println(iss.State())

	if iss.state.Amount() != page*pageSize {
		t.Error(fmt.Sprintf("actual != expected %d != %d", iss.state.Amount(), page*pageSize))
	}
	var dbAmount int64
	iss.db.Table(tableName).Count(&dbAmount)
	if err != nil {
		t.Error("db amount ", err)
	}

	if iss.state.Amount() != dbAmount {
		t.Error(fmt.Sprintf("total != db amount %d != %d", iss.state.Amount(), dbAmount))
	}
}
