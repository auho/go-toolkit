package source

import (
	"context"

	"github.com/auho/go-toolkit/redis/client"
)

type keyer interface {
	Size() (int64, error)
	Scan() map[string]string
}

type key struct {
	concurrency int
	pageSize    int64
	keyName     string
	client      *client.Redis
}

func (k *key) config(config Config) error {
	k.concurrency = config.Concurrency
	k.pageSize = config.PageSize
	k.keyName = config.Key

	var err error
	k.client, err = client.NewRedisClient(config.Options)
	if err != nil {
		return err
	}

	return nil
}

type sortedSetsKey struct {
	key
}

func newSortedSets(config Config) (*sortedSetsKey, error) {
	s := &sortedSetsKey{}
	err := s.config(config)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s sortedSetsKey) Size() (int64, error) {
	cmd := s.client.ZCard(context.Background(), s.keyName)
	if cmd.Err() != nil {
		return 0, cmd.Err()
	}

	return cmd.Val(), nil
}

func (s sortedSetsKey) Scan() map[string]string {

}
