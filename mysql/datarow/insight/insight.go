package insight

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	simpleDb "github.com/auho/go-simple-db/v2"
	"github.com/auho/go-toolkit/mysql/datarow/insight/analysis"
	"github.com/auho/go-toolkit/mysql/schema"
)

func Explore(db *simpleDb.SimpleDB, table string) (*analysis.Analysis, error) {
	return (&Insight{}).Explore(db, table)
}

type Insight struct{}

func (i *Insight) Explore(db *simpleDb.SimpleDB, table string) (*analysis.Analysis, error) {
	return i.analyse(db, table)
}

func (i *Insight) analyse(db *simpleDb.SimpleDB, table string) (*analysis.Analysis, error) {
	_cs, err := db.GetTableColumnsSchema(table)
	if err != nil {
		return nil, err
	}

	tableAly, err := i.analyseTable(db, table)
	if err != nil {
		return nil, err
	}

	columnsSchema := schema.NewColumnsFromSimpleDb(_cs)
	columnsAly, err := i.analyseColumns(db, tableAly, columnsSchema)
	if err != nil {
		return nil, err
	}
	a := analysis.NewAnalysis()
	a.Table = tableAly

	var fields []string
	for _, _ca := range columnsAly {
		a.Columns[_ca.Column.Name] = _ca
		fields = append(fields, _ca.Column.Name)
	}

	sort.Slice(fields, func(i, j int) bool {
		return fields[i] < fields[j]
	})

	a.FieldsName = fields

	return a, nil
}

func (i *Insight) analyseTable(db *simpleDb.SimpleDB, table string) (*analysis.Table, error) {
	amount, err := db.TableAmount(table)
	if err != nil {
		return nil, err
	}

	return &analysis.Table{
		Table:  schema.Table{Name: table},
		Amount: amount,
	}, err
}

func (i *Insight) analyseColumns(db *simpleDb.SimpleDB, tableAly *analysis.Table, columns schema.Columns) ([]analysis.Column, error) {
	var fields []string

	fields = append(fields, "COUNT(*) AS 'amount'")
	for _, column := range columns {
		switch column.DataType {
		case schema.DataTypeInt, schema.DataTypeFloat:
			fields = append(fields,
				fmt.Sprintf("SUM(IF(`%s` = 0, 1, 0)) AS `%s_empty`", column.Name, column.Name),
				fmt.Sprintf("SUM(IF(`%s` IS NULL, 1, 0)) AS `%s_null`", column.Name, column.Name),
			)
		case schema.DataTypeString:
			fields = append(fields,
				fmt.Sprintf("SUM(IF(`%s` = '', 1, 0)) AS `%s_empty`", column.Name, column.Name),
				fmt.Sprintf("SUM(IF(`%s` IS NULL, 1, 0)) AS `%s_null`", column.Name, column.Name),
			)
		default:

		}
	}

	var retRows []map[string]any
	sql := fmt.Sprintf("SELECT %s FROM `%s`", strings.Join(fields, ","), tableAly.Table.Name)
	err := db.Raw(sql).Scan(&retRows).Error
	if err != nil {
		return nil, err
	}

	ret := retRows[0]

	var columnsAly []analysis.Column
	for _, column := range columns {
		_ca := analysis.Column{
			Column: column,
			Amount: tableAly.Amount,
		}

		if v, ok := ret[column.Name+"_empty"]; ok {
			_v, err1 := strconv.Atoi(v.(string))
			if err1 != nil {
				return nil, err1
			}
			_ca.Empty = _v
		}

		if v, ok := ret[column.Name+"_null"]; ok {
			_v, err1 := strconv.Atoi(v.(string))
			if err1 != nil {
				return nil, err1
			}
			_ca.Null = _v
		}

		columnsAly = append(columnsAly, _ca)
	}

	return columnsAly, nil
}
