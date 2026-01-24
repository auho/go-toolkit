package contrast

import "fmt"

type ByColumns struct {
	excel

	columns *ByColumnsColumns
}

func NewByColumns() *ByColumns {
	return &ByColumns{
		columns: &ByColumnsColumns{},
	}
}

func (c *ByColumns) CompareFilePath(input InputFilePath) (ByColumnsOutput, error) {
	var ret ByColumnsOutput

	inputExcelFile, err := c.inputOpenFile(input)
	if err != nil {
		return ret, fmt.Errorf("inputOpenFile: %w", err)
	}

	defer func() { _ = inputExcelFile.close() }()

	ret, err = c.Compare(inputExcelFile)
	if err != nil {
		return ret, fmt.Errorf("compare: %w", err)
	}

	return ret, nil
}

func (c *ByColumns) Compare(input Input) (ByColumnsOutput, error) {
	var ret ByColumnsResult

	sheetsRet := c.compareSheets(input)

	for _, sheet := range sheetsRet.bothHave {
		sheetRet, err := c.compareSheet(input, sheet)
		if err != nil {
			return ByColumnsOutput{}, fmt.Errorf("compareSheet[%s]: %w", sheet, err)
		}

		if sheetRet.hasChanged {
			ret.sheets = append(ret.sheets, sheetRet)
			sheetsRet.addModified(sheet)
		} else {
			sheetsRet.addUnchanged(sheet)
		}
	}

	ret.catalog = sheetsRet

	return newByColumnsOutput(ret), nil
}

func (c *ByColumns) compareSheet(input Input, sheet string) (ByColumnsSheetResult, error) {
	var sheetRet ByColumnsSheetResult
	sheetRet.sheetName = sheet

	sheetData, err := input.sheetColsData(sheet)
	if err != nil {
		return sheetRet, fmt.Errorf("sheetColsData: %w", err)
	}

	columnsRet, hasChanged := c.columns.compare(sheetData.base, sheetData.compare)
	sheetRet.hasChanged = hasChanged
	sheetRet.columns = columnsRet

	return sheetRet, nil
}
