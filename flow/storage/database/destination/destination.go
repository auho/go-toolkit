package destination

import (
	"errors"
	"fmt"
	"sync"

	"github.com/auho/go-simple-db/simple"
	"github.com/auho/go-toolkit/flow/storage"
	"github.com/auho/go-toolkit/flow/storage/database"
	"github.com/auho/go-toolkit/time/timing"
)

var _ storage.Destinationer[storage.MapEntry] = (*Destination[storage.MapEntry])(nil)
var _ database.Databaseor = (*Destination[storage.MapEntry])(nil)

type destinationer[E storage.Entry] interface {
	desFunc(driver simple.Driver, tableName string, items []E) error
}

func withConfig[E storage.Entry](config Config) func(*Destination[E]) error {
	return func(d *Destination[E]) error {
		err := d.config(config)
		if err != nil {
			return err
		}

		return nil
	}
}

func withDestinationer[E storage.Entry](der destinationer[E]) func(*Destination[E]) {
	return func(d *Destination[E]) {
		d.destinationer = der
	}
}

type Destination[E storage.Entry] struct {
	storage.Storage
	simple.Driver
	isDone        bool
	isTruncate    bool
	concurrency   int
	pageSize      int64
	tableName     string
	itemsChan     chan []E
	doWg          sync.WaitGroup
	state         *storage.State
	destinationer destinationer[E]
}

func newDestination[E storage.Entry](c func(*Destination[E]) error, os ...func(*Destination[E])) (*Destination[E], error) {
	d := &Destination[E]{}
	err := c(d)
	if err != nil {
		return nil, err
	}

	d.options(os)

	return d, nil
}

func (s *Destination[E]) GetDriver() simple.Driver {
	return s.Driver
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

	d.state = storage.NewState()
	d.state.Concurrency = d.concurrency
	d.state.Title = d.Title()
	d.state.StatusConfig()

	return
}

func (d *Destination[E]) Accept() (err error) {
	d.state.StatusAccept()
	d.state.DurationStart()

	if d.isTruncate {
		err = d.Truncate(d.tableName)
		if err != nil {
			return
		}
	}

	d.itemsChan = make(chan []E, d.concurrency)

	for i := 0; i < d.concurrency; i++ {
		d.doWg.Add(1)
		go func() {
			d.do()

			d.doWg.Done()
		}()
	}

	return nil
}

func (d *Destination[E]) Receive(items []E) {
	d.itemsChan <- items
}

func (d *Destination[E]) Done() {
	d.state.StatusDone()

	if d.isDone {
		return
	}

	d.isDone = true

	close(d.itemsChan)
}

func (d *Destination[E]) Finish() {
	d.doWg.Wait()

	d.Close()

	d.state.StatusFinish()
	d.state.DurationStop()
}

func (d *Destination[E]) do() {
	duration := timing.NewDuration()
	duration.Start()
	var descItems []E

	for items := range d.itemsChan {
		duration.Begin()

		itemsLen := int64(len(items))
		if itemsLen <= 0 {
			continue
		}

		for start := int64(0); start < itemsLen; start += d.pageSize {
			end := start + d.pageSize
			if end >= itemsLen {
				descItems = items[start:]
			} else {
				descItems = items[start:end]
			}

			err := d.destinationer.desFunc(d.Driver, d.tableName, descItems)
			if err != nil {
				panic(err)
			}
		}

		d.state.AddAmount(itemsLen)

		duration.End()
	}

	duration.Stop()
}

func (d *Destination[E]) Title() string {
	return fmt.Sprintf("Destinationer driver[%s]", d.DriverName())
}

func (d *Destination[E]) Summary() []string {
	return []string{fmt.Sprintf("%s Concurrency:%d", d.Title(), d.concurrency)}
}

func (d *Destination[E]) State() []string {
	return []string{d.state.Overview()}
}

func (d *Destination[E]) options(os []func(*Destination[E])) {
	for _, o := range os {
		o(d)
	}
}
