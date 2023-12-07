package analysis

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
