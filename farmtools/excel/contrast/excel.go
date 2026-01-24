package contrast

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

type excel struct{}

func (e *excel) compareSheets(input Input) CatalogResult {
	var catalogRet CatalogResult

	input.initSheets()

	baseSheets := input.getSheetList(input.Base)
	compareSheets := input.getSheetList(input.Compare)

	var baseMap = make(map[string]struct{})
	for _, sheet := range baseSheets {
		baseMap[sheet] = struct{}{}
	}

	var compareMap = make(map[string]struct{})
	for _, sheet := range compareSheets {
		compareMap[sheet] = struct{}{}
	}

	// base
	for _, sheet := range baseSheets {
		if _, ok := compareMap[sheet]; !ok {
			catalogRet.addDeleted(sheet)
		} else {
			catalogRet.addBothHave(sheet)
		}
	}

	// compare
	for _, sheet := range compareSheets {
		if _, ok := baseMap[sheet]; !ok {
			catalogRet.addAdded(sheet)
		}
	}

	return catalogRet
}

func (e *excel) inputOpenFile(input InputFilePath) (Input, error) {
	var err error
	var inputExcelFile Input

	inputExcelFile.Sheets = input.Sheets

	inputExcelFile.Base, err = excelize.OpenFile(input.Base)
	if err != nil {
		return inputExcelFile, fmt.Errorf("open excel base file: %w", err)
	}

	inputExcelFile.Compare, err = excelize.OpenFile(input.Compare)
	if err != nil {
		return inputExcelFile, fmt.Errorf("open compare excel file: %w", err)
	}

	return inputExcelFile, nil
}
