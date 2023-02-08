package maps

import (
	"testing"

	"github.com/auho/go-toolkit/farmtools/sort"
)

var _ Comparable[int] = (*sorterStruct)(nil)

type sorterStruct struct {
	b int
}

func (s sorterStruct) SortedVal() int {
	return s.b
}

func Test_SorterStruct(t *testing.T) {
	m := map[string]sorterStruct{"2": {2}, "1": {1}, "3": {3}}

	nm := make(SorterMapStruct[string, int], len(m))
	for k := range m {
		nm[k] = m[k]
	}

	s := NewSorterStructByValue[string](nm, sort.SortedOrderAsc)
	ss := s.Sort()
	if ss[0].Val != 1 {
		t.Error("value asc error")
	}

	s = NewSorterStructByValue[string](nm, sort.SortedOrderDesc)
	ss = s.Sort()
	if ss[0].Val != 3 {
		t.Error("value desc error")
	}

	s = NewSorterStructByKey[string](nm, sort.SortedOrderAsc)
	ss = s.Sort()
	if ss[0].Key != "1" {
		t.Error("key desc error")
	}

	s = NewSorterStructByKey[string](nm, sort.SortedOrderDesc)
	ss = s.Sort()
	if ss[0].Key != "3" {
		t.Error("key desc error")
	}
}
