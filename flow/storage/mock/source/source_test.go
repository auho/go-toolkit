package source

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestMock(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	factor := rand.Intn(10) + 1
	total := factor * 100
	pageSize := factor*factor + 1
	s := NewSource(WithTotal(int64(total)), WithPageSize(int64(pageSize)))

	s.Scan()
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
