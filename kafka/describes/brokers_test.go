package describes

import "testing"

func TestDescribeBrokers(t *testing.T) {
	_, err := DescribeBrokers(_network, _address)
	if err != nil {
		t.Fatal(err)
	}
}
