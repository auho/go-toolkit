package convert

import (
	"reflect"
	"testing"
)

func Test__slice(t *testing.T) {
	type args struct {
		value reflect.Value
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{"bool", args{reflect.ValueOf([]bool{false, false, true, true})}, ""},
		{"int", args{reflect.ValueOf([]int{1, 2, 3, 4, 5})}, ""},
		{"string", args{reflect.ValueOf([]string{"s1", "s2", "s3", "s4", "s5"})}, ""},
		{"struct", args{reflect.ValueOf([]structType{_struct})}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatSlice(tt.args.value); got != tt.want {
				_assert(t, got, tt.want)
			}
		})
	}
}
