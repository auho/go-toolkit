package convert

// struct

type emptyStructType struct {
}

type structType struct {
	bool bool

	uint       uint
	uint8      uint8
	uint16     uint16
	uint32     uint32
	uint64     uint64
	int        int
	int8       int8
	int16      int16
	int32      int32
	int64      int64
	float32    float32
	float64    float64
	complex64  complex64
	complex128 complex128

	string string

	arrayBool   [3]bool
	arrayInt    [3]int
	arrayString [3]string

	sliceBool   []bool
	sliceInt    []int
	sliceString []string

	structField     emptyStructType
	ptrStructNotNil *emptyStructType // not nil
	ptrStructIsNil  *emptyStructType // is nil
}

type interfaceStruct struct {
	notNil interface{}
	isNIl  interface{}
}

// string

var _string = "string content"

// array

var _arrayBool = [3]bool{false, true, false}
var _arrayInt = [3]int{1, 3, 5}
var _arrayString = [3]string{"s1", "s3", "s5"}

// slice

var _sliceBool = []bool{false, true, false}
var _sliceInt = []int{1, 3, 5}
var _sliceString = []string{"s1", "s3", "s5"}

// map
var _mapBoolEmpty = map[bool]bool{}
var _mapBool = map[bool]bool{false: true, true: false}
var _mapInt = map[int]int{1: 1, 2: 2}
var _mapString = map[string]string{"s1": "s1", "s2": "s2", "s3": "s3"}

var _mapMultiPtr1 = map[anyMultiPtr1]anyMultiPtr1{_anyMultiPtr1: _anyMultiPtr1, _anyMultiPtr1: _anyMultiPtr1}
var _mapMultiPtr2 = map[anyMultiPtr2]anyMultiPtr2{_anyMultiPtr2: _anyMultiPtr2, _anyMultiPtr2: _anyMultiPtr2}
var _mapMultiPtr3 = map[anyMultiPtr3]anyMultiPtr3{_anyMultiPtr3: _anyMultiPtr3, _anyMultiPtr3: _anyMultiPtr3}
var _mapMultiPtr4 = map[anyMultiPtr4]anyMultiPtr4{_anyMultiPtr4: _anyMultiPtr4, _anyMultiPtr4: _anyMultiPtr4}
var _mapMultiPtr5 = map[anyMultiPtr5]anyMultiPtr5{_anyMultiPtr5: _anyMultiPtr5, _anyMultiPtr5: _anyMultiPtr5}

var _mapAnyNil = map[interface{}]interface{}{nil: nil}
var _mapAny = map[interface{}]interface{}{_anyMultiPtr1: _anyMultiPtr2}

// interface
type anyMultiPtr1 = *structType
type anyMultiPtr2 = **structType
type anyMultiPtr3 = ***structType
type anyMultiPtr4 = ****structType
type anyMultiPtr5 = *****structType

var _anyMultiPtr1 = &_struct
var _anyMultiPtr2 = &_anyMultiPtr1
var _anyMultiPtr3 = &_anyMultiPtr2
var _anyMultiPtr4 = &_anyMultiPtr3
var _anyMultiPtr5 = &_anyMultiPtr4

//
//

var _structEmpty = structType{}

var _struct = structType{
	bool:        false,
	uint:        1,
	uint8:       8,
	uint16:      16,
	uint32:      32,
	uint64:      64,
	int:         -1,
	int8:        -8,
	int16:       -16,
	int32:       -32,
	int64:       -64,
	float32:     1e-32,
	float64:     1e-64,
	complex64:   complex(1e-10, 1e-10),
	complex128:  complex(1e-10, 1e-10),
	string:      _string,
	arrayBool:   _arrayBool,
	arrayInt:    _arrayInt,
	arrayString: _arrayString,
	sliceBool:   _sliceBool,
	sliceInt:    _sliceInt,
	sliceString: _sliceString,

	structField:     emptyStructType{},
	ptrStructNotNil: &emptyStructType{},
	ptrStructIsNil:  nil,
}
