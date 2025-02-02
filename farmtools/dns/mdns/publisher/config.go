package publisher

import (
	"time"

	"dns/mdns/zone"
)

type Config struct {
	enableIpv4        bool
	enableIpv6        bool
	broadcastInterval time.Duration

	recordsChangeChan chan zone.Records
}

func (c *Config) config() error {
	err := c.checkArgs()
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) close() {
	close(c.recordsChangeChan)
}

func (c *Config) checkArgs() error {
	c.recordsChangeChan = make(chan zone.Records)

	if !c.enableIpv4 && !c.enableIpv6 {
		//return fmt.Errorf("neigher ipv4 nor ipv6 set")
		c.enableIpv4 = true
	}

	if c.broadcastInterval <= time.Second*30 {
		c.broadcastInterval = time.Second * 30
	}

	if c.broadcastInterval > time.Second*900 {
		c.broadcastInterval = time.Second * 900
	}

	return nil
}
