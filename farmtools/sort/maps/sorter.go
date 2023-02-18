package maps

import (
	sort2 "github.com/auho/go-toolkit/farmtools/sort"
)

func SorterByKeyAsc[keyE sort2.KeyEntity, valE sort2.ValEntity](m map[keyE]valE) []Item[keyE, valE] {
	return newSorter[keyE, valE](m, sortedByKey, sort2.SortedOrderAsc)
}

func SorterByKeyDesc[keyE sort2.KeyEntity, valE sort2.ValEntity](m map[keyE]valE) []Item[keyE, valE] {
	return newSorter[keyE, valE](m, sortedByKey, sort2.SortedOrderDesc)
}

func SorterByValueAsc[keyE sort2.KeyEntity, valE sort2.ValEntity](m map[keyE]valE) []Item[keyE, valE] {
	return newSorter[keyE, valE](m, sortedByValue, sort2.SortedOrderAsc)
}

func SorterByValueDesc[keyE sort2.KeyEntity, valE sort2.ValEntity](m map[keyE]valE) []Item[keyE, valE] {
	return newSorter[keyE, valE](m, sortedByValue, sort2.SortedOrderDesc)
}

func newSorter[keyE sort2.KeyEntity, valE sort2.ValEntity](m map[keyE]valE, sortedBy string, sortedOrder string) []Item[keyE, valE] {
	ms := &sorter[keyE, valE]{}
	ms.sortedBy = sortedBy
	ms.sortedOrder = sortedOrder

	ms.items = make([]Item[keyE, valE], 0, len(m))
	for k, v := range m {
		ms.items = append(ms.items, Item[keyE, valE]{Key: k, Val: v})
	}

	return ms.sort()
}
