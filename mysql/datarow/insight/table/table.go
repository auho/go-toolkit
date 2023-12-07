package table

import (
	"fmt"
	"strings"

	simpleDb "github.com/auho/go-simple-db/v2"
	analysis2 "github.com/auho/go-toolkit/mysql/datarow/insight/analysis"
	"github.com/auho/go-toolkit/mysql/schema"
)

type Table struct {
	Name string // table name
	DB   *simpleDb.SimpleDB
}

func (t *Table) analyse() (*analysis2.Analysis, error) {
	_cs, err := t.DB.GetTableColumnsSchema(t.Name)
	if err != nil {
		return nil, err
	}

	tableAly, err := t.analyseTable()
	if err != nil {
		return nil, err
	}

	columnsSchema := schema.NewColumnsFromSimpleDb(_cs)
	columnsAly, err := t.analyseColumns(tableAly, columnsSchema)
	if err != nil {
		return nil, err
	}
	a := &analysis2.Analysis{
		Table: tableAly,
	}

	for _, _ca := range columnsAly {
		a.Columns[_ca.Column.Name] = _ca
	}

	return a, nil
}

func (t *Table) analyseTable() (*analysis2.Table, error) {
	amount, err := t.DB.TableAmount(t.Name)
	if err != nil {
		return nil, err
	}

	return &analysis2.Table{
		Table:  schema.Table{Name: t.Name},
		Amount: amount,
	}, err
}

func (t *Table) analyseColumns(tableAly *analysis2.Table, columns schema.Columns) ([]analysis2.Column, error) {
	var fields []string

	fields = append(fields, "COUNT(*) AS 'amount'")
	for _, column := range columns {
		switch column.DataType {
		case schema.DataTypeInt, schema.DataTypeFloat:
			fields = append(fields,
				fmt.Sprintf("SUM(IF(`%s` = 0, 0, 1)) AS `%s_empty`", column.Name, column.Name),
			)
		case schema.DataTypeString:
			fields = append(fields,
				fmt.Sprintf("SUM(IF(`%s` = '', 0, 1)) AS `%s_empty`", column.Name, column.Name),
			)
		default:

		}
	}

	var retRows []map[string]int
	sql := fmt.Sprintf("SELECT %s FROM `%s`", strings.Join(fields, ","), t.Name)
	err := t.DB.GormDB().Raw(sql).Scan(&retRows).Error
	if err != nil {
		return nil, err
	}

	ret := retRows[0]

	var columnsAly []analysis2.Column
	for _, column := range columns {
		_ca := analysis2.Column{
			Column: column,
			Amount: tableAly.Amount,
		}

		if v, ok := ret[column.Name+"_empty"]; ok {
			_ca.Empty = v
		}

		columnsAly = append(columnsAly, _ca)
	}

	return columnsAly, nil
}
