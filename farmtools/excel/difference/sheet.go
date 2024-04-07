package difference

import (
	"fmt"
	"github.com/xuri/excelize/v2"
)

type SheetDiff struct {
	data          Data
	addedRows     map[int][]string
	deletedRows   map[int][]string
	modifiedCells []CellDiff
	cellStyles    map[int]int
}

func (s *SheetDiff) ToSheet(f *excelize.File, sheetName string) error {
	err := s.toSheet(f, sheetName)
	if err != nil {
		return fmt.Errorf("%s to sheet %w", sheetName, err)
	}

	return nil
}

func (s *SheetDiff) buildStyles(f *excelize.File) error {
	var err error

	s.cellStyles = make(map[int]int)

	s.cellStyles[actionAdd], err = s.buildStyle(f, "70AD47")
	if err != nil {
		return err
	}

	s.cellStyles[actionModify], err = s.buildStyle(f, "FFC000")
	if err != nil {
		return err
	}

	s.cellStyles[actionDelete], err = s.buildStyle(f, "A6A6A6")
	if err != nil {
		return err
	}

	return nil
}

func (s *SheetDiff) toSheet(f *excelize.File, sheetName string) error {
	err := s.buildStyles(f)
	if err != nil {
		return fmt.Errorf("build style %w", err)
	}

	err = s.buildSheet(f, sheetName)
	if err != nil {
		return fmt.Errorf("build sheet %w", err)
	}

	return nil
}

func (s *SheetDiff) buildSheet(f *excelize.File, sheetName string) error {
	_, err := f.NewSheet(sheetName)
	if err != nil {
		return fmt.Errorf("new sheet %s %w", sheetName, err)
	}

	// modify
	err = s.setCell(f, sheetName, s.modifiedCells)
	if err != nil {
		return fmt.Errorf("set modify cell %w", err)
	}

	// added
	err = s.setRows(f, sheetName, s.addedRows, actionAdd)
	if err != nil {
		return fmt.Errorf("set added rows %w", err)
	}

	// deleted
	err = s.setRows(f, sheetName, s.deletedRows, actionDelete)
	if err != nil {
		return fmt.Errorf("set deleted rows %w", err)
	}

	return nil
}

func (s *SheetDiff) setCell(f *excelize.File, sheetName string, slice []CellDiff) error {
	var err error
	var cell string
	for _, cellDiff := range slice {
		cell, err = excelize.CoordinatesToCellName(cellDiff.row, cellDiff.col)
		if err != nil {
			return fmt.Errorf("%+v to cell name %w", cellDiff, err)
		}

		err = f.SetCellValue(sheetName, cell, cellDiff.value)
		if err != nil {
			return fmt.Errorf("%#v set cell value %w", cellDiff, err)
		}

		err = f.SetCellStyle(sheetName, cell, cell, s.cellStyles[cellDiff.action])
		if err != nil {
			return fmt.Errorf("%#v set cell style %w", cellDiff, err)
		}
	}

	return nil
}

func (s *SheetDiff) setRows(f *excelize.File, sheetName string, m map[int][]string, action int) error {
	var cell string
	for rowNo, rowValues := range m {
		rowStartCell, err := excelize.CoordinatesToCellName(rowNo, 1)
		if err != nil {
			return fmt.Errorf("%d 1 to cell name %w", rowNo, err)
		}

		err = f.SetSheetRow(sheetName, rowStartCell, rowValues)
		if err != nil {
			return fmt.Errorf("%d set row %#v %w", rowNo, rowValues, err)
		}

		_len := len(rowValues)
		for i := 1; i <= _len; i++ {
			cell, err = excelize.CoordinatesToCellName(rowNo, i)
			if err != nil {
				return fmt.Errorf("%d:%d to cell name %w", rowNo, i, err)
			}

			err = f.SetCellStyle(sheetName, cell, cell, s.cellStyles[action])
			if err != nil {
				return fmt.Errorf("%d:%d set style %w", rowNo, i, err)
			}
		}
	}

	return nil
}

func (s *SheetDiff) buildStyle(f *excelize.File, color string) (int, error) {
	return f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "gradient",
			Color:   []string{color, color},
			Shading: 0,
		},
	})
}
