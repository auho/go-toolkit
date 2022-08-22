package source

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestSliceMap(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	factor := rand.Intn(10) + 1
	total := factor * 100
	pageSize := factor*factor + 1
	s := NewSliceMap(Config{
		PageSize: int64(pageSize),
		Total:    int64(total),
	})

	err := s.Scan()
	if err != nil {
		t.Fatal(err)
	}

	amount := 0
	for items := range s.ReceiveChan() {
		amount = amount + len(items)
	}

	fmt.Println(s.Summary())
	fmt.Println(s.State())

	if s.amount != int64(amount) {
		t.Error(" amount ")
	}
}
