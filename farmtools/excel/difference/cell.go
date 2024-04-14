package difference

const (
	actionAdd int = iota
	actionModify
	actionDelete
)

type CellDiff struct {
	row    int
	col    int
	action int
	value  string
}
