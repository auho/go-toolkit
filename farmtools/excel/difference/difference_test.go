package difference

import (
	"bytes"
	_ "embed"
	"github.com/xuri/excelize/v2"
	"testing"
)

//go:embed testdata/base.xlsx
var baseFile []byte

//go:embed testdata/compare.xlsx
var compareFile []byte

func TestDataDiffToSheet(t *testing.T) {
	baseExcle, err := excelize.OpenReader(bytes.NewReader(baseFile))
	if err != nil {
		t.Fatal(err)
	}

	baseData, err := baseExcle.GetRows("Sheet1")
	if err != nil {
		t.Fatal(err)
	}

	compareExcle, err := excelize.OpenReader(bytes.NewReader(compareFile))
	if err != nil {
		t.Fatal(err)
	}

	compareData, err := compareExcle.GetRows("Sheet1")
	if err != nil {
		t.Fatal(err)
	}

	data := &ByRows{
		baseData:    baseData,
		compareData: compareData,
	}

	ret, err := data.diff()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(ret)
}
