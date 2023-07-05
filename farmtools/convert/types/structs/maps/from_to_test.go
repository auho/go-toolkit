package maps

import (
	"testing"
)

type redisMap struct {
	Int     int     `json:"int"`
	Int64   int64   `json:"int64"`
	Float32 float32 `json:"float32"`
	Float64 float64 `json:"float64"`
	String  string  `json:"string"`
}

var err error
var s map[string]interface{}
var rm = redisMap{
	Int:     1,
	Int64:   2,
	Float32: 3.3,
	Float64: 4.4,
	String:  "5",
}

func Test_from(t *testing.T) {
	// not pointer
	s, err = MapStringAnyFromStruct(rm)
	if err != nil {
		t.Fatal(err)
	}

	// pointer
	s, err = MapStringAnyFromStruct(&rm)
	if err != nil {
		t.Fatal(err)
	}

	var rm1 redisMap
	err = MapStringAnyToStruct(&rm1, s)
	if err != nil {
		t.Fatal(err)
	}

	if rm != rm1 {
		t.Error("not equal")
	}
}

// error type
func Test_error(t *testing.T) {
	_i := 2
	s, err = MapStringAnyFromStruct(_i)
	if err == nil {
		t.Fatal("error type _i")
	}

	s, err = MapStringAnyFromStruct(&_i)
	if err == nil {
		t.Fatal("error type &_i")
	}

	var _i1 interface{} = 2
	s, err = MapStringAnyFromStruct(_i1)
	if err == nil {
		t.Fatal("error type _i1")
	}

	s, err = MapStringAnyFromStruct(&_i1)
	if err == nil {
		t.Fatal("error type &_i1")
	}
}
