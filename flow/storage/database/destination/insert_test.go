package destination

import "testing"

func TestNewInsertSliceSlice(t *testing.T) {
	iss, err := NewInsertSliceSlice(Config{}, []string{})
	if err != nil {
		t.Fatal(err)
	}

	items := [][]interface{}{}

	iss.Receive(items)
}
