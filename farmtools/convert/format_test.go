package convert

import (
	"reflect"
	"testing"
)

func Test_format(t *testing.T) {
	type args struct {
		value reflect.Value
	}

	var structAnyEmpty interface{} = _structEmpty
	var structAny interface{} = _struct
	var structPtrAnyNil interface{} = (*structType)(nil)
	var structPtrAnyEmpty interface{} = &_structEmpty
	var structPtrAny interface{} = &_struct

	tests := []struct {
		name string
		args args
		want string
	}{
		{"bool false", args{reflect.ValueOf(false)}, ""},
		{"bool true", args{reflect.ValueOf(true)}, ""},
		{"uint", args{reflect.ValueOf(1)}, ""},
		{"int", args{reflect.ValueOf(-1)}, ""},

		{"struct interface empty", args{reflect.ValueOf(structAnyEmpty)}, ""},
		{"struct interface", args{reflect.ValueOf(structAny)}, ""},
		{"struct pointer interface empty", args{reflect.ValueOf(structPtrAnyEmpty)}, ""},
		{"struct pointer interface nil", args{reflect.ValueOf(structPtrAnyNil)}, ""},
		{"struct pointer interface", args{reflect.ValueOf(structPtrAny)}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := format(tt.args.value); got != tt.want {
				_assert(t, got, tt.want)
			}
		})
	}
}
