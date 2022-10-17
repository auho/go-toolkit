package showdiff

import (
	"os"
)

type File struct {
	Differ
}

// ShowDiff 输出差异
// filePath  文件路径
func (f *File) ShowDiff(filePath1, filePath2 string) (string, error) {
	j1, j2, err := f.readContent(filePath1, filePath2)
	if err != nil {
		return "", err
	}

	return f.Compare(j1, j2)
}

func (f *File) ShowDiffAndHead(filePath1, filePath2 string, line int) (string, error) {
	j1, j2, err := f.readContent(filePath1, filePath2)
	if err != nil {
		return "", err
	}

	return f.CompareAndHead(j1, j2, line)
}

// 读取文件内容
// filePath  文件路径
func (f *File) readContent(filePath1, filePath2 string) ([]byte, []byte, error) {
	var f1, f2 []byte
	var err error

	f1, err = os.ReadFile(filePath1)
	if err != nil {
		return nil, nil, err
	}

	f2, err = os.ReadFile(filePath2)
	if err != nil {
		return nil, nil, err
	}

	return f1, f2, nil
}
