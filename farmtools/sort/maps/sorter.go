package maps

import (
	"sort"

	sort2 "github.com/auho/go-toolkit/farmtools/sort"
)

func NewSorterByKey[keyE sort2.KeyEntity, valE sort2.ValEntity](m map[keyE]valE) *Sorter[keyE, valE] {
	return newSorter[keyE, valE](m, sortedByKey)
}

func NewSorterByValue[keyE sort2.KeyEntity, valE sort2.ValEntity](m map[keyE]valE) *Sorter[keyE, valE] {
	return newSorter[keyE, valE](m, sortedByValue)
}

func newSorter[keyE sort2.KeyEntity, valE sort2.ValEntity](m map[keyE]valE, sortedBy string) *Sorter[keyE, valE] {
	ms := &Sorter[keyE, valE]{}
	ms.sortedBy = sortedBy
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
		return s.items[i].Key < s.items[j].Key
	} else {
		return s.items[i].Val < s.items[j].Val
	}
}

func (s *Sorter[keyE, valE]) Swap(i, j int) {
	s.items[i], s.items[j] = s.items[j], s.items[i]
}
