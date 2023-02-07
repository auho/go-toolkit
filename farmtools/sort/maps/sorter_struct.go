package maps

import (
	"sort"

	sort2 "github.com/auho/go-toolkit/farmtools/sort"
)

type SorterStruct[keyE sort2.KeyEntity, valE sort2.ValEntity] struct {
	Sorter[keyE, valE]
	origin map[keyE]Comparable[valE]
}

func NewSorterStructByKey[keyE sort2.KeyEntity, valE sort2.ValEntity](m map[keyE]Comparable[valE]) *SorterStruct[keyE, valE] {
	return newSorterStruct[keyE, valE](m, sortedByKey)
}

func NewSorterStructByValue[keyE sort2.KeyEntity, valE sort2.ValEntity](m map[keyE]Comparable[valE]) *SorterStruct[keyE, valE] {
	return newSorterStruct[keyE, valE](m, sortedByValue)
}

func newSorterStruct[keyE sort2.KeyEntity, valE sort2.ValEntity](m map[keyE]Comparable[valE], sortedBy string) *SorterStruct[keyE, valE] {
	ms := &SorterStruct[keyE, valE]{}
	ms.sortedBy = sortedBy
	ms.items = make([]Item[keyE, valE], 0, len(m))
	for k, v := range m {
		ms.items = append(ms.items, Item[keyE, valE]{Key: k, Val: v.SortedVal()})
	}

	return ms
}

func (ms *SorterStruct[keyE, valE]) Sort() []Item[keyE, valE] {
	sort.Sort(ms)
	return ms.items
}

func (ms *SorterStruct[keyE, valE]) Len() int {
	return len(ms.items)
}

func (ms *SorterStruct[keyE, valE]) Less(i, j int) bool {
	if ms.sortedBy == sortedByKey {
		return ms.items[i].Key < ms.items[j].Key
	} else {
		return ms.items[i].Val < ms.items[j].Val
	}
}

func (ms *SorterStruct[keyE, valE]) Swap(i, j int) {
	ms.items[i], ms.items[j] = ms.items[j], ms.items[i]
}
