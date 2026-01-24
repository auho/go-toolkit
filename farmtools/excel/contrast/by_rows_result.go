package contrast

type ByRowsResult struct {
	catalog CatalogResult
	sheets  []ByRowsSheetResult
}

type ByRowsSheetResult struct {
	hasChanged bool
	sheetName  string
	rows       ByRowsRowsResult
}

type ByRowsRowsResult struct {
	modified []RowResult
	added    []RowData
	deleted  []RowData
}

func (r *ByRowsRowsResult) addAddedRow(rowIndex int, data []string) {
	r.added = append(r.added, RowData{rowIndex, data})
}

func (r *ByRowsRowsResult) addDeletedRow(rowIndex int, data []string) {
	r.deleted = append(r.deleted, RowData{rowIndex, data})
}

func (r *ByRowsRowsResult) addModifiedRows(rowIndex int, cellsRet []RowCellResult) {
	r.modified = append(r.modified, RowResult{rowIndex, cellsRet})
}

func (r *ByRowsRowsResult) modifiedToCellResult() []CellResult {
	var cellsRet []CellResult

	for _, row := range r.modified {
		cells := row.toCell()
		cellsRet = append(cellsRet, cells...)
	}

	return cellsRet
}
