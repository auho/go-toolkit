package difference

type ByRowResult struct {
	modifiedRows []RowResult
	addedRows    []RowData
	deletedRows  []RowData
}

func (r *ByRowResult) addAddedRow(rowIndex int, data []string) {
	r.addedRows = append(r.addedRows, RowData{rowIndex, data})
}

func (r *ByRowResult) addDeletedRow(rowIndex int, data []string) {
	r.deletedRows = append(r.deletedRows, RowData{rowIndex, data})
}

func (r *ByRowResult) addModifiedRows(rowIndex int, cellsRet []RowCellResult) {
	r.modifiedRows = append(r.modifiedRows, RowResult{rowIndex, cellsRet})
}

func (r *ByRowResult) modifiedToCellResult() []CellResult {
	var cellsRet []CellResult

	for _, row := range r.modifiedRows {
		cells := row.toCell()
		cellsRet = append(cellsRet, cells...)
	}

	return cellsRet
}
