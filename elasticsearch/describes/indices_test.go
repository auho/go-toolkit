package describes

import (
	"github.com/elastic/go-elasticsearch/v7"
	"testing"
)

func TestDescribeAllIndices(t *testing.T) {
	_, err := DescribeAllIndices(elasticsearch.Config{
		Addresses: []string{_address},
	})
	if err != nil {
		t.Error(err)
	}
}
