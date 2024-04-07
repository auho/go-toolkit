package difference

import (
	"fmt"
	"github.com/xuri/excelize/v2"
)

func DataDiffToSheet(excelFile *excelize.File, sheetName string, baseData, diffData [][]string) error {
	data := &Data{
		baseData: baseData,
		diffData: diffData,
	}

	sheet, err := data.diff()
	if err != nil {
		return fmt.Errorf("data init %w", err)
	}

	err = sheet.ToSheet(excelFile, sheetName)
	if err != nil {
		return fmt.Errorf("to sheet %w", err)
	}

	return nil
}
