package maps

import "github.com/auho/go-toolkit/farmtools/sort"

func SorterTableAsc[KT sort.KeyEntity, VT sort.ValEntity](x map[KT]VT, val func(key KT) VT) ([]KT, []VT) {
	return newSortByValue(x, val, sort.SortedOrderAsc)
}

func SorterTableDesc[KT sort.KeyEntity, VT sort.ValEntity](x map[KT]VT, val func(key KT) VT) ([]KT, []VT) {
	return newSortByValue(x, val, sort.SortedOrderDesc)
}
