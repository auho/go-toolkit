package destination

import (
	"fmt"
	"github.com/auho/go-toolkit/flow/storage"
	"sync"
	"sync/atomic"
)

type Destination[E storage.Entry] struct {
	isDone    bool
	amount    int64
	itemsChan chan []E
	chanWg    sync.WaitGroup
}

func (d *Destination[E]) Accept() error {
	d.itemsChan = make(chan []E)

	d.chanWg.Add(1)
	go func() {
		for items := range d.itemsChan {
			var descItems []E
			descItems = items[0:]
			atomic.AddInt64(&d.amount, int64(len(descItems)))
		}

		d.chanWg.Done()
	}()

	return nil
}

func (d *Destination[E]) Done() {
	if d.isDone {
		return
	}

	d.isDone = true
	close(d.itemsChan)
}

func (d *Destination[E]) Finish() {
	d.chanWg.Wait()
}

func (d *Destination[E]) Summary() []string {
	return []string{fmt.Sprintf("%s", d.title())}
}

func (d *Destination[E]) State() []string {
	return []string{fmt.Sprintf("Amount: %d", d.amount)}
}

func (d *Destination[E]) title() string {
	return "Mock:desc"
}
