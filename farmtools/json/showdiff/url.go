package showdiff

import (
	"io"
	"net/http"
)

type Url struct {
	Differ
}

// ShowDiff 输出差异
// url  http url
func (u *Url) ShowDiff(url1, url2 string) (string, error) {
	b1, b2, err := u.readBody(url1, url2)
	if err != nil {
		return "", err
	}

	return u.Compare(b1, b2)
}

func (u *Url) ShowDiffAndHead(url1, url2 string, line int) (string, error) {
	b1, b2, err := u.readBody(url1, url2)
	if err != nil {
		return "", err
	}

	return u.CompareAndHead(b1, b2, line)
}

func (u *Url) readBody(url1, url2 string) ([]byte, []byte, error) {
	r1, err := http.Get(url1)
	if err != nil {
		return nil, nil, err
	}

	defer func() { _ = r1.Body.Close() }()

	b1, err := io.ReadAll(r1.Body)
	if err != nil {
		return nil, nil, err
	}

	r2, err := http.Get(url2)
	if err != nil {
		return nil, nil, err
	}

	defer func() { _ = r2.Body.Close() }()

	b2, err := io.ReadAll(r2.Body)
	if err != nil {
		return nil, nil, err
	}

	return b1, b2, nil
}
