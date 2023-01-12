package convert

import (
	"reflect"
	"strings"
)

func formatMap(value reflect.Value) string {
	elmsString := make([]string, value.Len())

	i := 0
	iterator := value.MapRange()
	for iterator.Next() {
		k := iterator.Key()
		v := iterator.Value()

		ks, isLiteral := underlyingKindString(k)
		if !isLiteral {
			ks = `"- ` + ks + ` -"`
		}

		elmsString[i] = ks + `: ` + format(v)
		i++
	}

	return "{" + strings.Join(elmsString, ", ") + "}"
}

func underlyingKindString(value reflect.Value) (s string, isLiteral bool) {
	isLiteral = false

	kind := value.Kind()

	switch {
	case _isLiteral(kind):
		s = format(value)
		isLiteral = true
	case kind == reflect.Pointer:
		s, _ = underlyingKindString(value.Elem())
		s = "*" + s
	default:
		if value.Type().Kind() == reflect.Interface {
			if value.IsNil() {
				s = CNull
				isLiteral = true
			} else {
				s = value.Type().Kind().String()
			}
		} else {
			s = value.Type().String()
		}
	}

	return s, isLiteral
}
