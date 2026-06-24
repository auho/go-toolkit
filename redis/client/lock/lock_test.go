package lock

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

var locker *RedisLocker

func setup() {
	var err error

	client := redis.NewClient(&redis.Options{
		Network:  "tcp",
		Addr:     "127.0.0.1:6379",
		Password: "123456",
		DB:       0,
	})

	locker, err = NewRedisLocker(client)

	if err != nil {
		panic(err)
	}
}

func teardown() {}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func Test_Once(t *testing.T) {
	ttl := time.Millisecond * 500

	key := fmt.Sprintf("test:redis:locker:once:%d", time.Now().Nanosecond())
	lock, ok, err := locker.ObtainLockOnce(context.Background(), key, ttl)
	if err != nil {
		t.Fatal("err", err)
	}

	if !ok {
		t.Fatal("ok", ok)
	}

	defer lock.Release()

	var sw sync.WaitGroup
	for i := 0; i < 10; i++ {
		sw.Add(1)

		go func() {
			lock1, ok1, err1 := locker.ObtainLockOnce(context.Background(), key, ttl)
			if err1 != nil {
				t.Error("err1", err1)
				return
			}

			if ok1 {
				t.Error("can not obtain lock")
				lock1.Release()
			}

			sw.Done()
		}()

		sw.Wait()
	}
}

func Test_Deadline(t *testing.T) {
	wholeTll := time.Millisecond * 500
	ttl := time.Millisecond * 100
	wait := time.Millisecond * 100
	backoff := time.Millisecond * 50
	key := fmt.Sprintf("test:redis:locker:deadline:%d", time.Now().Nanosecond())

	lock, cancel, ok, err := locker.ObtainLockDeadline(context.Background(), key, wholeTll, wait, backoff)
	if err != nil {
		t.Fatal("err", err)
	}
	if !ok {
		t.Fatal("ok", ok)
	}

	defer cancel()
	defer lock.Release()

	c := time.NewTicker(backoff)
	_max := int(wholeTll/wait) - 1
	i := 0
	for range c.C {
		if i >= _max {
			c.Stop()
			break
		}

		lock1, cancel1, ok1, err1 := locker.ObtainLockDeadline(context.Background(), key, ttl, wait, backoff)
		if err1 != nil {
			t.Error("err1", err1)
		}

		if ok1 {
			t.Error("can not obtain lock")
			lock1.Release()
			cancel1()
		}

		i++
	}

	lock2, cancel2, ok2, err2 := locker.ObtainLockDeadline(context.Background(), key, ttl, wait, backoff)
	if err2 != nil {
		t.Fatal("err2", err2)
	}

	if !ok2 {
		t.Fatal("ok2", ok2)
	}

	lock2.Release()
	cancel2()
}

func Test_Retry(t *testing.T) {
	ttl := time.Millisecond * 500
	retryNum := 3
	backoff := time.Millisecond * 50
	key := fmt.Sprintf("test:redis:locker:retry:%d", time.Now().Nanosecond())

	lock, ok, err := locker.ObtainLockRetry(context.Background(), key, ttl, retryNum, backoff)
	if err != nil {
		t.Fatal("err", err)
	}
	if !ok {
		t.Fatal("ok", ok)
	}

	defer lock.Release()

	c := time.NewTicker(backoff)
	_max := int(ttl/backoff)/retryNum - 1
	i := 0
	for range c.C {
		if i >= _max {
			c.Stop()
			break
		}

		lock1, ok1, err1 := locker.ObtainLockRetry(context.Background(), key, ttl, retryNum, backoff)
		if err1 != nil {
			t.Fatal("err1", err1)
		}

		if ok1 {
			lock1.Release()
			t.Error("can not obtain lock")
		}

		i++
	}

	lock2, ok2, err2 := locker.ObtainLockRetry(context.Background(), key, ttl, retryNum, backoff)
	if err2 != nil {
		t.Fatal("err2", err2)
	}
	if !ok2 {
		t.Fatal("ok2", ok2)
	}

	lock2.Release()
}
