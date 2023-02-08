package sort

const SortedOrderAsc = "asc"
const SortedOrderDesc = "desc"

type intEntity interface {
	int8 | int16 | int32 | int64 | int
}

type uintEntity interface {
	uint8 | uint16 | uint32 | uint64 | uint
}

type floatEntity interface {
	float32 | float64
}

type KeyEntity interface {
	intEntity | uintEntity | floatEntity | string
}

type ValEntity interface {
	KeyEntity
}
