package summary

import (
	"github.com/auho/go-toolkit/elasticsearch/restapi"
	"github.com/elastic/go-elasticsearch/v7"
)

type Summary struct {
	indicesApi *restapi.Indices
}

func NewSummary(config elasticsearch.Config) (*Summary, error) {
	c, err := restapi.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &Summary{indicesApi: restapi.NewIndices(c)}, nil
}

func (s *Summary) AllIndices() ([]Response, error) {
	_indices, err := s.indicesApi.CatIndices()
	if err != nil {
		return nil, err
	}

	var rets []Response
	for _, _index := range _indices {
		ret, err := s.Index(_index.Index)
		if err != nil {
			return nil, err
		}

		rets = append(rets, ret)
	}

	return rets, err
}

func (s *Summary) Index(index string) (Response, error) {
	var resp Response
	resp.Index = index

	_s, err := s.indicesApi.Mapping(index)
	if err != nil {
		return resp, err
	}

	resp.Mapping = _s

	_s, err = s.indicesApi.Settings(index)
	if err != nil {
		return resp, err
	}

	resp.Settings = _s

	return resp, nil
}
