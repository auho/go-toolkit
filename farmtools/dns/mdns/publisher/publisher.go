package publisher

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"dns/mdns/zone"
)

type Publisher struct {
	mutex             sync.Mutex
	records           zone.Records
	recordsChan       chan zone.Records
	recordsUpdateChan chan struct{}
	shutdown          chan struct{}

	wg sync.WaitGroup
}

func newPublisher() (*Publisher, error) {
	return &Publisher{
		recordsChan:       make(chan zone.Records),
		recordsUpdateChan: make(chan struct{}),
		shutdown:          make(chan struct{}),
	}, nil
}

func runPublisher(c *Config) error {
	p, err := newPublisher()
	if err != nil {
		return err
	}

	p.start(c)

	return nil
}

func (p *Publisher) start(c *Config) {

	p.wg.Add(3)
	go p.loop(c)
	go p.handleRecordsChange(c)
	go p.multicast(c)

	p.wg.Wait()

	c.close()
}

func (p *Publisher) loop(c *Config) {
	defer p.wg.Done()

	var jitter = time.Millisecond * 100 * time.Duration(rand.Intn(10))
	for {
		select {
		case <-time.NewTimer(c.broadcastInterval + jitter).C:
			p.recordsUpdateChan <- struct{}{}
		case <-p.shutdown:
			return
		}
	}
}

func (p *Publisher) multicast(c *Config) {
	defer p.wg.Done()

	for {
		select {
		case <-p.recordsUpdateChan:
			p.mutex.Lock()

			if len(p.records) > 0 {
				z, zoneCancel, err := zone.NewZone(zone.Config{
					EnableIpv4: c.enableIpv4,
					EnableIpv6: c.enableIpv6,
				})
				if err != nil {
					log.Fatal(err)
				}

				err = z.BroadcastRecords(p.records)
				if err != nil {
					log.Println(err)
				}

				zoneCancel()

				log.Println(p.records)
			}

			p.mutex.Unlock()

		case <-p.shutdown:
			return
		}
	}
}

func (p *Publisher) handleRecordsChange(c *Config) {
	defer p.wg.Done()

	for {
		select {
		case records := <-c.recordsChangeChan:
			p.mutex.Lock()
			p.records = records
			p.mutex.Unlock()

			p.recordsUpdateChan <- struct{}{}

		case <-p.shutdown:
			return
		}
	}
}
