package describes

import (
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/jedib0t/go-pretty/v6/table"
)

func DescribeAllShards(config elasticsearch.Config) (table.Writer, error) {
	so, err := NewShardsOverview(config)
	if err != nil {
		return nil, err
	}

	return so.All()
}

type ShardsOverview struct {
	overview
}

func NewShardsOverview(config elasticsearch.Config) (*ShardsOverview, error) {
	so := &ShardsOverview{}
	err := so.newClient(config)
	if err != nil {
		return nil, err
	}

	return so, nil
}

func (so *ShardsOverview) All() (table.Writer, error) {
	shards, err := so.CatShards()
	if err != nil {
		return nil, err
	}

	_table := table.NewWriter()
	_table.AppendHeader(table.Row{"index", "shard", "prirep", "state", "docs", "store", "ip", "node"})
	for _, shard := range shards {
		if so.filterSysIndex(shard.Index) {
			continue
		}

		_table.AppendRow(table.Row{shard.Index, shard.Shard, shard.Prirep, shard.State, shard.Docs, shard.Ip, shard.Node})
		_table.AppendSeparator()
	}

	return _table, nil
}
