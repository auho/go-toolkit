package describes

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/segmentio/kafka-go"
	"sort"
)

func DescribeBrokers(network, address string) (table.Writer, error) {
	conn, err := kafka.Dial(network, address)
	if err != nil {
		return nil, err
	}

	bs, err := conn.Brokers()
	if err != nil {
		return nil, err
	}

	sort.SliceStable(bs, func(i, j int) bool {
		return bs[i].ID < bs[j].ID
	})

	_table := table.NewWriter()
	_table.AppendHeader(table.Row{"id", "address", "rack"})
	for _, b := range bs {
		_table.AppendRow(table.Row{b.ID, fmt.Sprintf("%s:%d", b.Host, b.Port), b.Rack})
	}

	return _table, nil
}
