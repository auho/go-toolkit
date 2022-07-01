package destination

import (
	"errors"
	"fmt"
	"github.com/auho/go-simple-db/simple"
	"github.com/auho/go-toolkit/flow/storage"
	"github.com/auho/go-toolkit/flow/storage/database"
	"github.com/auho/go-toolkit/time/timing"
	"sync"
	"sync/atomic"
)

type SliceEntries [][]interface{}
type MapEntries []map[string]interface{}

type entries interface {
	SliceEntries | MapEntries
}

type desFunc[E entries] func(driver simple.Driver, tableName string, items E) error

type Destination[E entries] struct {
	storage.Storage
	simple.Driver
	isDone      bool
	isTruncate  bool
	concurrency int
	pageSize    int64
	tableName   string
	itemsChan   chan E
	doWg        sync.WaitGroup
	state       *database.State
	desFunc     desFunc[E]
}

func (d *Destination[E]) config(config Config) (err error) {
	d.isTruncate = config.IsTruncate
	d.concurrency = config.Concurrency
	d.pageSize = config.PageSize
	d.tableName = config.TableName

	d.Driver, err = simple.NewDriver(config.Driver, config.Dsn)
	if err != nil {
		return err
	}

	if d.concurrency <= 0 {
		err = errors.New(fmt.Sprintf("concurrency[%d] is error", d.concurrency))
		return
	}

	if d.pageSize <= 0 {
		err = errors.New(fmt.Sprintf("page size[%d] is error", d.pageSize))
		return
	}

	d.state = database.NewState()
	d.state.Concurrency = d.concurrency
	d.state.PageSize = d.pageSize
	d.state.Title = d.Title()
	d.state.Status = "config"

	return
}

func (d *Destination[E]) Accept() (err error) {
	d.state.Status = "accept"
	d.state.Duration.Start()

	if d.isTruncate {
		err = d.Truncate(d.tableName)
		if err != nil {
			return
		}
	}

	d.itemsChan = make(chan E, d.concurrency)

	for i := 0; i < d.concurrency; i++ {
		d.doWg.Add(1)
		go func() {
			d.do()

			d.doWg.Done()
		}()
	}

	return nil
}

func (d *Destination[E]) Receive(items E) {
	d.itemsChan <- items
}

func (d *Destination[E]) Done() {
	d.state.Status = "done"

	if d.isDone {
		return
	}

	d.isDone = true

	close(d.itemsChan)
}

func (d *Destination[E]) Finish() {
	d.doWg.Wait()

	d.Close()

	d.state.Status = "finish"
	d.state.Duration.Stop()
}

func (d *Destination[E]) do() {
	duration := timing.NewDuration()
	duration.Start()
	var targetItems E

	for {
		if items, ok := <-d.itemsChan; ok {
			duration.Begin()

			itemsLen := int64(len(items))
			if itemsLen <= 0 {
				continue
			}

			for start := int64(0); start < itemsLen; start += d.pageSize {
				end := start + d.pageSize
				if end >= itemsLen {
					targetItems = items[start:]
				} else {
					targetItems = items[start:end]
				}

				err := d.desFunc(d.Driver, d.tableName, targetItems)
				if err != nil {
					panic(err)
				}
			}

			atomic.AddInt64(&d.state.Amount, itemsLen)
			duration.End()
		} else {
			break
		}
	}

	duration.Stop()
}

func (d *Destination[E]) Title() string {
	return fmt.Sprintf("Destination driver[%s]", d.DriverName())
}
