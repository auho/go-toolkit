package convert

import (
	"reflect"
	"testing"
)

func Test_map(t *testing.T) {
	type args struct {
		value reflect.Value
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{"map bool empty", args{reflect.ValueOf(_mapBoolEmpty)}, ""},
		{"map bool", args{reflect.ValueOf(_mapBool)}, ""},
		{"map int", args{reflect.ValueOf(_mapInt)}, ""},
		{"map string", args{reflect.ValueOf(_mapString)}, ""},

		{"map pointer 1", args{reflect.ValueOf(_mapMultiPtr1)}, ""},
		{"map pointer 2", args{reflect.ValueOf(_mapMultiPtr2)}, ""},
		{"map pointer 3", args{reflect.ValueOf(_mapMultiPtr3)}, ""},
		{"map pointer 4", args{reflect.ValueOf(_mapMultiPtr4)}, ""},
		{"map pointer 5", args{reflect.ValueOf(_mapMultiPtr5)}, ""},

		{"map interface nil", args{reflect.ValueOf(_mapAnyNil)}, ""},
		{"map interface", args{reflect.ValueOf(_mapAny)}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := format(tt.args.value); got != tt.want {
				_assert(t, got, tt.want)
			}
		})
	}
}
