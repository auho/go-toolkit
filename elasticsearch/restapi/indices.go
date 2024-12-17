package restapi

import (
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
)

type Indices struct {
	client *elasticsearch.Client
}

func NewIndices(client *elasticsearch.Client) *Indices {
	return &Indices{
		client: client,
	}
}

func (i *Indices) FilterSysIndex(index string) bool {
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
