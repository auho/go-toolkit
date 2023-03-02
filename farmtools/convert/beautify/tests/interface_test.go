package tests

import (
	"reflect"
	"testing"
)

func Test_interface(t *testing.T) {
	_testInterface(t, nil)
}

func Benchmark_interface(b *testing.B) {
	_testInterface(nil, b)

}

func _testInterface(t *testing.T, b *testing.B) {
	_anyMultiPtr3 = &_anyMultiPtr2

	type args struct {
		value interface{}
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{"interface struct empty", args{interfaceStruct{}}, ""},

		{"interface struct with int", args{interfaceStruct{1, nil}}, ""},

		{"interface with slice", args{interfaceStruct{_sliceInt, nil}}, ""},
		{"interface with slice 1", args{reflect.ValueOf(struct {
			sliceInt      interface{}
			sliceSliceInt interface{}
		}{_sliceInt, [][]int{_sliceInt}})}, ""},

		{"interface map", args{interfaceStruct{_mapInt, nil}}, ""},
		{"interface map 1", args{reflect.ValueOf(struct {
			mapIntInt       interface{}
			mapIntMapIntInt interface{}
		}{_mapInt, map[int]map[int]int{1: _mapInt}})}, ""},

		{"interface struct pointer", args{interfaceStruct{(*structType)(nil), nil}}, ""},
		{"interface struct pointer 1", args{interfaceStruct{(anyMultiPtr3)(nil), nil}}, ""},
		{"interface struct pointer 2", args{&struct {
			interfaceStruct        interface{}
			interfaceMultiPointer1 interface{}
			interfaceMultiPointer2 interface{}
			interfaceMultiPointer3 interface{}
			interfaceMultiPointer4 interface{}
			interfaceMultiPointer5 interface{}
		}{_struct, _anyMultiPtr1, _anyMultiPtr2, _anyMultiPtr3, _anyMultiPtr4, _anyMultiPtr5}}, ""},
	}

	for _, tt := range tests {
		if t != nil {
			t.Run(tt.name, func(t *testing.T) {
				if got := _decoder.Decode(tt.args.value); got != tt.want {
					_assert(t, got, tt.want)
				}
			})
		}

		if b != nil {
			b.Run(tt.name, func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_decoder.Decode(tt.args.value)
				}
			})
		}
	}
}
