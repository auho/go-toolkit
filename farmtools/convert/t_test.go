package convert

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_t(t *testing.T) {
	type s struct {
		a int
	}

	var s3 ***s

	s1 := &s{a: 1}
	s2 := &s1
	s3 = &s2

	ref := reflect.ValueOf(s3)

	fmt.Println(ref.Type().Elem().Kind())
}
