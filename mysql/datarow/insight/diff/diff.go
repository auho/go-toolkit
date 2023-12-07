package diff

import (
	"fmt"

	"github.com/auho/go-toolkit/mysql/datarow/insight/analysis"
)

func Diff(as ...*analysis.Analysis) *Differ {
	d := &Differ{}
	d.diff(as...)

	return d
}

type Differ struct {
	ss []string
}

func (d *Differ) DifferenceToString() []string {
	return d.ss
}

func (d *Differ) diff(as ...*analysis.Analysis) {
	var ss []string

	_las := as[0]
	_ras := as[1]

	_title := "table amount"
	if _las.Table.Amount == _ras.Table.Amount {
		ss = append(ss, d.success(fmt.Sprintf("%s: %d", _title, _las.Table.Amount)))
	} else {
		ss = append(ss, d.failure(fmt.Sprintf("%s[%d != %d]", _title, _las.Table.Amount, _ras.Table.Amount)))
	}

	for _, _lc := range _las.Columns {
		_title = fmt.Sprintf("  column[%s]", _lc.Column.Name)

		if _rc, ok := _ras.Columns[_lc.Column.Name]; ok {
			if _lc.Amount == _rc.Amount {
				ss = append(ss, d.success(fmt.Sprintf("%s: %d", _title, _lc.Amount)))
			} else {
				ss = append(ss, d.failure(fmt.Sprintf("%s[%d != %d]", _title, _lc.Amount, _rc.Amount)))
			}

			if _lc.Empty != _rc.Empty {
				ss = append(ss, d.failure(fmt.Sprintf("%s -> empty[%d != %d]", _title, _lc.Empty, _rc.Empty)))
			}

			if _lc.Null != _rc.Null {
				ss = append(ss, d.failure(fmt.Sprintf("%s -> null[%d != %d]", _title, _lc.Null, _rc.Null)))
			}

		} else {
			ss = append(ss, d.warningAndFailure(fmt.Sprintf("%s[%d != 0]", _title, _lc.Amount)))
		}
	}

	for _, _rc := range _ras.Columns {
		_title = fmt.Sprintf("  column[%s]", _rc.Column.Name)

		if _, ok := _las.Columns[_rc.Column.Name]; !ok {
			ss = append(ss, d.failureAndWarning(fmt.Sprintf("%s[0 != %d]", _title, _rc.Amount)))
		}
	}

	d.ss = ss
}

func (d *Differ) success(s string) string {
	return "✅  " + s
}

func (d *Differ) warning(s string) string {
	return "⚠️  " + s
}

func (d *Differ) failure(s string) string {
	return "❌  " + s
}

func (d *Differ) warningAndFailure(s string) string {
	return "⚠️❌" + s
}

func (d *Differ) failureAndWarning(s string) string {
	return "❌⚠️" + s
}
