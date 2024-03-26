package describes

import "testing"

func TestDescribeAllTopics(t *testing.T) {
	_, err := DescribeAllTopics(_network, _address)
	if err != nil {
		t.Fatal(err)
	}
}
