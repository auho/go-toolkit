package redis

import (
	"context"
	"fmt"
	redis2 "github.com/go-redis/redis/v8"
	"os"
	"sync"
	"testing"
	"time"
)

var redisClient *Redis

func setup() {
	var err error
	redisClient, err = NewRedisClient(&redis2.Options{
		Network:  "tcp",
		Addr:     "127.0.0.1:6379",
		Password: "123456",
		DB:       0,
	})

	if err != nil {
		panic(err)
	}
}

func teardown() {
	_ = redisClient.Close()
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func Test_Once(t *testing.T) {
	ttl := time.Millisecond * 500

	key := fmt.Sprintf("test:redis:locker:once:%d", time.Now().Nanosecond())
	lock, lockErr := redisClient.ObtainLockOnce(context.Background(), key, ttl)
	if lockErr.Err() != nil {
		t.Fatal(lockErr.Error())
	}

	defer lock.Release()

	var sw sync.WaitGroup
	for i := 0; i < 10; i++ {
		sw.Add(1)

		go func() {
			lock1, lockErr := redisClient.ObtainLockOnce(context.Background(), key, ttl)
			if !lockErr.ErrNotObtained() {
				t.Error(lockErr.Error())
				return
			} else {
				if !lockErr.ErrNotObtained() {
					lock1.Release()
				}
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

	lock, cancel, lockErr := redisClient.ObtainLockDeadline(context.Background(), key, wholeTll, wait, backoff)
	if lockErr.Err() != nil {
		t.Fatal(lockErr.Error())
	}

	defer cancel()
	defer lock.Release()

	c := time.NewTicker(backoff)
	max := int(wholeTll/wait) - 1
	i := 0
	for range c.C {
		if i >= max {
			c.Stop()
			break
		}

		lock1, cancel1, lockErr := redisClient.ObtainLockDeadline(context.Background(), key, ttl, wait, backoff)
		if lockErr.ErrNotObtained() {

		} else if lockErr.Err() != nil {
			t.Error(lockErr.Error())

		} else {
			lock1.Release()
			cancel1()
			t.Error("can not obtain lock")

		}

		i++
	}

	lock1, cancel1, lockErr := redisClient.ObtainLockDeadline(context.Background(), key, ttl, wait, backoff)
	if lockErr.ErrNotObtained() {
		t.Fatal(lockErr.Error())
	} else if lockErr.Err() != nil {
		t.Fatal(lockErr.Error())
	} else {
		cancel1()
		lock1.Release()
	}
}

func Test_Retry(t *testing.T) {
	ttl := time.Millisecond * 500
	retryNum := 3
	backoff := time.Millisecond * 50
	key := fmt.Sprintf("test:redis:locker:retry:%d", time.Now().Nanosecond())

	lock, lockErr := redisClient.ObtainLockRetry(context.Background(), key, ttl, retryNum, backoff)
	if lockErr.Err() != nil {
		t.Fatal(lockErr.Error())
	}

	defer lock.Release()

	c := time.NewTicker(backoff)
	max := int(ttl/backoff)/retryNum - 1
	i := 0
	for range c.C {
		if i >= max {
			c.Stop()
			break
		}

		lock1, lockErr := redisClient.ObtainLockRetry(context.Background(), key, ttl, retryNum, backoff)
		if lockErr.ErrNotObtained() {

		} else if lockErr.Err() != nil {
			t.Fatal(lockErr.Error())
		} else {
			lock1.Release()
			t.Fatal("failure")
		}

		i++
	}

	lock1, lockErr := redisClient.ObtainLockRetry(context.Background(), key, ttl, retryNum, backoff)
	if lockErr.ErrNotObtained() {
		t.Fatal(lockErr.Error())
	} else if lockErr.Err() != nil {
		t.Fatal(lockErr.Error())
	} else {
		lock1.Release()
	}
}
