package slices

import "testing"

func Test_slices(t *testing.T) {
	s := []int{10, 9, 8, 7, 6, 1, 2, 3, 4, 5, 0}

	NewSorterAsc(s)
	if s[0] != 0 {
		t.Error("asc error")
	}

	NewSorterDesc(s)
	if s[0] != 10 {
		t.Error("desc error")
	}
}
