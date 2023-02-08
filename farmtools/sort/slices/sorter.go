package slices

import (
	"sort"

	sort2 "github.com/auho/go-toolkit/farmtools/sort"
)

var _ sort.Interface = (*sorter[string])(nil)

type Comparable[valE sort2.ValEntity] interface {
	SortedVal() valE
}

type sorter[valE sort2.ValEntity] struct {
	items       []valE
	sortedOrder string
}

func (s *sorter[valE]) Len() int {
	return len(s.items)
}

func (s *sorter[valE]) Less(i, j int) bool {
	iv := s.items[i]
	jv := s.items[j]

	if s.sortedOrder == sort2.SortedOrderAsc {
		return iv < jv
	} else {
		return iv > jv
	}
}

func (s *sorter[valE]) Swap(i, j int) {
	s.items[i], s.items[j] = s.items[j], s.items[i]
}

func SorterAsc[valE sort2.ValEntity](s []valE) {
	newSorter(s, sort2.SortedOrderAsc)
}

func SorterDesc[valE sort2.ValEntity](s []valE) {
	newSorter(s, sort2.SortedOrderDesc)
}

func newSorter[valE sort2.ValEntity](s []valE, sortedOrder string) {
	ss := &sorter[valE]{}
	ss.sortedOrder = sortedOrder
	ss.items = s

	sort.Sort(ss)
}
