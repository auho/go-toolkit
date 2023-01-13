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

		ks, _isLiteral := underlyingKindString(k)
		if !_isLiteral {
			ks = typeStringToSymbolString(ks)
		}

		elmsString[i] = ks + `: ` + format(v)
		i++
	}

	return addBraces(strings.Join(elmsString, ", "))
}

func underlyingKindString(value reflect.Value) (s string, _isLiteral bool) {
	_isLiteral = false

	kind := value.Kind()

	switch {
	case isLiteral(kind):
		s = format(value)
		_isLiteral = true
	case kind == reflect.Pointer:
		s, _ = underlyingKindString(value.Elem())
		s = addPointerSymbol(s)
	default:
		if value.Type().Kind() == reflect.Interface {
			if value.IsNil() {
				s = CNull
				_isLiteral = true
			} else {
				s = value.Type().Kind().String()
			}
		} else {
			s = value.Type().String()
		}
	}

	return s, _isLiteral
}
