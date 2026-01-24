package contrast

type ByColumnsResult struct {
	catalog CatalogResult
	sheets  []ByColumnsSheetResult
}

type ByColumnsSheetResult struct {
	hasChanged bool
	sheetName  string
	columns    ByColumnsColumnsResult
}

type ByColumnsColumnsResult struct {
	modified []ColResult
	added    []ColData
	deleted  []ColData
}

func (r *ByColumnsColumnsResult) addAddedColumn(colIndex int, data []string) {
	r.added = append(r.added, ColData{colIndex, data})
}

func (r *ByColumnsColumnsResult) addDeletedColumn(colIndex int, data []string) {
	r.deleted = append(r.deleted, ColData{colIndex, data})
}

func (r *ByColumnsColumnsResult) addModifiedColumn(colIndex int, cellsRet []ColCellResult) {
	r.modified = append(r.modified, ColResult{colIndex, cellsRet})
}

func (r *ByColumnsColumnsResult) modifiedToCellResult() []CellResult {
	var cellsRet []CellResult

	for _, col := range r.modified {
		cells := col.toCell()
		cellsRet = append(cellsRet, cells...)
	}

	return cellsRet
}

type ColResult struct {
	col   int
	cells []ColCellResult
}

func (r *ColResult) toCell() []CellResult {
	var cells []CellResult
	for _, cell := range r.cells {
		cells = append(cells, CellResult{
			row:    cell.row,
			col:    r.col,
			action: cell.action,
			value:  cell.value,
		})
	}

	return cells
}

type ColCellResult struct {
	row    int
	action int
	value  string
}
