package analysis

type Analysis struct {
	Table      *Table
	FieldsName []string
	Columns    map[string]Column
}
