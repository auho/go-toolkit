package convert

import (
	"reflect"
	"testing"
)

func Test_struct(t *testing.T) {
	type args struct {
		itemElem reflect.Value
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{"struct empty", args{reflect.ValueOf(_structEmpty)}, ""},
		{"struct", args{reflect.ValueOf(_struct)}, ""},
		{"struct pointer empty", args{reflect.ValueOf(&_structEmpty)}, ""},
		{"struct pointer", args{reflect.ValueOf(&_struct)}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := format(tt.args.itemElem); got != tt.want {
				_assert(t, got, tt.want)
			}
		})
	}
}
