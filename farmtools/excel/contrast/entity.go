package contrast

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

//
// row col starts from 1
//

// SheetMappingItem represents a single sheet mapping
type SheetMappingItem struct {
	BaseSheet    string
	CompareSheet string
}

// SheetMapping interface defines how to map sheets for comparison
type SheetMapping interface {
	// GetMappings returns sheet mappings for comparison
	GetMappings(baseSheets, compareSheets []string) []SheetMappingItem
}

// Ensure all mapping types implement SheetMapping interface
var _ SheetMapping = (*DefaultMapping)(nil)
var _ SheetMapping = (*NameMapping)(nil)
var _ SheetMapping = (*IndexMapping)(nil)

// DefaultMapping default mapping, compares sheets with the same name, supports include and exclude
type DefaultMapping struct {
	IncludeSheets []string
	ExcludeSheets []string
}

// GetMappings implements GetMappings method for DefaultMapping
func (m *DefaultMapping) GetMappings(baseSheets, compareSheets []string) []SheetMappingItem {
	// First get all sheets with the same name
	baseSheetMap := make(map[string]struct{})
	for _, sheet := range baseSheets {
		baseSheetMap[sheet] = struct{}{}
	}

	compareSheetMap := make(map[string]struct{})
	for _, sheet := range compareSheets {
		compareSheetMap[sheet] = struct{}{}
	}

	// Find all sheets with the same name
	mappings := make(map[string]string)
	for _, sheet := range baseSheets {
		if _, ok := compareSheetMap[sheet]; ok {
			mappings[sheet] = sheet
		}
	}

	// Apply include rules
	if len(m.IncludeSheets) > 0 {
		includedMappings := make(map[string]string)
		for _, sheet := range m.IncludeSheets {
			if compareSheet, ok := mappings[sheet]; ok {
				includedMappings[sheet] = compareSheet
			}
		}
		mappings = includedMappings
	}

	// Apply exclude rules
	if len(m.ExcludeSheets) > 0 {
		excludeSet := make(map[string]struct{})
		for _, sheet := range m.ExcludeSheets {
			excludeSet[sheet] = struct{}{}
		}

		filteredMappings := make(map[string]string)
		for baseSheet, compareSheet := range mappings {
			if _, ok := excludeSet[baseSheet]; !ok {
				filteredMappings[baseSheet] = compareSheet
			}
		}
		mappings = filteredMappings
	}

	// Convert to slice
	result := make([]SheetMappingItem, 0, len(mappings))
	for _, baseSheet := range baseSheets {
		if compareSheet, ok := mappings[baseSheet]; ok {
			result = append(result, SheetMappingItem{BaseSheet: baseSheet, CompareSheet: compareSheet})
		}
	}

	return result
}

// NameMapping mapping by specified names
type NameMapping struct {
	Mappings map[string]string
}

// GetMappings implements GetMappings method for NameMapping
func (m *NameMapping) GetMappings(baseSheets, compareSheets []string) []SheetMappingItem {
	// Validate mappings
	compareSheetMap := make(map[string]struct{})
	for _, sheet := range compareSheets {
		compareSheetMap[sheet] = struct{}{}
	}

	baseSheetMap := make(map[string]struct{})
	for _, sheet := range baseSheets {
		baseSheetMap[sheet] = struct{}{}
	}

	// Filter valid mappings
	mappings := make(map[string]string)
	for baseSheet, compareSheet := range m.Mappings {
		if _, ok := baseSheetMap[baseSheet]; ok {
			if _, ok := compareSheetMap[compareSheet]; ok {
				mappings[baseSheet] = compareSheet
			}
		}
	}

	// Convert to slice, maintaining the order of base sheets
	result := make([]SheetMappingItem, 0, len(mappings))
	for _, baseSheet := range baseSheets {
		if compareSheet, ok := mappings[baseSheet]; ok {
			result = append(result, SheetMappingItem{BaseSheet: baseSheet, CompareSheet: compareSheet})
		}
	}

	return result
}

// IndexMapping mapping by specified indices
type IndexMapping struct {
	Mappings map[int]int
}

