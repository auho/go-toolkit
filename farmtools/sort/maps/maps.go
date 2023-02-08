package maps

import sort2 "github.com/auho/go-toolkit/farmtools/sort"

const sortedByKey = "key"
const sortedByValue = "value"

type Sorter[keyE sort2.KeyEntity, valE sort2.ValEntity] struct {
	items       []Item[keyE, valE]
	sortedBy    string
	sortedOrder string
}

type Item[keyE sort2.KeyEntity, valE sort2.ValEntity] struct {
	Key keyE
	Val valE
}
