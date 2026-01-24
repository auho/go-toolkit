package contrast

import (
	"fmt"
	"testing"
	"time"
)

func TestByRows_CompareFilePath(t *testing.T) {
	input := genInput(t)

	byRows := NewByRows()
	_output, err := byRows.Compare(input)
	if err != nil {
		t.Fatal(err)
	}

	err = _output.SaveAs(fmt.Sprintf("./testdata/test.rows.%s.xlsx", time.Now().Format("20060102.150405")))
	if err != nil {
		t.Fatal(err)
	}
}
