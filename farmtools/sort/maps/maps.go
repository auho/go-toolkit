package maps

import (
	"sort"

	sort2 "github.com/auho/go-toolkit/farmtools/sort"
)

const sortedByKey = "key"
const sortedByValue = "value"

type Item[keyE sort2.KeyEntity, valE sort2.ValEntity] struct {
	Key keyE
	Val valE
}

type Items[keyE sort2.KeyEntity, valE sort2.ValEntity] []Item[keyE, valE]

func (is Items[keyE, valE]) Keys() []keyE {
	keys := make([]keyE, 0, len(is))

	for _, item := range is {
		keys = append(keys, item.Key)
	}

	return keys
}

func (is Items[keyE, valE]) Values() []valE {
	values := make([]valE, 0, len(is))

	for _, item := range is {
		values = append(values, item.Val)
	}

	return values
}

type sorter[keyE sort2.KeyEntity, valE sort2.ValEntity] struct {
	items       Items[keyE, valE]
	sortedBy    string
	sortedOrder string
}

func (s *sorter[keyE, valE]) sort() Items[keyE, valE] {
	sort.Sort(s)
	return s.items
}

func (s *sorter[keyE, valE]) Len() int {
	return len(s.items)
}

func (s *sorter[keyE, valE]) Less(i, j int) bool {
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

func (s *sorter[keyE, valE]) Swap(i, j int) {
	s.items[i], s.items[j] = s.items[j], s.items[i]
}
