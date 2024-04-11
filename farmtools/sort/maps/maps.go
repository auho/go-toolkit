package maps

import (
	"github.com/auho/go-toolkit/farmtools/sort"
	sort2 "sort"
)

func newSortByKey[KT sort.KeyEntity, VT sort.ValEntity](x map[KT]VT, sortedOrder string) []KT {
	var _ks []KT

	for _k := range x {
		_ks = append(_ks, _k)
	}

	sort2.SliceStable(_ks, func(i, j int) bool {
		if sortedOrder == sort.SortedOrderAsc {
			return _ks[i] < _ks[j]
		} else {
			return _ks[i] > _ks[j]
		}
	})

	return _ks
}

func newSortByValue[KT sort.KeyEntity, VT sort.ValEntity](x map[KT]VT, val func(key KT) VT, sortedOrder string) ([]KT, []VT) {
	var _ks []KT
	_vs := make(map[KT]VT)

	for _k := range x {
		_ks = append(_ks, _k)
		_vs[_k] = val(_k)
	}

	sort2.SliceStable(_ks, func(i, j int) bool {
		_ki := _ks[i]
		_kj := _ks[j]

		if sortedOrder == sort.SortedOrderAsc {
			return _vs[_ki] < _vs[_kj]
		} else {
			return _vs[_ki] > _vs[_kj]
		}
	})

	var _nvs []VT
	for _, _k := range _ks {
		_nvs = append(_nvs, _vs[_k])
	}

	return _ks, _nvs
}
