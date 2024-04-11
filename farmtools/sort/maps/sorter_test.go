package maps

import (
	"testing"
)

func TestSorter(t *testing.T) {
	m := map[string]int{"2": 2, "3": 3, "1": 1, "4": 4}

	// SortKeyAsc
	sKey := SortKeyAsc[string, int](m)
	if sKey[0] != "1" {
		t.Error("key asc error")
	}

	// SortKeyDesc
	sKey = SortKeyDesc[string, int](m)
	if sKey[0] != "4" {
		t.Error("key desc error")
	}

	// SortValueAsc
	sKey, sVal := SortValueAsc[string, int](m)
	if sKey[0] != "1" {
		t.Error("table asc key error")
	}

	if sVal[0] != 1 {
		t.Error("table asc val error")
	}

	// SortValueDesc
	sKey, sVal = SortValueDesc[string, int](m)
	if sKey[0] != "4" {
		t.Error("table asc key error")
	}

	if sVal[0] != 4 {
		t.Error("table asc val error")
	}
}
