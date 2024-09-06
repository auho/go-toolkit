package restapi

import (
	"github.com/auho/go-toolkit/elasticsearch/restapi/entity"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

func (i *Indices) CatIndices() ([]entity.Index, error) {
	var err error
	var indices []entity.Index

	indices, err = DoResponse(
		func() (*esapi.Response, error) {
			return i.client.Cat.Indices(i.client.Cat.Indices.WithFormat("json"))
		},
		indices,
	)
	if err != nil {
		return nil, err
	}

	var newIndices []entity.Index
	for _, _index := range indices {
		if i.FilterSysIndex(_index.Index) {
			continue
		}

		newIndices = append(newIndices, _index)
	}

	return newIndices, nil
}

func (i *Indices) CatShards() ([]entity.Shard, error) {
	var err error
	var shards []entity.Shard

	shards, err = DoResponse(
		func() (*esapi.Response, error) {
			return i.client.Cat.Shards(i.client.Cat.Shards.WithFormat("json"))
		},
		shards,
	)
	if err != nil {
		return nil, err
	}

	return shards, nil
}
