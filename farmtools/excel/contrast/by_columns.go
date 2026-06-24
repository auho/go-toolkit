package contrast

import "fmt"

type ByColumns struct {
	excel

	columns *ByColumnComparator
}

func NewByColumns() *ByColumns {
	return &ByColumns{
		columns: &ByColumnComparator{},
	}
}

// CompareFilePath Default comparison method, using default mapping
func (c *ByColumns) CompareFilePath(input InputFilePath) (ByColumnsOutput, error) {
	return c.CompareFilePathWithMapping(input, nil)
}

// CompareFilePathWithMapping Comparison method using specified mapping
func (c *ByColumns) CompareFilePathWithMapping(input InputFilePath, mapping SheetMapping) (ByColumnsOutput, error) {
	var ret ByColumnsOutput

	inputExcelFile, err := c.inputOpenFile(input)
	if err != nil {
		return ret, fmt.Errorf("inputOpenFile: %w", err)
	}

	defer func() { _ = inputExcelFile.close() }()

	ret, err = c.CompareWithMapping(inputExcelFile, mapping)
	if err != nil {
		return ret, fmt.Errorf("compare: %w", err)
	}

	return ret, nil
}

// Compare Default comparison method, using default mapping
func (c *ByColumns) Compare(input Input) (ByColumnsOutput, error) {
	return c.CompareWithMapping(input, nil)
}

// CompareWithMapping Comparison method using specified mapping
func (c *ByColumns) CompareWithMapping(input Input, mapping SheetMapping) (ByColumnsOutput, error) {
	var ret ByColumnsResult

	sheetsRet := c.compareSheets(input, mapping)

	// Check if no mappings found
	if len(sheetsRet.mappings) == 0 {
		return ByColumnsOutput{}, fmt.Errorf("no sheets to compare")
	}

	// Use mappings from CatalogResult
	for _, item := range sheetsRet.mappings {
		baseSheet := item.BaseSheet
		compareSheet := item.CompareSheet
		sheetRet, err := c.compareSheet(input, baseSheet, compareSheet)
		if err != nil {
			return ByColumnsOutput{}, fmt.Errorf("compareSheet[%s,%s]: %w", baseSheet, compareSheet, err)
		}

		if sheetRet.hasChanged {
			ret.sheets = append(ret.sheets, sheetRet)
			sheetsRet.addModified(baseSheet)
		} else {
			sheetsRet.addUnchanged(baseSheet)
		}
	}

	ret.catalog = sheetsRet

	return newByColumnsOutput(ret), nil
}

func (c *ByColumns) compareSheet(input Input, baseSheet, compareSheet string) (ByColumnsSheetResult, error) {
	var sheetRet ByColumnsSheetResult
	sheetRet.sheetName = baseSheet

	sheetData, err := input.sheetColsData(baseSheet, compareSheet)
	if err != nil {
		return sheetRet, fmt.Errorf("sheetColsData: %w", err)
	}

	columnsRet, hasChanged := c.columns.compare(sheetData.base, sheetData.compare)
	sheetRet.hasChanged = hasChanged
	sheetRet.columns = columnsRet

	return sheetRet, nil
}
