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

func TestByRows_CompareWithDefaultMapping(t *testing.T) {
	input := genInput(t)

	byRows := NewByRows()
	// Use default mapping (no include or exclude)
	_output, err := byRows.CompareWithMapping(input, nil)
	if err != nil {
		t.Fatal(err)
	}

	err = _output.SaveAs(fmt.Sprintf("./testdata/test.rows.default.%s.xlsx", time.Now().Format("20060102.150405")))
	if err != nil {
		t.Fatal(err)
	}
}

func TestByRows_CompareWithDefaultMappingIncludeExclude(t *testing.T) {
	input := genInput(t)

	byRows := NewByRows()
	// Use default mapping, include and exclude some sheets
	mapping := &DefaultMapping{
		IncludeSheets: []string{"Sheet1"}, // only include Sheet1
		ExcludeSheets: []string{},          // exclude no sheets
	}
	_output, err := byRows.CompareWithMapping(input, mapping)
	if err != nil {
		t.Fatal(err)
	}

	err = _output.SaveAs(fmt.Sprintf("./testdata/test.rows.default.include.%s.xlsx", time.Now().Format("20060102.150405")))
	if err != nil {
		t.Fatal(err)
	}
}

func TestByRows_CompareWithNameMapping(t *testing.T) {
	input := genInput(t)

	byRows := NewByRows()
	// Use name mapping, assume Sheet1 in base file corresponds to Sheet1 in compare file
	mapping := &NameMapping{
		Mappings: map[string]string{
			"Sheet1": "Sheet1",
		},
	}
	_output, err := byRows.CompareWithMapping(input, mapping)
	if err != nil {
		t.Fatal(err)
	}

	err = _output.SaveAs(fmt.Sprintf("./testdata/test.rows.name.%s.xlsx", time.Now().Format("20060102.150405")))
	if err != nil {
		t.Fatal(err)
	}
}

func TestByRows_CompareWithIndexMapping(t *testing.T) {
	input := genInput(t)

	byRows := NewByRows()
	// Use index mapping, assume first sheet in base file corresponds to first sheet in compare file
	mapping := &IndexMapping{
		Mappings: map[int]int{
			1: 1, // index starts from 1
		},
	}
	_output, err := byRows.CompareWithMapping(input, mapping)
	if err != nil {
		t.Fatal(err)
	}

	err = _output.SaveAs(fmt.Sprintf("./testdata/test.rows.index.%s.xlsx", time.Now().Format("20060102.150405")))
	if err != nil {
		t.Fatal(err)
	}
}
