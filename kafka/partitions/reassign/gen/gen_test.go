package gen

import (
	"fmt"
	"testing"
)

func TestGen_combination(t *testing.T) {
	gen := &Gen{}

	ids := []int{1, 2, 3}
	rets := gen.combination(ids, 3, 2)
	fmt.Println(rets)

	ids = []int{1, 2, 3}
	rets = gen.combination(ids, 3, 3)
	fmt.Println(rets)

	ids = []int{1, 2, 3}
	rets = gen.combination(ids, 3, 4)
	fmt.Println(rets)

	ids = []int{1, 2, 3, 4, 5}
	rets = gen.combination(ids, 3, 2)
	fmt.Println(rets)
}
