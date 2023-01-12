package convert

import (
	"reflect"
	"testing"
)

func Test_interface(t *testing.T) {
	_anyMultiPtr3 = &_anyMultiPtr2

	type args struct {
		itemElem reflect.Value
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{"interface struct empty", args{reflect.ValueOf(interfaceStruct{})}, ""},

		{"interface struct with int", args{reflect.ValueOf(interfaceStruct{1, nil})}, ""},

		{"interface with slice", args{reflect.ValueOf(interfaceStruct{_sliceInt, nil})}, ""},
		{"interface with slice 1", args{reflect.ValueOf(struct {
			sliceInt      interface{}
			sliceSliceInt interface{}
		}{_sliceInt, [][]int{_sliceInt}})}, ""},

		{"interface map", args{reflect.ValueOf(interfaceStruct{_mapInt, nil})}, ""},
		{"interface map 1", args{reflect.ValueOf(struct {
			mapIntInt       interface{}
			mapIntMapIntInt interface{}
		}{_mapInt, map[int]map[int]int{1: _mapInt}})}, ""},

		{"interface struct pointer", args{reflect.ValueOf(interfaceStruct{(*structType)(nil), nil})}, ""},
		{"interface struct pointer 1", args{reflect.ValueOf(interfaceStruct{(anyMultiPtr3)(nil), nil})}, ""},
		{"interface struct pointer 2", args{reflect.ValueOf(&struct {
			interfaceStruct        interface{}
			interfaceMultiPointer1 interface{}
			interfaceMultiPointer2 interface{}
			interfaceMultiPointer3 interface{}
			interfaceMultiPointer4 interface{}
			interfaceMultiPointer5 interface{}
		}{_struct, _anyMultiPtr1, _anyMultiPtr2, _anyMultiPtr3, _anyMultiPtr4, _anyMultiPtr5})}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := format(tt.args.itemElem); got != tt.want {
				_assert(t, got, tt.want)
			}
		})
	}
}
