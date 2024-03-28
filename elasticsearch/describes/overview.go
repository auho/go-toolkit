package describes

import (
	"github.com/auho/go-toolkit/elasticsearch/restapi"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"sort"
	"strings"
)

type Index struct {
	DocsCount    string `json:"docs.count"`
	DocsDeleted  string `json:"docs.deleted"`
	Health       string `json:"health"`
	Index        string `json:"index"`
	Pri          string `json:"pri"`
	PriStoreSize string `json:"pri.store.size"`
	Rep          string `json:"rep"`
	Status       string `json:"status"`
	StoreSize    string `json:"store.size"`
	Uuid         string `json:"uuid"`
	Shards       []Shard
}

type Shard struct {
	Index  string `json:"index"`
	Shard  string `json:"shard"`
	Prirep string `json:"prirep"`
	State  string `json:"state"`
	Docs   string `json:"docs"`
	Store  string `json:"store"`
	Ip     string `json:"ip"`
	Node   string `json:"node"`
}
type IndicesShards []IndexShards

type IndexShards struct {
	Index    string
	ShardsNo []string           // []shard no
	Shards   map[string][]Shard // map[shard no][]Shard
}

func (is *IndexShards) addShard(shard Shard) {
	if _, ok := is.Shards[shard.Shard]; ok {
		is.Shards[shard.Shard] = append(is.Shards[shard.Shard], shard)
	} else {
		is.ShardsNo = append(is.ShardsNo, shard.Shard)
		is.Shards[shard.Shard] = append(is.Shards[shard.Shard], shard)
	}
}

func (is *IndexShards) hasShardNo(shardNo string) bool {
	_, ok := is.Shards[shardNo]

	return ok
}

func (is *IndexShards) sort() {
	sort.SliceStable(is.ShardsNo, func(i, j int) bool {
		return is.ShardsNo[i] < is.ShardsNo[j]
	})

	for _, shards := range is.Shards {
		sort.SliceStable(shards, func(i, j int) bool {
			return shards[i].Prirep < shards[j].Prirep
		})
	}
}

type overview struct {
	client *elasticsearch.Client
}

func (o *overview) newClient(config elasticsearch.Config) error {
	var err error
	o.client, err = elasticsearch.NewClient(config)
	if err != nil {
		return err
	}
	return nil
}

func (o *overview) filterSysIndex(index string) bool {
	prefixFilters := []string{".monitoring", ".reporting", ".kibana"}

	var ok bool
	for _, filter := range prefixFilters {
		if strings.HasPrefix(index, filter) {
			ok = true
			break
		}
	}

	return ok
}

func (o *overview) CatShards() ([]Shard, error) {
	var err error
	var shards []Shard

	shards, err = restapi.Do(
		func() (*esapi.Response, error) {
			return o.client.Cat.Shards(o.client.Cat.Shards.WithFormat("json"))
		},
		shards,
	)
	if err != nil {
		return nil, err
	}

	return shards, nil
}

func (o *overview) CatIndices() ([]Index, error) {
	var err error
	var indices []Index

	indices, err = restapi.Do(
		func() (*esapi.Response, error) {
			return o.client.Cat.Indices(o.client.Cat.Indices.WithFormat("json"))
		},
		indices,
	)
	if err != nil {
		return nil, err
	}

	return indices, nil
}
