package convert

import (
	"reflect"
	"strings"
)

func formatMap(value reflect.Value) string {
	var b strings.Builder
	iterator := value.MapRange()
	for iterator.Next() {
		k := iterator.Key()
		v := iterator.Value()

		ks, _isLiteral := underlyingKindString(k)
		if !_isLiteral {
			ks = typeStringToSymbolString(ks)
		}

		b.WriteString(", " + ks + `: `)
		b.WriteString(format(v))
	}

	if b.Len() <= 0 {
		return addBraces("")
	} else {
		return addBraces(b.String()[2:])
	}
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
				s = value.Kind().String()
			}
		} else {
			s = value.Type().String()
		}
	}

	return s, _isLiteral
}