// GetMappings implements GetMappings method for IndexMapping
func (m *IndexMapping) GetMappings(baseSheets, compareSheets []string) []SheetMappingItem {
	// First create a map for quick lookup
	mappingMap := make(map[string]string)

	for baseIndex, compareIndex := range m.Mappings {
		// Validate indices (indices start from 1)
		if baseIndex >= 1 && baseIndex <= len(baseSheets) {
			if compareIndex >= 1 && compareIndex <= len(compareSheets) {
				baseSheet := baseSheets[baseIndex-1]          // convert to 0-based index
				compareSheet := compareSheets[compareIndex-1] // convert to 0-based index
				mappingMap[baseSheet] = compareSheet
			}
		}
	}

	// Convert to slice, maintaining the order of base sheets
	result := make([]SheetMappingItem, 0, len(mappingMap))
	for _, baseSheet := range baseSheets {
		if compareSheet, ok := mappingMap[baseSheet]; ok {
			result = append(result, SheetMappingItem{BaseSheet: baseSheet, CompareSheet: compareSheet})
		}
	}

	return result
}

type InputFilePath struct {
	Base    string
	Compare string
}

type Input struct {
	Base    *excelize.File
	Compare *excelize.File
}

func (i *Input) close() error {
	_ = i.Base.Close()
	_ = i.Compare.Close()

	return nil
}

func (i *Input) sheetRowsData(baseSheet, compareSheet string) (SheetData, error) {
	var err error

	sd := SheetData{
		sheetName: baseSheet,
	}

	sd.base, err = i.Base.GetRows(baseSheet)
	if err != nil {
		return sd, fmt.Errorf("base GetRows: %w", err)
	}

	sd.compare, err = i.Compare.GetRows(compareSheet)
	if err != nil {
		return sd, fmt.Errorf("compare GetRows: %w", err)
	}

	return sd, nil
}

func (i *Input) sheetColsData(baseSheet, compareSheet string) (SheetData, error) {
	var err error

	sd := SheetData{
		sheetName: baseSheet,
	}

	sd.base, err = i.Base.GetCols(baseSheet)
	if err != nil {
		return sd, fmt.Errorf("base GetCols: %w", err)
	}

	sd.compare, err = i.Compare.GetCols(compareSheet)
	if err != nil {
		return sd, fmt.Errorf("compare GetCols: %w", err)
	}

	return sd, nil
}

//
// sheet

type CatalogResult struct {
	bothHave  []string
	unchanged []CatalogItemResult
	modified  []CatalogItemResult
	added     []CatalogItemResult
	deleted   []CatalogItemResult
	mappings  []SheetMappingItem // store sheet mappings in the order of sheets in base file
}

func (r *CatalogResult) addBothHave(sheetName string) {
	r.bothHave = append(r.bothHave, sheetName)
}

func (r *CatalogResult) addUnchanged(sheetName string) {
	r.unchanged = append(r.unchanged, CatalogItemResult{actionUnchanged, sheetName})
}

func (r *CatalogResult) addModified(sheetName string) {
	r.modified = append(r.modified, CatalogItemResult{actionModify, sheetName})
}

func (r *CatalogResult) addAdded(sheetName string) {
	r.added = append(r.added, CatalogItemResult{actionAdd, sheetName})
}

func (r *CatalogResult) addDeleted(sheetName string) {
	r.deleted = append(r.deleted, CatalogItemResult{actionDelete, sheetName})
}

type CatalogItemResult struct {
	action    int
	sheetName string
}

//
// rows

type RowResult struct {
	row   int
	cells []RowCellResult
}

func (r *RowResult) toCell() []CellResult {
	var cells []CellResult
	for _, cell := range r.cells {
		cells = append(cells, CellResult{
			row:    r.row,
			col:    cell.col,
			action: cell.action,
			value:  cell.value,
		})
	}

	return cells
}

type RowCellResult struct {
	col    int
	action int
	value  string
}

// cell
//

type CellResult struct {
	row    int
	col    int
	action int
	value  string
}

// data
//

type SheetData struct {
	sheetName string
	base      [][]string
	compare   [][]string
}

type RowData struct {
	row  int
	data []string
}

type ColData struct {
	col  int
	data []string
}
