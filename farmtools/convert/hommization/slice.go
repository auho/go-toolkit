package hommization

import (
	"reflect"
	"strings"
)

func formatSlice(value reflect.Value) string {
	_len := value.Len()

	var b strings.Builder
	for i := 0; i < _len; i++ {
		b.WriteString(", ")
		b.WriteString(Format(value.Index(i)))
	}

	if b.Len() <= 0 {
		return addBrackets("")
	} else {
		return addBrackets(b.String()[2:])
	}
}
