package schema

import (
	"strings"

	"github.com/auho/go-simple-db/v2/schema"
)

type Column struct {
	Name      string
	FieldType FieldType
	DataType  DataType
}

type Columns []Column

func NewColumnsFromSimpleDb(columns []schema.Column) Columns {
	var cs Columns
	for _, _c := range columns {
		ft := FieldType(strings.ToLower(_c.FieldType))
		cs = append(cs, Column{
			Name:      _c.Name,
			FieldType: ft,
			DataType:  FileTypeToDataType(ft),
		})
	}

	return cs
}
