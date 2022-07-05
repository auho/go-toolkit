package source

import (
	"fmt"
	"testing"
)

func TestSection(t *testing.T) {
	s, err := NewSectionFromTable(
		FromTableConfig{
			Config: Config{
				Concurrency: 4,
				Maximum:     10e4,
				StartId:     0,
				EndId:       0,
				PageSize:    0,
				TableName:   "",
				IdName:      "",
				Driver:      "",
				Dsn:         "",
			},
			Fields: []string{},
		})

	if err != nil {
		t.Error(err)
	}

	s.Scan()

	amount := 0
	for items := range s.ReceiveChan() {
		l := len(items)
		amount = amount + l
	}

	fmt.Println(s.Summary())
	fmt.Println(s.State())

	if s.state.Amount != int64(amount) {
		t.Error(" amount ")
	}
}
