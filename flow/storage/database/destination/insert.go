package destination

import (
	"errors"
	"github.com/auho/go-simple-db/simple"
)

func iss(fields []string) desFunc[SliceEntries] {
	return func(sd simple.Driver, tableName string, items SliceEntries) error {

		_, err := sd.BulkInsertFromSliceSlice(tableName, fields, items)

		return err
	}
}

func NewInsertSliceSlice(config Config, fields []string) (*Destination[SliceEntries], error) {
	if len(fields) <= 0 {
		return nil, errors.New("fields is error")
	}

	d := &Destination[SliceEntries]{}
	err := d.config(config)
	if err != nil {
		return nil, err
	}

	d.desFunc = iss(fields)

	return d, nil
}

func ism() desFunc[MapEntries] {
	return func(sd simple.Driver, tableName string, items MapEntries) error {
		_, err := sd.BulkInsertFromSliceMap(tableName, items)

		return err
	}
}

func NewInsertSliceMap(config Config) (*Destination[MapEntries], error) {
	d := &Destination[MapEntries]{}
	err := d.config(config)
	if err != nil {
		return nil, err
	}

	d.desFunc = ism()

	return d, nil
}
