package contrast

type ByColumnsColumns struct {
	by
}

// bool: columns has changed
func (c *ByColumnsColumns) compare(baseData, compareData [][]string) (ByColumnsColumnsResult, bool) {
	var ret ByColumnsColumnsResult
	var hasChanged bool

	baseTotalCols := len(baseData)
	compareTotalCols := len(compareData)

	smallerTotalCols := min(baseTotalCols, compareTotalCols)

	// Compare common columns
	for i := 0; i < smallerTotalCols; i++ {
		baseCol := baseData[i]
		compareCol := compareData[i]

		cellsRet, ok := c.compareColumn(baseCol, compareCol)
		if !ok {
			continue
		}

		ret.addModifiedColumn(c.indexToNo(i), cellsRet)
		hasChanged = true
	}

	// base columns > compare columns
	if baseTotalCols > compareTotalCols {
		for i := smallerTotalCols; i < baseTotalCols; i++ {
			baseCol := baseData[i]
			ret.addDeletedColumn(c.indexToNo(i), baseCol)
			hasChanged = true
		}

	} else if baseTotalCols < compareTotalCols {
		// base columns < compare columns
		for i := smallerTotalCols; i < compareTotalCols; i++ {
			compareCol := compareData[i]
			ret.addAddedColumn(c.indexToNo(i), compareCol)
			hasChanged = true
		}

	}

	return ret, hasChanged
}

// bool: has changed
func (c *ByColumnsColumns) compareColumn(baseCol, compareCol []string) ([]ColCellResult, bool) {
	var cellsRet []ColCellResult
	var hasChanged bool

	baseLen := len(baseCol)
	compareLen := len(compareCol)

	smallerLen := min(baseLen, compareLen)

	// common rows
	for i := 0; i < smallerLen; i++ {
		cellRet := ColCellResult{
			row: c.indexToNo(i),
		}

		if baseCol[i] == compareCol[i] {
			cellRet.action = actionUnchanged
			cellRet.value = baseCol[i]
		} else {
			cellRet.action = actionModify
			cellRet.value = compareCol[i]

			hasChanged = true
		}

		cellsRet = append(cellsRet, cellRet)
	}

	// base len > compare len
	if baseLen > compareLen {
		for i := smallerLen; i < baseLen; i++ {
			cellsRet = append(cellsRet, ColCellResult{
				row:    c.indexToNo(i),
				action: actionDelete,
				value:  baseCol[i],
			})
		}

		hasChanged = true
	} else if baseLen < compareLen {
		// base len < compare len
		for i := smallerLen; i < compareLen; i++ {
			cellsRet = append(cellsRet, ColCellResult{
				row:    c.indexToNo(i),
				action: actionAdd,
				value:  compareCol[i],
			})
		}

		hasChanged = true
	}

	return cellsRet, hasChanged
}
