package describes

import (
	"testing"

	"github.com/elastic/go-elasticsearch/v7"
)

func TestDescribeAllIndices(t *testing.T) {
	_, err := DescribeAllIndices(elasticsearch.Config{
		Addresses: []string{_address},
	})
	if err != nil {
		t.Error(err)
	}
}
