package difference

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

type output struct {
	cellStyles map[int]int
}

func (o *output) setCellsResult(f *excelize.File, sheetName string, cellsRet []CellResult) error {
	var err error
	var cell string

	for _, cellRet := range cellsRet {
		cell, err = excelize.CoordinatesToCellName(cellRet.row, cellRet.col)
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
	}

	return nil
}

func (o *output) setRowsData(f *excelize.File, sheetName string, action int, rows []RowData) error {
	var cell string
	for rowNo, row := range rows {
		startCell, err := excelize.CoordinatesToCellName(rowNo, 1)
		if err != nil {
			return fmt.Errorf("%d 1 to cell name %w", rowNo, err)
		}

		// set value
		err = f.SetSheetRow(sheetName, startCell, row.data)
		if err != nil {
			return fmt.Errorf("%d set row %#v %w", rowNo, row, err)
		}

		// set style
		_len := len(row.data)
		for i := 1; i <= _len; i++ {
			cell, err = excelize.CoordinatesToCellName(rowNo, i)
			if err != nil {
				return fmt.Errorf("%d:%d to cell name %w", rowNo, i, err)
			}

			err = f.SetCellStyle(sheetName, cell, cell, o.cellStyles[action])
			if err != nil {
				return fmt.Errorf("%d:%d set style %w", rowNo, i, err)
			}
		}
	}

	return nil
}

func (o *output) buildStyles(f *excelize.File) error {
	var err error

	o.cellStyles = make(map[int]int)

	o.cellStyles[actionAdd], err = o.buildStyle(f, "70AD47") // 绿
	if err != nil {
		return err
	}

	o.cellStyles[actionModify], err = o.buildStyle(f, "FFC000") // 橙
	if err != nil {
		return err
	}

	o.cellStyles[actionDelete], err = o.buildStyle(f, "A6A6A6") // 灰
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
