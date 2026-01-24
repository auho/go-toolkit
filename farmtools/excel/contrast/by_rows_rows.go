package contrast

type ByRowsRows struct {
	by
}

// bool: rows has changed
func (r *ByRowsRows) compare(baseData, compareData [][]string) (ByRowsRowsResult, bool) {
	var ret ByRowsRowsResult
	var hasChanged bool

	baseTotalRows := len(baseData)
	compareTotalRows := len(compareData)

	smallerTotalRows := min(baseTotalRows, compareTotalRows)

	var rowIndex int

	// both have rows
	for i := 0; i < smallerTotalRows; i++ {
		cellsRet, ok := r.compareRow(baseData[i], compareData[i])
		if !ok {
			continue
		}

		ret.addModifiedRows(r.indexToNo(rowIndex), cellsRet)
		hasChanged = true
		rowIndex++

	}

	// base rows > compare rows
	if baseTotalRows > compareTotalRows {
		for i := smallerTotalRows; i < baseTotalRows; i++ {
			ret.addDeletedRow(r.indexToNo(rowIndex), baseData[i])
			hasChanged = true
			rowIndex++
		}

	} else if baseTotalRows < compareTotalRows {
		// base rows < compare rows
		for i := smallerTotalRows; i < compareTotalRows; i++ {
			ret.addAddedRow(r.indexToNo(rowIndex), compareData[i])
			hasChanged = true
			rowIndex++
		}

	}

	return ret, hasChanged
}

// bool: has changed
func (r *ByRowsRows) compareRow(baseRow, compareRow []string) ([]RowCellResult, bool) {
	var cellsRet []RowCellResult
	var hasChanged bool

	baseLen := len(baseRow)
	compareLen := len(compareRow)

	smallerLen := min(baseLen, compareLen)

	// common
	for i := 0; i < smallerLen; i++ {
		cellRet := RowCellResult{
			col: r.indexToNo(i),
		}

		if baseRow[i] == compareRow[i] {
			cellRet.action = actionUnchanged
			cellRet.value = baseRow[i]
		} else {
			cellRet.action = actionModify
			cellRet.value = compareRow[i]

			hasChanged = true
		}

		cellsRet = append(cellsRet, cellRet)
	}

	// base len > compare len
	if baseLen > compareLen {
		for i := smallerLen; i < baseLen; i++ {
			cellsRet = append(cellsRet, RowCellResult{
				col:    r.indexToNo(i),
				action: actionDelete,
				value:  baseRow[i],
			})
		}

		hasChanged = true
	} else if baseLen < compareLen {
		// base len < compare len
		for i := smallerLen; i < compareLen; i++ {
			cellsRet = append(cellsRet, RowCellResult{
				col:    r.indexToNo(i),
				action: actionAdd,
				value:  compareRow[i],
			})
		}

		hasChanged = true
	}

	return cellsRet, hasChanged
}
