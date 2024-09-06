package restapi

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

func NewClient(config elasticsearch.Config) (*elasticsearch.Client, error) {
	return elasticsearch.NewClient(config)
}

func GetResponseBody(resp *esapi.Response, err error) (string, error) {
	_b, err := getResponseBytes(resp, err)
	if err != nil {
		return "", err
	}

	return string(_b), nil
}

func getResponseBytes(resp *esapi.Response, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	return io.ReadAll(resp.Body)
}

func DoResponse[T any](fn func() (*esapi.Response, error), a T) (T, error) {
	resp, err := fn()
	body, err := getResponseBytes(resp, err)
	if err != nil {
		return a, err
	}

	err = json.Unmarshal(body, &a)
	if err != nil {
		return a, fmt.Errorf("json unmarshal; %w", err)
	}

	return a, nil
}
