package maps

import (
	"testing"

	"github.com/auho/go-toolkit/farmtools/sort"
)

func Test_map(t *testing.T) {
	m := map[string]int{"2": 2, "3": 3, "1": 1}

	s := NewSorterByValue[string, int](m, sort.SortedOrderAsc)
	ss := s.Sort()
	if ss[0].Val != 1 {
		t.Error("value asc error")
	}

	s = NewSorterByValue[string, int](m, sort.SortedOrderDesc)
	ss = s.Sort()
	if ss[0].Val != 3 {
		t.Error("value asc error")
	}

	s = NewSorterByKey[string, int](m, sort.SortedOrderAsc)
	ss = s.Sort()
	if ss[0].Key != "1" {
		t.Error("key asc error")
	}

	s = NewSorterByKey[string, int](m, sort.SortedOrderDesc)
	ss = s.Sort()
	if ss[0].Key != "3" {
		t.Error("key asc error")
	}

}
