package convert

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
		return addBraces("")
	} else {
		return addBraces(b.String()[2:])
	}
}
