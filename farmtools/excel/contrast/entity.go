package contrast

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

//
// row col 从 1 开始
//

type InputFilePath struct {
	base    string
	compare string
}

type Input struct {
	base    *excelize.File
	compare *excelize.File
}

func (i *Input) close() error {
	_ = i.base.Close()
	_ = i.compare.Close()

	return nil
}

func (i *Input) sheetData(sheet string) (SheetData, error) {
	var err error

	sd := SheetData{
		sheetName: sheet,
	}

	sd.base, err = i.base.GetRows(sheet)
	if err != nil {
		return sd, fmt.Errorf("base GetRows: %w", err)
	}

	sd.compare, err = i.compare.GetRows(sheet)
	if err != nil {
		return sd, fmt.Errorf("compare GetRows: %w", err)
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
