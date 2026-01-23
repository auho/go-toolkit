package difference

const (
	actionAdd    = 1
	actionModify = 2
	actionDelete = 3
)

type CellResult struct {
	row    int
	col    int
	action int
	value  string
}
