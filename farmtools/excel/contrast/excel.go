package contrast

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

type excel struct{}

// Compare sheets between base and compare files
func (e *excel) compareSheets(input Input, mapping SheetMapping) CatalogResult {
	var catalogRet CatalogResult

	// Get all sheet lists
	baseSheets := input.Base.GetSheetList()
	compareSheets := input.Compare.GetSheetList()

	// Get sheet mappings for comparison
	var mappingItems []SheetMappingItem

	if mapping != nil {
		mappingItems = mapping.GetMappings(baseSheets, compareSheets)
	} else {
		// Use default mapping
		defaultMapping := &DefaultMapping{}
		mappingItems = defaultMapping.GetMappings(baseSheets, compareSheets)
	}

	// Build compare file sheet map for existence check
	compareMap := make(map[string]struct{})
	for _, sheet := range compareSheets {
		compareMap[sheet] = struct{}{}
	}

	// Build base file sheet map for existence check
	baseMap := make(map[string]struct{})
	for _, sheet := range baseSheets {
		baseMap[sheet] = struct{}{}
	}

	// Process mapped sheets
	mappedBaseSheets := make(map[string]struct{})
	mappedCompareSheets := make(map[string]struct{})

	for _, item := range mappingItems {
		catalogRet.addBothHave(item.BaseSheet)
		mappedBaseSheets[item.BaseSheet] = struct{}{}
		mappedCompareSheets[item.CompareSheet] = struct{}{}
	}

	// Store sheet mappings
	catalogRet.mappings = mappingItems

	// Process unmapped sheets in base file (deleted)
	for _, sheet := range baseSheets {
		if _, ok := mappedBaseSheets[sheet]; !ok {
			catalogRet.addDeleted(sheet)
		}
	}

	// Process unmapped sheets in compare file (added)
	for _, sheet := range compareSheets {
		if _, ok := mappedCompareSheets[sheet]; !ok {
			catalogRet.addAdded(sheet)
		}
	}

	return catalogRet
}

// Open Excel files for comparison
func (e *excel) inputOpenFile(input InputFilePath) (Input, error) {
	var err error
	var inputExcelFile Input

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
