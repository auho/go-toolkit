package showdiff

import (
	"io"
	"net/http"
)

type Request struct {
	Differ
}

func (u *Request) ShowDiff(req1, req2 *http.Request) (string, error) {
	b1, b2, err := u.readBody(req1, req2)
	if err != nil {
		return "", err
	}

	return u.Compare(b1, b2)
}

func (u *Request) ShowDiffAndHead(req1, req2 *http.Request, line int) (string, error) {
	b1, b2, err := u.readBody(req1, req2)
	if err != nil {
		return "", err
	}

	return u.CompareAndHead(b1, b2, line)
}

func (u *Request) readBody(req1, req2 *http.Request) ([]byte, []byte, error) {
	_client := &http.Client{}

	r1, err := _client.Do(req1)
	if err != nil {
		return nil, nil, err
	}

	defer func() { _ = r1.Body.Close() }()

	b1, err := io.ReadAll(r1.Body)
	if err != nil {
		return nil, nil, err
	}

	r2, err := _client.Do(req2)
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
