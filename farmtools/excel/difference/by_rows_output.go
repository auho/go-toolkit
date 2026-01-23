package difference

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

type ByRowsOutput struct {
	output
}

func (o *ByRowsOutput) ToSheet(f *excelize.File, sheetName string, byRowsRet ByRowResult) error {
	err := o.saveResult(f, sheetName, byRowsRet)
	if err != nil {
		return fmt.Errorf("%s to sheet %w", sheetName, err)
	}

	return nil
}

func (o *ByRowsOutput) saveResult(f *excelize.File, sheetName string, byRowsRet ByRowResult) error {
	err := o.buildStyles(f)
	if err != nil {
		return fmt.Errorf("build style %w", err)
	}

	err = o.buildSheet(f, sheetName, byRowsRet)
	if err != nil {
		return fmt.Errorf("build sheet %w", err)
	}

	return nil
}

func (o *ByRowsOutput) buildSheet(f *excelize.File, sheetName string, byRowsRet ByRowResult) error {
	_, err := f.NewSheet(sheetName)
	if err != nil {
		return fmt.Errorf("new sheet %s %w", sheetName, err)
	}

	// modify
	err = o.setCellsResult(f, sheetName, byRowsRet.modifiedToCellResult())
	if err != nil {
		return fmt.Errorf("set modify cell %w", err)
	}

	// added
	err = o.setRowsData(f, sheetName, actionAdd, byRowsRet.addedRows)
	if err != nil {
		return fmt.Errorf("set added rows %w", err)
	}

	// deleted
	err = o.setRowsData(f, sheetName, actionDelete, byRowsRet.deletedRows)
	if err != nil {
		return fmt.Errorf("set deleted rows %w", err)
	}

	return nil
}
