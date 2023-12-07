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
	if _las.Table.Amount != _ras.Table.Amount {
		ss = append(ss, d.failure(fmt.Sprintf("%s[%d != %d]", _title, _las.Table.Amount, _ras.Table.Amount)))
	} else {
		ss = append(ss, _title)
	}

	for _, _lc := range _las.Columns {
		_s := ""
		_isOk := true

		_title = fmt.Sprintf("  column[%s]", _lc.Column.Name)

		if _rc, ok := _ras.Columns[_lc.Column.Name]; ok {
			if _lc.Empty != _rc.Empty {
				_s = d.failure(fmt.Sprintf("%s empty[%d != %d]", _title, _lc.Empty, _rc.Empty))
				_isOk = false
			}

			if _lc.Null != _rc.Null {
				_s = d.failure(fmt.Sprintf("%s null[%d != %d]", _title, _lc.Empty, _rc.Empty))
				_isOk = false
			}

			if _isOk {
				_s = d.success(fmt.Sprintf("%s", _lc.Column.Name))
			}
		} else {
			if _lc.Empty != 0 {
				_s = d.warningAndFailure(fmt.Sprintf("%s empty[%d != 0]", _title, _lc.Empty))
			}

			if _lc.Null != 0 {
				_s = d.warningAndFailure(fmt.Sprintf("%s null[%d != 0]", _title, _lc.Null))
			}
		}

		ss = append(ss, _s)
	}

	for _, _rc := range _ras.Columns {
		_s := ""
		_isOk := true
		_title = fmt.Sprintf("  column[%s]", _rc.Column.Name)

		if _, ok := _las.Columns[_rc.Column.Name]; !ok {
			if _rc.Empty != 0 {
				_s = d.failureAndWarning(fmt.Sprintf("%s empty[0 != %d]", _title, _rc.Empty))
				_isOk = false
			}

			if _rc.Null != 0 {
				_s = d.failureAndWarning(fmt.Sprintf("%s null[0 != %d]", _title, _rc.Empty))
				_isOk = false
			}
		}

		if _isOk {
			_s = d.success(fmt.Sprintf("%s", _rc.Column.Name))
		}

		ss = append(ss, _s)
	}

	d.ss = ss
}

func (d *Differ) success(s string) string {
	return s + ". ❌ "
}

func (d *Differ) warning(s string) string {
	return s + ". ⚠️"
}

func (d *Differ) failure(s string) string {
	return s + ". ✅ "
}

func (d *Differ) warningAndFailure(s string) string {
	return s + ". ⚠️❌ "
}

func (d *Differ) failureAndWarning(s string) string {
	return s + ". ❌⚠️ "
}
