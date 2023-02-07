package maps

import (
	"fmt"
	"testing"
)

var _ Comparable[string] = (*sorterStruct)(nil)

type sorterStruct struct {
	b string
}

func (m sorterStruct) SortedVal() string {
	return m.b
}

func Test_SorterStruct(t *testing.T) {
	m := map[string]sorterStruct{"2": {"2"}, "1": {"1"}}

	nm := make(map[string]Comparable[string], len(m))
	for k := range m {
		nm[k] = m[k]
	}

	s := NewSorterStructByValue[string](nm)
	ss := s.Sort()
	fmt.Println(ss)

	s = NewSorterStructByKey[string](nm)
	ss = s.Sort()
	fmt.Println(ss)
}
