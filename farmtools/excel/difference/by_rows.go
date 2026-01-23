package difference

type ByRows struct {
	baseData    [][]string
	compareData [][]string
}

func (d *ByRows) diff() (SheetResult, error) {
	addedRows := make(map[int][]string)
	deletedRows := make(map[int][]string)
	var modifiedCells []CellResult

	baseTotalRows := len(d.baseData)
	compareTotalRows := len(d.compareData)

	smallerTotalRows := min(baseTotalRows, compareTotalRows)

	// common rows
	for i := 0; i < smallerTotalRows; i++ {

		rowRet := d.row(i)
	}

	// base rows > compare rows
	if baseTotalRows > compareTotalRows {
		for i := smallerTotalRows; i < baseTotalRows; i++ {
			rowRet := d.rowWithAction(i, d.baseData[i], actionDelete)
		}

	} else if baseTotalRows < compareTotalRows {
		// base rows < compare rows
		for i := smallerTotalRows; i < compareTotalRows; i++ {
			rowRet := d.rowWithAction(i, d.compareData[i], actionAdd)
		}
	}

	// base rows < compare rows

	if baseTotalRows < compareTotalRows { // length of rows of base less than diff. new cells (row)
		for i := baseTotalRows + 1; i < compareTotalRows; i++ {
			addedRows[i] = d.compareData[i]
		}
	} else { // length of rows of base greater than diff. new cells (row)
		for baseRowNo, baseRowValues := range d.baseData { // iteration row
			// rows
			if baseRowNo > compareTotalRows { // delete cells (row)
				deletedRows[baseRowNo] = baseRowValues
			} else { // change cells (row)

			}
		}
	}

	return SheetResult{
		addedRows:     addedRows,
		deletedRows:   deletedRows,
		modifiedCells: modifiedCells,
	}, nil
}

func (d *ByRows) row(rowIndex int) []CellResult {
	var cellsRet []CellResult

	baseRow := d.baseData[rowIndex]
	compareRow := d.compareData[rowIndex]

	baseLen := len(baseRow)
	compareLen := len(compareRow)

	smallerLen := min(baseLen, compareLen)

	// common
	for i := 0; i < smallerLen; i++ {
		if baseRow[i] != compareRow[i] {
			cellsRet = append(cellsRet, CellResult{
				row:    rowIndex,
				col:    i,
				action: actionModify,
				value:  compareRow[i],
			})
		}
	}

	// base len > compare len
	if baseLen > compareLen {
		for i := smallerLen; i < baseLen; i++ {
			cellsRet = append(cellsRet, CellResult{
				row:    rowIndex,
				col:    i,
				action: actionDelete,
				value:  baseRow[i],
			})
		}
	} else if baseLen < compareLen {
		// base len < compare len
		for i := smallerLen; i < compareLen; i++ {
			cellsRet = append(cellsRet, CellResult{
				row:    rowIndex,
				col:    i,
				action: actionAdd,
				value:  compareRow[i],
			})
		}
	}

	return cellsRet
}

func (d *ByRows) rowWithAction(rowIndex int, row []string, action int) []CellResult {
	var cellsRet []CellResult
	for i, v := range row {
		cellsRet = append(cellsRet, CellResult{
			row:    rowIndex,
			col:    i,
			action: action,
			value:  v,
		})
	}

	return cellsRet
}
