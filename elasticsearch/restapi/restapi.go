package restapi

import (
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"io"
)

func Do[T any](fn func() (*esapi.Response, error), a T) (T, error) {
	resp, err := fn()
	if err != nil {
		return a, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return a, err
	}

	err = json.Unmarshal(body, &a)
	if err != nil {
		return a, fmt.Errorf("json unmarshal; %w", err)
	}

	return a, nil
}
