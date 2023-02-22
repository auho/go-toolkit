package slices

import (
	"sort"

	sort2 "github.com/auho/go-toolkit/farmtools/sort"
)

type sorterStruct[valE sort2.ValEntity, sortE sort2.ValueSorter[valE]] struct {
	sorter[valE]
	origin []sortE
}

func (s *sorterStruct[valE, sortE]) Swap(i, j int) {
	s.items[i], s.items[j] = s.items[j], s.items[i]
	s.origin[i], s.origin[j] = s.origin[j], s.origin[i]
}

func SorterStructAsc[valE sort2.ValEntity, sortE sort2.ValueSorter[valE]](s []sortE) {
	newSorterStruct[valE, sortE](s, sort2.SortedOrderAsc)
}

func SorterStructDesc[valE sort2.ValEntity, sortE sort2.ValueSorter[valE]](s []sortE) {
	newSorterStruct[valE, sortE](s, sort2.SortedOrderDesc)
}

func newSorterStruct[valE sort2.ValEntity, sortE sort2.ValueSorter[valE]](s []sortE, sortedOrder string) {
	ss := &sorterStruct[valE, sortE]{}
	ss.sortedOrder = sortedOrder
	ss.origin = s

	ss.items = make([]valE, 0, len(s))
	for _, v := range s {
		ss.items = append(ss.items, v.SortedVal())
	}

	sort.Sort(ss)
}
