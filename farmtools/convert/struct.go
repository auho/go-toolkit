package convert

import (
	"reflect"
	"strings"
)

func formatStruct(value reflect.Value) string {
	_filedNum := value.NumField()
	elmsString := make([]string, _filedNum)

	for i := 0; i < _filedNum; i++ {
		fieldElem := value.Field(i)
		fieldName := value.Type().Field(i).Name

		elmsString[i] = `"` + fieldName + `": ` + format(fieldElem)
	}

	return addBraces(strings.Join(elmsString, ", "))
}
