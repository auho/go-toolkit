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

func TestByColumns_CompareWithDefaultMapping(t *testing.T) {
	input := genInput(t)

	byColumns := NewByColumns()
	// Use default mapping (no include or exclude)
	_output, err := byColumns.CompareWithMapping(input, nil)
	if err != nil {
		t.Fatal(err)
	}

	err = _output.SaveAs(fmt.Sprintf("./testdata/test.columns.default.%s.xlsx", time.Now().Format("20060102.150405")))
	if err != nil {
		t.Fatal(err)
	}
}

func TestByColumns_CompareWithDefaultMappingIncludeExclude(t *testing.T) {
	input := genInput(t)

	byColumns := NewByColumns()
	// Use default mapping, include and exclude some sheets
	mapping := &DefaultMapping{
		IncludeSheets: []string{"Sheet1"}, // only include Sheet1
		ExcludeSheets: []string{},          // exclude no sheets
	}
	_output, err := byColumns.CompareWithMapping(input, mapping)
	if err != nil {
		t.Fatal(err)
	}

	err = _output.SaveAs(fmt.Sprintf("./testdata/test.columns.default.include.%s.xlsx", time.Now().Format("20060102.150405")))
	if err != nil {
		t.Fatal(err)
	}
}

func TestByColumns_CompareWithNameMapping(t *testing.T) {
	input := genInput(t)

	byColumns := NewByColumns()
	// Use name mapping, assume Sheet1 in base file corresponds to Sheet1 in compare file
	mapping := &NameMapping{
		Mappings: map[string]string{
			"Sheet1": "Sheet1",
		},
	}
	_output, err := byColumns.CompareWithMapping(input, mapping)
	if err != nil {
		t.Fatal(err)
	}

	err = _output.SaveAs(fmt.Sprintf("./testdata/test.columns.name.%s.xlsx", time.Now().Format("20060102.150405")))
	if err != nil {
		t.Fatal(err)
	}
}

func TestByColumns_CompareWithIndexMapping(t *testing.T) {
	input := genInput(t)

	byColumns := NewByColumns()
	// Use index mapping, assume first sheet in base file corresponds to first sheet in compare file
	mapping := &IndexMapping{
		Mappings: map[int]int{
			1: 1, // index starts from 1
		},
	}
	_output, err := byColumns.CompareWithMapping(input, mapping)
	if err != nil {
		t.Fatal(err)
	}

	err = _output.SaveAs(fmt.Sprintf("./testdata/test.columns.index.%s.xlsx", time.Now().Format("20060102.150405")))
	if err != nil {
		t.Fatal(err)
	}
}
