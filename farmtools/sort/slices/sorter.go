package slices

import (
	"sort"

	sort2 "github.com/auho/go-toolkit/farmtools/sort"
)

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
