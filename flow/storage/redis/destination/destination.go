package destination

import (
	"github.com/auho/go-toolkit/redis/client"
)

type Destination struct {
	concurrency int
	pageSize    int64
	keyName     string
	client      *client.Redis
}

func (d *Destination) config(config *Config) error {
	var err error
	d.client, err = client.NewRedisClient(config.Options)
	if err != nil {
		return err
	}

	return nil
}
