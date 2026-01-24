package contrast

type ByRowComparator struct {
	by
}

// bool: rows has changed
func (r *ByRowComparator) compare(baseData, compareData [][]string) (ByRowsRowsResult, bool) {
	var result ByRowsRowsResult
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

		result.addModifiedRows(r.indexToNo(rowIndex), cellsRet)
		hasChanged = true
		rowIndex++

	}

	// base rows > compare rows
	if baseTotalRows > compareTotalRows {
		for i := smallerTotalRows; i < baseTotalRows; i++ {
			result.addDeletedRow(r.indexToNo(rowIndex), baseData[i])
			hasChanged = true
			rowIndex++
		}

	} else if baseTotalRows < compareTotalRows {
		// base rows < compare rows
		for i := smallerTotalRows; i < compareTotalRows; i++ {
			result.addAddedRow(r.indexToNo(rowIndex), compareData[i])
			hasChanged = true
			rowIndex++
		}

	}

	return result, hasChanged
}

// bool: has changed
func (r *ByRowComparator) compareRow(baseRow, compareRow []string) ([]RowCellResult, bool) {
	var cellResults []RowCellResult
	var hasChanged bool

	baseLen := len(baseRow)
	compareLen := len(compareRow)

	smallerLen := min(baseLen, compareLen)

	// Compare common cells (present in both rows)
	for i := 0; i < smallerLen; i++ {
		cellResult := RowCellResult{
			col: r.indexToNo(i),
		}

		if baseRow[i] == compareRow[i] {
			cellResult.action = actionUnchanged
			cellResult.value = baseRow[i]
		} else {
			cellResult.action = actionModify
			cellResult.value = compareRow[i]

			hasChanged = true
		}

		cellResults = append(cellResults, cellResult)
	}

	// Base row has more cells than compare row
	if baseLen > compareLen {
		for i := smallerLen; i < baseLen; i++ {
			cellResults = append(cellResults, RowCellResult{
				col:    r.indexToNo(i),
				action: actionDelete,
				value:  baseRow[i],
			})
		}

		hasChanged = true
	} else if baseLen < compareLen {
		// Base row has fewer cells than compare row
		for i := smallerLen; i < compareLen; i++ {
			cellResults = append(cellResults, RowCellResult{
				col:    r.indexToNo(i),
				action: actionAdd,
				value:  compareRow[i],
			})
		}

		hasChanged = true
	}

	return cellResults, hasChanged
}
