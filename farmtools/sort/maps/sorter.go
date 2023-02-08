package maps

import (
	"sort"

	sort2 "github.com/auho/go-toolkit/farmtools/sort"
)

func NewSorterByKey[keyE sort2.KeyEntity, valE sort2.ValEntity](m map[keyE]valE, sortedOrder string) *Sorter[keyE, valE] {
	return newSorter[keyE, valE](m, sortedByKey, sortedOrder)
}

func NewSorterByValue[keyE sort2.KeyEntity, valE sort2.ValEntity](m map[keyE]valE, sortedOrder string) *Sorter[keyE, valE] {
	return newSorter[keyE, valE](m, sortedByValue, sortedOrder)
}

func newSorter[keyE sort2.KeyEntity, valE sort2.ValEntity](m map[keyE]valE, sortedBy string, sortedOrder string) *Sorter[keyE, valE] {
	ms := &Sorter[keyE, valE]{}
	ms.sortedBy = sortedBy
	ms.sortedOrder = sortedOrder

	ms.items = make([]Item[keyE, valE], 0, len(m))
	for k, v := range m {
		ms.items = append(ms.items, Item[keyE, valE]{Key: k, Val: v})
	}

	return ms
}

func (s *Sorter[keyE, valE]) Sort() []Item[keyE, valE] {
	sort.Sort(s)
	return s.items
}

func (s *Sorter[keyE, valE]) Len() int {
	return len(s.items)
}

func (s *Sorter[keyE, valE]) Less(i, j int) bool {
	if s.sortedBy == sortedByKey {
		ik := s.items[i].Key
		jk := s.items[j].Key

		if s.sortedOrder == sort2.SortedOrderAsc {
			return ik < jk
		} else {
			return ik > jk
		}
	} else {
		iv := s.items[i].Val
		jv := s.items[j].Val

		if s.sortedOrder == sort2.SortedOrderAsc {
			return iv < jv
		} else {
			return iv > jv
		}
	}
}

func (s *Sorter[keyE, valE]) Swap(i, j int) {
	s.items[i], s.items[j] = s.items[j], s.items[i]
}
