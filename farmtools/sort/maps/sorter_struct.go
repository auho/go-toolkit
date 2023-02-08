package maps

import (
	sort2 "github.com/auho/go-toolkit/farmtools/sort"
)

type SorterMapStruct[keyE sort2.KeyEntity, valE sort2.ValEntity] map[keyE]Comparable[valE]

func NewSorterStructByKey[keyE sort2.KeyEntity, valE sort2.ValEntity](
	m SorterMapStruct[keyE, valE], sortedOrder string,
) *Sorter[keyE, valE] {
	return newSorterStruct[keyE, valE](m, sortedByKey, sortedOrder)
}

func NewSorterStructByValue[keyE sort2.KeyEntity, valE sort2.ValEntity](
	m SorterMapStruct[keyE, valE], sortedOrder string,
) *Sorter[keyE, valE] {
	return newSorterStruct[keyE, valE](m, sortedByValue, sortedOrder)
}

func newSorterStruct[keyE sort2.KeyEntity, valE sort2.ValEntity](
	m SorterMapStruct[keyE, valE], sortedBy string, sortedOrder string,
) *Sorter[keyE, valE] {
	ms := &Sorter[keyE, valE]{}
	ms.sortedBy = sortedBy
	ms.sortedOrder = sortedOrder

	ms.items = make([]Item[keyE, valE], 0, len(m))
	for k, v := range m {
		ms.items = append(ms.items, Item[keyE, valE]{Key: k, Val: v.SortedVal()})
	}

	return ms
}
