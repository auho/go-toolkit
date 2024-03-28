package describes

import (
	"github.com/auho/go-toolkit/elasticsearch/restapi"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func DescribeAllIndices(config elasticsearch.Config) (table.Writer, error) {
	io, err := NewIndicesOverview(config)
	if err != nil {
		return nil, err
	}

	return io.All()
}

type IndicesOverview struct {
	overview
}

func NewIndicesOverview(config elasticsearch.Config) (*IndicesOverview, error) {
	io := &IndicesOverview{}
	err := io.newClient(config)
	if err != nil {
		return nil, err
	}

	return io, nil
}

func (io *IndicesOverview) All() (table.Writer, error) {
	indices, err := io.CatIndices()
	if err != nil {
		return nil, err
	}

	shards, err := io.CatShards()
	if err != nil {
		return nil, err
	}

	indicesShardsMap := make(map[string]IndexShards)

	for _, shard := range shards {
		if _, ok := indicesShardsMap[shard.Index]; !ok {
			indicesShardsMap[shard.Index] = IndexShards{
				Index:    shard.Index,
				ShardsNo: nil,
				Shards:   make(map[string][]Shard),
			}
		}

		_is := indicesShardsMap[shard.Index]
		_is.addShard(shard)

		indicesShardsMap[shard.Index] = _is
	}

	for _, indexShards := range indicesShardsMap {
		indexShards.sort()
	}

	_table := table.NewWriter()
	_table.AppendHeader(table.Row{"index", "status", "health", "pri", "rep", "docs.count", "docs.deleted", "store.size", "pri.store.size", "uuid", "shard", "prirep", "state", "docs", "store", "ip", "node"})
	for _, index := range indices {
		if io.filterSysIndex(index.Index) {
			continue
		}

		var _color text.Color

		switch index.Health {
		case restapi.HealthGreen:
			_color = text.FgGreen
		case restapi.HealthYellow:
			_color = text.FgYellow
		default:
			_color = text.FgRed
		}

		_health := _color.Sprintf("%s", index.Health)

		_table.AppendRow(table.Row{index.Index, index.Status, _health, index.Pri, index.Rep, index.DocsCount, index.DocsDeleted, index.StoreSize, index.PriStoreSize, index.Uuid, "", "", "", "", "", "", ""})
		if indexShards, ok := indicesShardsMap[index.Index]; ok {
			for _, shardNo := range indexShards.ShardsNo {
				if _is, ok1 := indexShards.Shards[shardNo]; ok1 {
					for _, _shard := range _is {
						_table.AppendRow(table.Row{"", "", "", "", "", "", "", "", "", "", _shard.Shard, _shard.Prirep, _shard.State, _shard.Docs, _shard.Store, _shard.Ip, _shard.Node})
					}
				}
			}
		}

		_table.AppendSeparator()
	}

	return _table, nil
}
