package slices

import (
	"sort"

	sort2 "github.com/auho/go-toolkit/farmtools/sort"
)

const sortedByKey = "key"
const sortedByValue = "value"

var _ sort.Interface = (*SorterAsc[string])(nil)
var _ sort.Interface = (*SorterDesc[string])(nil)

type Comparable[valE sort2.ValEntity] interface {
	SortedVal() valE
}

type SorterAsc[valE sort2.ValEntity] []valE

func (s SorterAsc[valE]) Len() int {
	return len(s)
}

func (s SorterAsc[valE]) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s SorterAsc[valE]) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type SorterDesc[valE sort2.ValEntity] []valE

func (s SorterDesc[valE]) Len() int {
	return len(s)
}

func (s SorterDesc[valE]) Less(i, j int) bool {
	return s[i] > s[j]
}

func (s SorterDesc[valE]) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func NewSorterAsc[valE sort2.ValEntity](s SorterAsc[valE]) {
	sort.Sort(s)
}

func NewSorterDesc[valE sort2.ValEntity](s SorterDesc[valE]) {
	sort.Sort(s)
}
