package hommization

import (
	"reflect"
	"strings"
)

func formatStruct(value reflect.Value) string {
	_filedNum := value.NumField()
	_valueType := value.Type()

	var b strings.Builder
	for i := 0; i < _filedNum; i++ {
		fieldElem := value.Field(i)
		fieldName := _valueType.Field(i).Name

		b.WriteString(`, "` + fieldName + `": `)
		b.WriteString(Format(fieldElem))
	}

	if b.Len() <= 0 {
		return addBraces("")
	} else {
		return addBraces(b.String()[2:])
	}
}
