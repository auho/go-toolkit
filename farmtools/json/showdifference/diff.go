package showdifference

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/yudai/gojsondiff"
	"github.com/yudai/gojsondiff/formatter"
	"regexp"
	"strings"
)

type Differ struct {
}

// CompareAndHead 对比不同 json 的差异并截断
// line 截断前多少行
func (d *Differ) CompareAndHead(aJson, bJson []byte, s string, line int) (string, error) {
	s, err := d.Compare(aJson, bJson)
	if err != nil {
		return s, err
	}

	re := regexp.MustCompile("(?m)^\x1b.+")
	newLines := re.FindAllString(s, -1)
	if len(newLines) > line {
		return strings.Join(newLines[0:line], "\n"), nil
	} else {
		return strings.Join(newLines, "\n"), nil
	}
}

// Compare 对比不同 json 的差异
func (d *Differ) Compare(aJson, bJson []byte) (string, error) {
	differ := gojsondiff.New()
	difference, err := differ.Compare(aJson, bJson)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Failed to compare: %s", err.Error()))
	}

	if !difference.Modified() {
		return "", nil
	}

	var aJsonMap map[string]interface{}
	err = json.Unmarshal(aJson, &aJsonMap)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Failed to unmarshal: %s", err.Error()))
	}

	config := formatter.AsciiFormatterConfig{
		ShowArrayIndex: true,
		Coloring:       true,
	}

	f := formatter.NewAsciiFormatter(aJson, config)
	diffString, err := f.Format(difference)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Failed to format: %s", err.Error()))
	}

	return diffString, nil
}
