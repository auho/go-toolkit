package destination

import (
	"errors"
	"github.com/auho/go-simple-db/simple"
)

func iss(fields []string) desFunc[[][]interface{}] {
	return func(sd simple.Driver, tableName string, items [][]interface{}) error {

		_, err := sd.BulkInsertFromSliceSlice(tableName, fields, items)

		return err
	}
}

func NewInsertSliceSlice(config Config, fields []string) (*Destination[[][]interface{}], error) {
	if len(fields) <= 0 {
		return nil, errors.New("fields is error")
	}

	d := &Destination[[][]interface{}]{}
	err := d.config(config)
	if err != nil {
		return nil, err
	}

	d.desFunc = iss(fields)

	return d, nil
}

func ism() desFunc[[]map[string]interface{}] {
	return func(sd simple.Driver, tableName string, items []map[string]interface{}) error {
		_, err := sd.BulkInsertFromSliceMap(tableName, items)

		return err
	}
}

func NewInsertSliceMap(config Config) (*Destination[[]map[string]interface{}], error) {
	d := &Destination[[]map[string]interface{}]{}
	err := d.config(config)
	if err != nil {
		return nil, err
	}

	d.desFunc = ism()

	return d, nil
}
