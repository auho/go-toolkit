package sort

type KeyEntity interface {
	int | int64 | string
}

type ValEntity interface {
	KeyEntity
}
