package maps

import (
	"github.com/auho/go-toolkit/farmtools/sort"
)

func SortKeyAsc[KT sort.KeyEntity, VT sort.ValEntity](x map[KT]VT) []KT {
	return newSortByKey(x, sort.SortedOrderAsc)
}

func SortKeyDesc[KT sort.KeyEntity, VT sort.ValEntity](x map[KT]VT) []KT {
	return newSortByKey(x, sort.SortedOrderDesc)
}

func SortValueAsc[KT sort.KeyEntity, VT sort.ValEntity](x map[KT]VT) ([]KT, []VT) {
	return newSortByValue(x, func(key KT) VT {
		return x[key]
	}, sort.SortedOrderAsc)
}

func SortValueDesc[KT sort.KeyEntity, VT sort.ValEntity](x map[KT]VT) ([]KT, []VT) {
	return newSortByValue(x, func(key KT) VT {
		return x[key]
	}, sort.SortedOrderDesc)
}
