package publisher

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"dns/mdns/zone"
)

type Publisher struct {
	mutex    sync.Mutex
	records  zone.Records
	change   chan zone.Records
	shutdown chan struct{}

	wg sync.WaitGroup
}

func newPublisher() (*Publisher, error) {
	return &Publisher{}, nil
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
	z, zoneCancel, err := zone.NewZone(zone.Config{
		EnableIpv4: c.enableIpv4,
		EnableIpv6: c.enableIpv6,
	})
	if err != nil {
		log.Fatal(err)
	}

	p.wg.Add(2)
	go p.recordsChange(c)
	go p.multicast(c, z)

	p.wg.Wait()

	c.close()
	zoneCancel()
}

func (p *Publisher) multicast(c *Config, z *zone.Zone) {
	defer p.wg.Done()

	var jitter = time.Millisecond * 100 * time.Duration(rand.Intn(10))
	for {
		select {
		case <-time.NewTimer(c.broadcastInterval + jitter).C:
			p.mutex.Lock()
			if len(p.records) <= 0 {
				continue
			}

			err := z.BroadcastRecords(p.records)
			if err != nil {
				log.Println(err)
			}

			fmt.Println(p.records)

			p.mutex.Unlock()

		case <-p.shutdown:
			return
		}
	}
}

func (p *Publisher) recordsChange(c *Config) {
	defer p.wg.Done()

	for {
		select {
		case records := <-c.recordsChangeChan:
			p.mutex.Lock()
			p.records = records
			p.mutex.Unlock()
		case <-p.shutdown:
			return
		}
	}
}
