package analysis

import (
	"fmt"
)

type Analysis struct {
	Table      *Table
	FieldsName []string
	Columns    map[string]Column
}

func NewAnalysis() *Analysis {
	a := &Analysis{}
	a.Columns = make(map[string]Column)

	return a
}

func (a *Analysis) ToStrings() []string {
	var ss []string

	ss = append(ss, fmt.Sprintf("table[%d]: %d", a.Table.Amount, a.Table.Amount))

	for _, column := range a.Columns {
		_title := fmt.Sprintf("  %s [%s]", column.Column.Name, column.Column.FieldType)

		noRow := ""
		if column.Amount == 0 {
			noRow = fmt.Sprintf(" 0 ⚠️")
		}

		ss = append(ss, fmt.Sprintf("%-30s: %d;  %s", _title, column.Amount, noRow))
	}

	return ss
}
