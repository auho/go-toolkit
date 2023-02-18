package slices

import (
	"sort"

	sort2 "github.com/auho/go-toolkit/farmtools/sort"
)

type sorterStruct[valE sort2.ValEntity] struct {
	sorter[valE]
	origin []sort2.ValueSorter[valE]
}

func (s *sorterStruct[valE]) Swap(i, j int) {
	s.items[i], s.items[j] = s.items[j], s.items[i]
	s.origin[i], s.origin[j] = s.origin[j], s.origin[i]
}

func SorterStructAsc[valE sort2.ValEntity](s []sort2.ValueSorter[valE]) {
	newSorterStruct(s, sort2.SortedOrderAsc)
}

func SorterStructDesc[valE sort2.ValEntity](s []sort2.ValueSorter[valE]) {
	newSorterStruct(s, sort2.SortedOrderDesc)
}

func newSorterStruct[valE sort2.ValEntity](s []sort2.ValueSorter[valE], sortedOrder string) {
	ss := &sorter[valE]{}
	ss.sortedOrder = sortedOrder

	ss.items = make([]valE, 0, len(s))
	for _, v := range s {
		ss.items = append(ss.items, v.SortedVal())
	}

	sort.Sort(ss)
}
