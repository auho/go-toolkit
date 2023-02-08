package hommization

import "reflect"

func formatChan(value reflect.Value) string {
	return addDoubleQuotationMark(addTypeSymbol(value.Type().String()))
}
