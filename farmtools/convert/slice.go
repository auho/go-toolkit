package convert

import (
	"reflect"
	"strings"
)

func formatSlice(value reflect.Value) string {
	_len := value.Len()
	elmsString := make([]string, _len)

	for i := 0; i < _len; i++ {
		elmsString[i] = format(value.Index(i))
	}

	return "[" + strings.Join(elmsString, ", ") + "]"
}
