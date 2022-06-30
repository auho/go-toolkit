package showdifference

import (
	"log"
	"os"
)

type File struct {
	Differ
}

// ShowDiff 输出差异
// filePath  文件路径
func (f *File) ShowDiff(filePath1, filePath2 string) (string, error) {
	f1, err := os.ReadFile(filePath1)
	if err != nil {
		log.Fatal(err)
	}

	f2, err := os.ReadFile(filePath2)
	if err != nil {
		log.Fatal(err)
	}

	return f.Compare(f1, f2)
}
