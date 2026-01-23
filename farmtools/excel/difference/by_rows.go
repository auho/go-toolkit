package difference

type ByRows struct {
	by

	baseData    [][]string
	compareData [][]string
}

func (r *ByRows) compare() (ByRowResult, error) {
	var byRowsRet ByRowResult

	baseTotalRows := len(r.baseData)
	compareTotalRows := len(r.compareData)

	smallerTotalRows := min(baseTotalRows, compareTotalRows)

	var rowIndex int

	// common rows
	for i := 0; i < smallerTotalRows; i++ {
		cellsRet, ok := r.compareRow(i)
		if !ok {
			continue
		}

		byRowsRet.addModifiedRows(r.indexToNo(rowIndex), cellsRet)
		rowIndex++
	}

	// base rows > compare rows
	if baseTotalRows > compareTotalRows {
		for i := smallerTotalRows; i < baseTotalRows; i++ {
			byRowsRet.addDeletedRow(r.indexToNo(rowIndex), r.baseData[i])
			rowIndex++
		}

	} else if baseTotalRows < compareTotalRows {
		// base rows < compare rows
		for i := smallerTotalRows; i < compareTotalRows; i++ {
			byRowsRet.addAddedRow(r.indexToNo(rowIndex), r.compareData[i])
			rowIndex++
		}

	}

	return byRowsRet, nil
}

func (r *ByRows) compareRow(rowIndex int) ([]RowCellResult, bool) {
	var cellsRet []RowCellResult
	var hasChange bool

	baseRow := r.baseData[rowIndex]
	compareRow := r.compareData[rowIndex]

	baseLen := len(baseRow)
	compareLen := len(compareRow)

	smallerLen := min(baseLen, compareLen)

	// common
	for i := 0; i < smallerLen; i++ {
		cellRet := RowCellResult{
			col: r.indexToNo(i),
		}

		if baseRow[i] == compareRow[i] {
			cellRet.action = actionNull
			cellRet.value = baseRow[i]
		} else {
			cellRet.action = actionModify
			cellRet.value = compareRow[i]

			hasChange = true
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

		hasChange = true
	} else if baseLen < compareLen {
		// base len < compare len
		for i := smallerLen; i < compareLen; i++ {
			cellsRet = append(cellsRet, RowCellResult{
				col:    r.indexToNo(i),
				action: actionAdd,
				value:  compareRow[i],
			})
		}

		hasChange = true
	}

	return cellsRet, hasChange
}
