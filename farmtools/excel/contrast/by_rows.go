package contrast

import "fmt"

type ByRows struct {
	excel

	rows *ByRowComparator
}

func NewByRows() *ByRows {
	return &ByRows{
		rows: &ByRowComparator{},
	}
}

func (r *ByRows) CompareFilePath(input InputFilePath) (ByRowsOutput, error) {
	var ret ByRowsOutput

	inputExcelFile, err := r.inputOpenFile(input)
	if err != nil {
		return ret, fmt.Errorf("inputOpenFile: %w", err)
	}

	defer func() { _ = inputExcelFile.close() }()

	ret, err = r.Compare(inputExcelFile)
	if err != nil {
		return ret, fmt.Errorf("compare: %w", err)
	}

	return ret, nil
}

func (r *ByRows) Compare(input Input) (ByRowsOutput, error) {
	var ret ByRowsResult

	sheetsRet := r.compareSheets(input)

	for _, sheet := range sheetsRet.bothHave {
		sheetRet, err := r.compareSheet(input, sheet)
		if err != nil {
			return ByRowsOutput{}, fmt.Errorf("compareSheet[%s]: %w", sheet, err)
		}

		if sheetRet.hasChanged {
			ret.sheets = append(ret.sheets, sheetRet)
			sheetsRet.addModified(sheet)
		} else {
			sheetsRet.addUnchanged(sheet)
		}
	}

	ret.catalog = sheetsRet

	return newByRowsOutput(ret), nil
}

func (r *ByRows) compareSheet(input Input, sheet string) (ByRowsSheetResult, error) {
	var sheetRet ByRowsSheetResult
	sheetRet.sheetName = sheet

	sheetData, err := input.sheetRowsData(sheet)
	if err != nil {
		return sheetRet, fmt.Errorf("sheetRowsData: %w", err)
	}

	rowsRet, hasChanged := r.rows.compare(sheetData.base, sheetData.compare)
	sheetRet.hasChanged = hasChanged
	sheetRet.rows = rowsRet

	return sheetRet, nil
}
