package contrast

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

type output struct {
	by

	cellStyles map[int]int
}

func (o *output) setCatalogResult(f *excelize.File, ret CatalogResult) error {
	sheetName := fmt.Sprint("__CATALOG__")

	_, err := f.NewSheet(sheetName)
	if err != nil {
		return fmt.Errorf("NewSheet[%s]: %v", sheetName, err)
	}

	var rowIndex int

	// unchanged
	rowIndex, err = o.setCatalog(f, sheetName, rowIndex, ret.unchanged)
	if err != nil {
		return fmt.Errorf("setCatalog unchanged: %w", err)
	}

	// modified
	rowIndex, err = o.setCatalog(f, sheetName, rowIndex, ret.modified)
	if err != nil {
		return fmt.Errorf("setCatalog modified: %w", err)
	}

	// added
	rowIndex, err = o.setCatalog(f, sheetName, rowIndex, ret.added)
	if err != nil {
		return fmt.Errorf("setCatalog added: %w", err)
	}

	// deleted
	rowIndex, err = o.setCatalog(f, sheetName, rowIndex, ret.deleted)
	if err != nil {
		return fmt.Errorf("setCatalog deleted: %w", err)
	}

	return nil
}

// int: next row index
func (o *output) setCatalog(f *excelize.File, sheetName string, rowIndex int, ret []CatalogItemResult) (int, error) {
	var cellsRet []CellResult

	var rowNo = o.indexToNo(rowIndex)

	for cellIndex, _ir := range ret {
		cellsRet = append(cellsRet, CellResult{
			row:    rowNo,
			col:    o.indexToNo(cellIndex),
			action: _ir.action,
			value:  _ir.sheetName,
		})
	}

	rowIndex++

	err := o.setCellsResult(f, sheetName, cellsRet)
	if err != nil {
		return 0, fmt.Errorf("setCellsResult: %w", err)
	}

	return rowIndex, nil
}

func (o *output) setCellsResult(f *excelize.File, sheetName string, cellsRet []CellResult) error {
	var err error

	for _, cellRet := range cellsRet {
		err = o.setCellResult(f, sheetName, cellRet)
		if err != nil {
			return fmt.Errorf("setCellResult: %w", err)
		}
	}

	return nil
}

func (o *output) setCellResult(f *excelize.File, sheetName string, cellRet CellResult) error {
	cell, err := excelize.CoordinatesToCellName(cellRet.col, cellRet.row)
	if err != nil {
		return fmt.Errorf("%+v to cell name %w", cellRet, err)
	}

	err = f.SetCellValue(sheetName, cell, cellRet.value)
	if err != nil {
		return fmt.Errorf("%#v set cell value %w", cellRet, err)
	}

	err = f.SetCellStyle(sheetName, cell, cell, o.cellStyles[cellRet.action])
	if err != nil {
		return fmt.Errorf("%#v set cell style %w", cellRet, err)
	}

	return nil
}

func (o *output) setRowsData(f *excelize.File, sheetName string, action int, rows []RowData) error {
	for _, row := range rows {
		err := o.setRow(f, sheetName, action, row.row, row.data)
		if err != nil {
			return fmt.Errorf("setRow[%d] %#v: %w", row.row, row, err)
		}
	}

	return nil
}

func (o *output) setRow(f *excelize.File, sheetName string, action, rowNo int, data []string) error {
	startCell, err := excelize.CoordinatesToCellName(1, rowNo)
	if err != nil {
		return fmt.Errorf("%d:1 to cell name: %w", rowNo, err)
	}

	// set values
	err = f.SetSheetRow(sheetName, startCell, &data)
	if err != nil {
		return fmt.Errorf("%d set row %#v: %w", rowNo, data, err)
	}

	// set style for entire row at once
	_len := len(data)
	if _len > 0 {
		endCell, err := excelize.CoordinatesToCellName(_len, rowNo)
		if err != nil {
			return fmt.Errorf("%d:%d to cell name %w", rowNo, _len, err)
		}

		err = f.SetCellStyle(sheetName, startCell, endCell, o.cellStyles[action])
		if err != nil {
			return fmt.Errorf("set row style %w", err)
		}
	}

	return nil
}

func (o *output) setColumnsData(f *excelize.File, sheetName string, action int, cols []ColData) error {
	for _, col := range cols {
		err := o.setColumn(f, sheetName, action, col.col, col.data)
		if err != nil {
			return fmt.Errorf("setColumn[%d] %#v: %w", col.col, col, err)
		}
	}

	return nil
}

func (o *output) setColumn(f *excelize.File, sheetName string, action, colNo int, data []string) error {
	startCell, err := excelize.CoordinatesToCellName(colNo, 1)
	if err != nil {
		return fmt.Errorf("%d:1 to cell name: %w", colNo, err)
	}

	// set values
	err = f.SetSheetCol(sheetName, startCell, &data)
	if err != nil {
		return fmt.Errorf("set column %#v: %w", data, err)
	}

	// set style
	_len := len(data)
	if _len > 0 {
		endCell, err := excelize.CoordinatesToCellName(colNo, _len)
		if err != nil {
			return fmt.Errorf("%d:%d to cell name %w", _len, colNo, err)
		}

		err = f.SetCellStyle(sheetName, startCell, endCell, o.cellStyles[action])
		if err != nil {
			return fmt.Errorf("set column style %w", err)
		}
	}

	return nil
}

func (o *output) buildStyles(f *excelize.File) error {
	var err error

	o.cellStyles = make(map[int]int)

	o.cellStyles[actionAdd], err = o.buildStyle(f, "70AD47") // green
	if err != nil {
		return err
	}

	o.cellStyles[actionModify], err = o.buildStyle(f, "FFC000") // orange
	if err != nil {
		return err
	}

	o.cellStyles[actionDelete], err = o.buildStyle(f, "A6A6A6") // gray
	if err != nil {
		return err
	}

	return nil
}

func (o *output) buildStyle(f *excelize.File, color string) (int, error) {
	return f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "gradient",
			Color:   []string{color, color},
			Shading: 0,
		},
	})
}
