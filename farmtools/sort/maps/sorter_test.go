package maps

import (
	"testing"
)

func Test_map(t *testing.T) {
	m := map[string]int{"2": 2, "3": 3, "1": 1}

	ss := SorterByValueAsc[string, int](m)
	if ss[0].Val != 1 {
		t.Error("value asc error")
	}

	ss = SorterByValueDesc[string, int](m)
	if ss[0].Val != 3 {
		t.Error("value asc error")
	}

	ss = SorterByKeyAsc[string, int](m)
	if ss[0].Key != "1" {
		t.Error("key asc error")
	}

	ss = SorterByKeyDesc[string, int](m)
	if ss[0].Key != "3" {
		t.Error("key asc error")
	}
}
