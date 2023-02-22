package slices

import (
	"testing"

	"github.com/auho/go-toolkit/farmtools/sort"
)

var _ sort.ValueSorter[int] = (*sortedStruct)(nil)

type sortedStruct struct {
	b int
}

func (s sortedStruct) SortedVal() int {
	return s.b
}

func Test_sorterStruct(t *testing.T) {
	s := []sortedStruct{{2}, {1}, {3}}

	SorterStructAsc[int](s)
	if s[0].SortedVal() != 1 {
		t.Error("asc error")
	}

	SorterStructDesc[int](s)
	if s[0].SortedVal() != 3 {
		t.Error("desc error")
	}
}
