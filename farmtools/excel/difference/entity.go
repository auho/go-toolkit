package difference

//
// row col 从 1 开始
//

type CellResult struct {
	row    int
	col    int
	action int
	value  string
}

//
// rows
//

type RowCellResult struct {
	col    int
	action int
	value  string
}

type RowResult struct {
	row   int
	cells []RowCellResult
}

func (r *RowResult) toCell() []CellResult {
	var cells []CellResult
	for _, cell := range r.cells {
		cells = append(cells, CellResult{
			row:    r.row,
			col:    cell.col,
			action: cell.action,
			value:  cell.value,
		})
	}

	return cells
}

//
// data
//

type RowData struct {
	row  int
	data []string
}

type ColData struct {
	col  int
	data []string
}
