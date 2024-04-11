package maps

import (
	"testing"
)

func TestSorterTable(t *testing.T) {
	m := map[string]int{"2": 2, "3": 3, "1": 1, "4": 4}

	// SorterTableAsc
	sKey, sVal := SorterTableAsc[string, int](m, func(key string) int {
		return m[key]
	})
	if sKey[0] != "1" {
		t.Error("table asc key error")
	}

	if sVal[0] != 1 {
		t.Error("table asc val error")
	}

	// SorterTableDesc
	sKey, sVal = SorterTableDesc[string, int](m, func(key string) int {
		return m[key]
	})
	if sKey[0] != "4" {
		t.Error("table asc key error")
	}

	if sVal[0] != 4 {
		t.Error("table asc val error")
	}
}
