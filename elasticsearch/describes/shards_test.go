package describes

import (
	"github.com/elastic/go-elasticsearch/v7"
	"testing"
)

func TestDescribeAllShards(t *testing.T) {
	_, err := DescribeAllShards(elasticsearch.Config{
		Addresses: []string{_address},
	})
	if err != nil {
		t.Error(err)
	}
}
