package describes

import (
	"testing"

	"github.com/elastic/go-elasticsearch/v7"
)

func TestDescribeAllShards(t *testing.T) {
	_, err := DescribeAllShards(elasticsearch.Config{
		Addresses: []string{_address},
	})
	if err != nil {
		t.Error(err)
	}
}
