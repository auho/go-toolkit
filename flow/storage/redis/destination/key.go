package destination

import (
	"context"
	"fmt"
	"sync"

	"github.com/auho/go-toolkit/flow/storage"
	"github.com/auho/go-toolkit/flow/storage/redis"
	"github.com/auho/go-toolkit/redis/client"
)

var _ storage.Destinationer[storage.MapEntry] = (*key[storage.MapEntry])(nil)

type keyer[E storage.Entry] interface {
	redis.Keyer
	accept(itemsChan <-chan []E, c *client.Redis, key string, pageSize int64)
	stateAmount() int64
}

type key[E storage.Entry] struct {
	storage.Storage
	concurrency int
	isTruncate  bool
	pageSize    int64
	keyName     string
	isDone      bool
	itemsChan   chan []E
	doWg        sync.WaitGroup
	client      *client.Redis
	keyer       keyer[E]
	state       *storage.State
}

func newKey[E storage.Entry](config Config, keyer keyer[E]) (*key[E], error) {
	k := &key[E]{}
	k.keyer = keyer
	err := k.config(config)
	if err != nil {
		return nil, err
	}

	return k, nil
}

func (k *key[E]) GetClient() *client.Redis {
	return k.client
}

func (k *key[E]) Accept() error {
	k.state.StatusAccept()
	k.state.DurationStart()

	if k.isTruncate {
		_, err := k.keyer.Truncate(context.Background(), k.client, k.keyName)
		if err != nil {
			return err
		}
	}

	k.itemsChan = make(chan []E, k.concurrency)

	for i := 0; i < k.concurrency; i++ {
		k.doWg.Add(1)
		go func() {
			k.keyer.accept(k.itemsChan, k.client, k.keyName, k.pageSize)

			k.doWg.Done()
		}()
	}

	return nil
}

func (k *key[E]) Receive(items []E) {
	k.itemsChan <- items
}

func (k *key[E]) Done() {
	k.state.StatusDone()

	if k.isDone {
		return
	}

	k.isDone = true

	close(k.itemsChan)
}

func (k *key[E]) Finish() {
	k.doWg.Wait()

	k.state.StatusFinish()
	k.state.DurationStop()
}

func (k *key[E]) Summary() []string {
	return []string{fmt.Sprintf("%s Concurrency:%d; page size:%d", k.Title(), k.concurrency, k.pageSize)}
}

func (k *key[E]) State() []string {
	k.state.SetAmount(k.keyer.stateAmount())
	return []string{k.state.Overview()}
}

func (k *key[E]) Title() string {
	return fmt.Sprintf("Destiantion redis[%s] %s", k.keyer.Type(), k.keyName)
}

func (k *key[E]) Close() error {
	return k.client.Close()
}

func (k *key[E]) config(config Config) error {
	k.isTruncate = config.IsTruncate
	k.concurrency = config.Concurrency
	k.pageSize = config.PageSize
	k.keyName = config.Key

	if k.concurrency <= 0 {
		k.concurrency = 1
	}

	if k.pageSize <= 0 {
		k.pageSize = 20
	}

	if k.keyName == "" {
		k.LogFatalWithTitle("key name is empty")
	}

	var err error
	k.client, err = client.NewRedisClient(config.Options)
	if err != nil {
		return err
	}

	k.state = storage.NewState()
	k.state.Title = k.Title()
	k.state.StatusConfig()

	return nil
}
