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

// CompareFilePath Default comparison method, using default mapping
func (r *ByRows) CompareFilePath(input InputFilePath) (ByRowsOutput, error) {
	return r.CompareFilePathWithMapping(input, nil)
}

// CompareFilePathWithMapping Comparison method using specified mapping
func (r *ByRows) CompareFilePathWithMapping(input InputFilePath, mapping SheetMapping) (ByRowsOutput, error) {
	var ret ByRowsOutput

	inputExcelFile, err := r.inputOpenFile(input)
	if err != nil {
		return ret, fmt.Errorf("inputOpenFile: %w", err)
	}

	defer func() { _ = inputExcelFile.close() }()

	ret, err = r.CompareWithMapping(inputExcelFile, mapping)
	if err != nil {
		return ret, fmt.Errorf("compare: %w", err)
	}

	return ret, nil
}

// Compare Default comparison method, using default mapping
func (r *ByRows) Compare(input Input) (ByRowsOutput, error) {
	return r.CompareWithMapping(input, nil)
}

// CompareWithMapping Comparison method using specified mapping
func (r *ByRows) CompareWithMapping(input Input, mapping SheetMapping) (ByRowsOutput, error) {
	var ret ByRowsResult

	sheetsRet := r.compareSheets(input, mapping)

	// Check if no mappings found
	if len(sheetsRet.mappings) == 0 {
		return ByRowsOutput{}, fmt.Errorf("no sheets to compare")
	}

	// Use mappings from CatalogResult
	for _, item := range sheetsRet.mappings {
		baseSheet := item.BaseSheet
		compareSheet := item.CompareSheet
		sheetRet, err := r.compareSheet(input, baseSheet, compareSheet)
		if err != nil {
			return ByRowsOutput{}, fmt.Errorf("compareSheet[%s,%s]: %w", baseSheet, compareSheet, err)
		}

		if sheetRet.hasChanged {
			ret.sheets = append(ret.sheets, sheetRet)
			sheetsRet.addModified(baseSheet)
		} else {
			sheetsRet.addUnchanged(baseSheet)
		}
	}

	ret.catalog = sheetsRet

	return newByRowsOutput(ret), nil
}

func (r *ByRows) compareSheet(input Input, baseSheet, compareSheet string) (ByRowsSheetResult, error) {
	var sheetRet ByRowsSheetResult
	sheetRet.sheetName = baseSheet

	sheetData, err := input.sheetRowsData(baseSheet, compareSheet)
	if err != nil {
		return sheetRet, fmt.Errorf("sheetRowsData: %w", err)
	}

	rowsRet, hasChanged := r.rows.compare(sheetData.base, sheetData.compare)
	sheetRet.hasChanged = hasChanged
	sheetRet.rows = rowsRet

	return sheetRet, nil
}
