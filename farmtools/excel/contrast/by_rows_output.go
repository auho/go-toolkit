package contrast

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

type ByRowsOutput struct {
	output

	ret ByRowsResult
}

func newByRowsOutput(ret ByRowsResult) ByRowsOutput {
	return ByRowsOutput{
		ret: ret,
	}
}

func (o *ByRowsOutput) SaveAs(filePath string) error {
	excelFile := excelize.NewFile()

	defer func() {
		_ = excelFile.Close()
	}()

	err := o.build(excelFile)
	if err != nil {
		return fmt.Errorf("build: %w", err)
	}

	err = excelFile.SaveAs(filePath)
	if err != nil {
		return fmt.Errorf("SaveAs: %w", err)
	}

	return nil
}

func (o *ByRowsOutput) Save(f *excelize.File) error {
	err := o.build(f)
	if err != nil {
		return fmt.Errorf("build: %w", err)
	}

	err = f.Save()
	if err != nil {
		return fmt.Errorf("save: %w", err)
	}

	return nil
}
func (o *ByRowsOutput) build(f *excelize.File) error {
	err := o.buildStyles(f)
	if err != nil {
		return fmt.Errorf("buildStyles: %w", err)
	}

	err = o.buildResult(f, o.ret)
	if err != nil {
		return fmt.Errorf("buildResult: %w", err)
	}

	return nil
}

func (o *ByRowsOutput) buildResult(f *excelize.File, ret ByRowsResult) error {
	err := o.buildStyles(f)
	if err != nil {
		return fmt.Errorf("buildStyles: %w", err)
	}

	err = o.buildCatalog(f, ret.catalog)
	if err != nil {
		return fmt.Errorf("buildCatalog: %w", err)
	}

	err = o.buildSheets(f, ret.sheets)
	if err != nil {
		return fmt.Errorf("buildSheets: %w", err)
	}

	return nil
}

func (o *ByRowsOutput) buildCatalog(f *excelize.File, ret CatalogResult) error {
	err := o.setCatalogResult(f, ret)
	if err != nil {
		return fmt.Errorf("setCatalogResult: %w", err)
	}

	return nil
}

func (o *ByRowsOutput) buildSheets(f *excelize.File, ret []ByRowsSheetResult) error {
	var err error

	for _, sheetRet := range ret {
		err = o.buildSheet(f, sheetRet)
		if err != nil {
			return fmt.Errorf("buildSheet[%s]: %w", sheetRet.sheetName, err)
		}
	}

	return nil
}

func (o *ByRowsOutput) buildSheet(f *excelize.File, ret ByRowsSheetResult) error {
	_, err := f.NewSheet(ret.sheetName)
	if err != nil {
		return fmt.Errorf("NewSheet: %w", err)
	}

	// modified
	err = o.setCellsResult(f, ret.sheetName, ret.rows.modifiedToCellResult())
	if err != nil {
		return fmt.Errorf("set modify cell: %w", err)
	}

	// added
	err = o.setRowsData(f, ret.sheetName, actionAdd, ret.rows.added)
	if err != nil {
		return fmt.Errorf("set added rows: %w", err)
	}

	// deleted
	err = o.setRowsData(f, ret.sheetName, actionDelete, ret.rows.deleted)
	if err != nil {
		return fmt.Errorf("set deleted rows: %w", err)
	}

	return nil
}
