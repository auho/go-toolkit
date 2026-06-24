package contrast

import (
	"bytes"
	_ "embed"
	"testing"

	"github.com/xuri/excelize/v2"
)

//go:embed testdata/base.xlsx
var baseFile []byte

//go:embed testdata/compare.xlsx
var compareFile []byte

func genInput(t *testing.T) Input {
	var input Input
	var err error

	input.Base, err = excelize.OpenReader(bytes.NewReader(baseFile))
	if err != nil {
		t.Fatal(err)
	}

	input.Compare, err = excelize.OpenReader(bytes.NewReader(compareFile))
	if err != nil {
		t.Fatal(err)
	}

	return input
}
