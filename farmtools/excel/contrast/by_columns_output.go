package contrast

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

type ByColumnsOutput struct {
	output

	ret ByColumnsResult
}

func newByColumnsOutput(ret ByColumnsResult) ByColumnsOutput {
	return ByColumnsOutput{
		ret: ret,
	}
}

func (o *ByColumnsOutput) SaveAs(filePath string) error {
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

func (o *ByColumnsOutput) Save(f *excelize.File) error {
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
func (o *ByColumnsOutput) build(f *excelize.File) error {
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

func (o *ByColumnsOutput) buildResult(f *excelize.File, ret ByColumnsResult) error {
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

func (o *ByColumnsOutput) buildCatalog(f *excelize.File, ret CatalogResult) error {
	err := o.setCatalogResult(f, ret)
	if err != nil {
		return fmt.Errorf("setCatalogResult: %w", err)
	}

	return nil
}

func (o *ByColumnsOutput) buildSheets(f *excelize.File, ret []ByColumnsSheetResult) error {
	var err error

	for _, sheetRet := range ret {
		err = o.buildSheet(f, sheetRet)
		if err != nil {
			return fmt.Errorf("buildSheet[%s]: %w", sheetRet.sheetName, err)
		}
	}

	return nil
}

func (o *ByColumnsOutput) buildSheet(f *excelize.File, ret ByColumnsSheetResult) error {
	_, err := f.NewSheet(ret.sheetName)
	if err != nil {
		return fmt.Errorf("NewSheet: %w", err)
	}

	// modified
	err = o.setCellsResult(f, ret.sheetName, ret.columns.modifiedToCellResult())
	if err != nil {
		return fmt.Errorf("set modify cell: %w", err)
	}

	// added
	err = o.setColumnsData(f, ret.sheetName, actionAdd, ret.columns.added)
	if err != nil {
		return fmt.Errorf("set added columns: %w", err)
	}

	// deleted
	err = o.setColumnsData(f, ret.sheetName, actionDelete, ret.columns.deleted)
	if err != nil {
		return fmt.Errorf("set deleted columns: %w", err)
	}

	return nil
}
