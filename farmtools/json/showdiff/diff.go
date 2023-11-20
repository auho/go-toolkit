package showdiff

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/yudai/gojsondiff"
	"github.com/yudai/gojsondiff/formatter"
)

type Differ struct {
}

// CompareAndHead 对比不同 json 的差异并截断
// line 截断前多少行
//
// [...] array
// {...} object
func (d *Differ) CompareAndHead(leftJson, rightJson []byte, line int) (string, error) {
	s, err := d.Compare(leftJson, rightJson)
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
//
// [...] array
// {...} object
func (d *Differ) Compare(leftJson, rightJson []byte) (string, error) {
	var err error
	var leftAny any
	var difference gojsondiff.Diff

	differ := gojsondiff.New()

	if string(leftJson[0]) == "[" {
		var leftSlice, rightSlice []any
		err = json.Unmarshal(leftJson, &leftSlice)
		if err != nil {
			return "", fmt.Errorf("failed to left slice unmarshal; %w", err)
		}

		err = json.Unmarshal(rightJson, &rightSlice)
		if err != nil {
			return "", fmt.Errorf("failed to right slice unmarshal; %w", err)
		}

		difference = differ.CompareArrays(leftSlice, rightSlice)
		leftAny = leftSlice
	} else if string(leftJson[0]) == "{" {
		var leftMap, rightMap map[string]any
		err = json.Unmarshal(leftJson, &leftMap)
		if err != nil {
			return "", fmt.Errorf("failed to left map unmarshal; %w", err)
		}

		err = json.Unmarshal(rightJson, &rightMap)
		if err != nil {
			return "", fmt.Errorf("failed to right map unmarshal; %w", err)
		}

		difference = differ.CompareObjects(leftMap, rightMap)
		leftAny = leftMap
	} else {
		return "", fmt.Errorf("failed json unmarshal")
	}

	if !difference.Modified() {
		return "", nil
	}

	config := formatter.AsciiFormatterConfig{
		ShowArrayIndex: true,
		Coloring:       true,
	}

	f := formatter.NewAsciiFormatter(leftAny, config)
	diffString, err := f.Format(difference)
	if err != nil {
		return "", fmt.Errorf("failed to format; %w", err)
	}

	return diffString, nil
}
