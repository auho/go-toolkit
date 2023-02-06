package sort

import "sort"

const SortedByKey = "key"
const SortedByValue = "value"

type MapEntity interface {
	int | int64 | string
}

type MapSorter[keyE MapEntity, valE MapEntity] struct {
	items    []Item[keyE, valE]
	sortedBy string
}

type Item[keyE MapEntity, valE MapEntity] struct {
	Key keyE
	Val valE
}

func NewMapSorterByKey[keyE MapEntity, valE MapEntity](m map[keyE]valE) *MapSorter[keyE, valE] {
	return newMapSorter[keyE, valE](m, SortedByKey)
}

func NewMapSorterByValue[keyE MapEntity, valE MapEntity](m map[keyE]valE) *MapSorter[keyE, valE] {
	return newMapSorter[keyE, valE](m, SortedByValue)
}

func newMapSorter[keyE MapEntity, valE MapEntity](m map[keyE]valE, sortedBy string) *MapSorter[keyE, valE] {
	ms := &MapSorter[keyE, valE]{}
	ms.sortedBy = sortedBy
	ms.items = make([]Item[keyE, valE], 0, len(m))
	for k, v := range m {
		ms.items = append(ms.items, Item[keyE, valE]{Key: k, Val: v})
	}

	return ms
}

func (ms *MapSorter[keyE, valE]) Sort() []Item[keyE, valE] {
	sort.Sort(ms)
	return ms.items
}

func (ms *MapSorter[keyE, valE]) Len() int {
	return len(ms.items)
}

func (ms *MapSorter[keyE, valE]) Less(i, j int) bool {
	if ms.sortedBy == SortedByKey {
		return ms.items[i].Key < ms.items[j].Key
	} else {
		return ms.items[i].Val > ms.items[j].Val
	}
}

func (ms *MapSorter[keyE, valE]) Swap(i, j int) {
	ms.items[i], ms.items[j] = ms.items[j], ms.items[i]
}
