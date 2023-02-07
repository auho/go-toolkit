package maps

import (
	"fmt"
	"testing"
)

func Test_map(t *testing.T) {
	m := map[string]int{"2": 2, "3": 3, "1": 1}

	s := NewSorterByKey[string, int](m)
	ss := s.Sort()
	fmt.Println(ss)

	s = NewSorterByValue[string, int](m)
	ss = s.Sort()
	fmt.Println(ss)
}
