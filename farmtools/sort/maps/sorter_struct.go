package maps

import (
	sort2 "github.com/auho/go-toolkit/farmtools/sort"
)

type SorterMapStruct[keyE sort2.KeyEntity, valE sort2.ValEntity] map[keyE]sort2.ValueSorter[valE]

func SorterStructByKeyAsc[keyE sort2.KeyEntity, valE sort2.ValEntity](m SorterMapStruct[keyE, valE]) Items[keyE, valE] {
	return newSorterStruct[keyE, valE](m, sortedByKey, sort2.SortedOrderAsc)
}

func SorterStructByKeyDesc[keyE sort2.KeyEntity, valE sort2.ValEntity](m SorterMapStruct[keyE, valE]) Items[keyE, valE] {
	return newSorterStruct[keyE, valE](m, sortedByKey, sort2.SortedOrderDesc)
}

func SorterStructByValueAsc[keyE sort2.KeyEntity, valE sort2.ValEntity](m SorterMapStruct[keyE, valE]) Items[keyE, valE] {
	return newSorterStruct[keyE, valE](m, sortedByValue, sort2.SortedOrderAsc)
}

func SorterStructByValueDesc[keyE sort2.KeyEntity, valE sort2.ValEntity](m SorterMapStruct[keyE, valE]) Items[keyE, valE] {
	return newSorterStruct[keyE, valE](m, sortedByValue, sort2.SortedOrderDesc)
}

func newSorterStruct[keyE sort2.KeyEntity, valE sort2.ValEntity](
	m SorterMapStruct[keyE, valE], sortedBy string, sortedOrder string,
) []Item[keyE, valE] {
	ms := &sorter[keyE, valE]{}
	ms.sortedBy = sortedBy
	ms.sortedOrder = sortedOrder

	ms.items = make([]Item[keyE, valE], 0, len(m))
	for k, v := range m {
		ms.items = append(ms.items, Item[keyE, valE]{Key: k, Val: v.SortedVal()})
	}

	return ms.sort()
}
