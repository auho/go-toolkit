package contrast

type ByColumnComparator struct {
	by
}

// bool: columns has changed
func (c *ByColumnComparator) compare(baseData, compareData [][]string) (ByColumnsColumnsResult, bool) {
	var result ByColumnsColumnsResult
	var hasChanged bool

	baseTotalCols := len(baseData)
	compareTotalCols := len(compareData)

	smallerTotalCols := min(baseTotalCols, compareTotalCols)

	// Compare common columns (present in both columns)
	for i := 0; i < smallerTotalCols; i++ {
		baseCol := baseData[i]
		compareCol := compareData[i]

		cellsRet, ok := c.compareColumn(baseCol, compareCol)
		if !ok {
			continue
		}

		result.addModifiedColumn(c.indexToNo(i), cellsRet)
		hasChanged = true
	}

	// base columns > compare columns
	if baseTotalCols > compareTotalCols {
		for i := smallerTotalCols; i < baseTotalCols; i++ {
			baseCol := baseData[i]
			result.addDeletedColumn(c.indexToNo(i), baseCol)
			hasChanged = true
		}

	} else if baseTotalCols < compareTotalCols {
		// base columns < compare columns
		for i := smallerTotalCols; i < compareTotalCols; i++ {
			compareCol := compareData[i]
			result.addAddedColumn(c.indexToNo(i), compareCol)
			hasChanged = true
		}

	}

	return result, hasChanged
}

// bool: has changed
func (c *ByColumnComparator) compareColumn(baseCol, compareCol []string) ([]ColCellResult, bool) {
	var cellResults []ColCellResult
	var hasChanged bool

	baseLen := len(baseCol)
	compareLen := len(compareCol)

	smallerLen := min(baseLen, compareLen)

	// Compare common cells (present in both columns)
	for i := 0; i < smallerLen; i++ {
		cellResult := ColCellResult{
			row: c.indexToNo(i),
		}

		if baseCol[i] == compareCol[i] {
			cellResult.action = actionUnchanged
			cellResult.value = baseCol[i]
		} else {
			cellResult.action = actionModify
			cellResult.value = compareCol[i]

			hasChanged = true
		}

		cellResults = append(cellResults, cellResult)
	}

	// Base column has more cells than compare column
	if baseLen > compareLen {
		for i := smallerLen; i < baseLen; i++ {
			cellResults = append(cellResults, ColCellResult{
				row:    c.indexToNo(i),
				action: actionDelete,
				value:  baseCol[i],
			})
		}

		hasChanged = true
	} else if baseLen < compareLen {
		// Base column has fewer cells than compare column
		for i := smallerLen; i < compareLen; i++ {
			cellResults = append(cellResults, ColCellResult{
				row:    c.indexToNo(i),
				action: actionAdd,
				value:  compareCol[i],
			})
		}

		hasChanged = true
	}

	return cellResults, hasChanged
}
