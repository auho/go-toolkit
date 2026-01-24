package contrast

import (
	"fmt"
	"testing"
	"time"
)

func TestByColumns_CompareFilePath(t *testing.T) {
	input := genInput(t)

	byColumns := NewByColumns()
	_output, err := byColumns.Compare(input)
	if err != nil {
		t.Fatal(err)
	}

	err = _output.SaveAs(fmt.Sprintf("./testdata/test.columns.%s.xlsx", time.Now().Format("20060102.150405")))
	if err != nil {
		t.Fatal(err)
	}
}
