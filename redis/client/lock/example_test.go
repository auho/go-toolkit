package lock

import (
	"context"
	redis2 "github.com/go-redis/redis/v8"
	"log"
	"time"
)

func ExampleRedisLocker_ObtainLockOnce() {
	var err error
	client := redis2.NewClient(&redis2.Options{
		Network:  "tcp",
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	locker, err = NewRedisLocker(client)

	if err != nil {
		panic(err)
	}

	lock, ok, err := locker.ObtainLockOnce(context.Background(), "redis:key", time.Millisecond*500)
	if err != nil {
		// redis 操作错误
		log.Fatal(err)
	}

	if !ok {
		// 没有获取到锁
		log.Fatal("lock not obtained")
	}

	// 释放锁
	lock.Release()

	// 返回剩余的生存时间
	_, _ = lock.TTL()
}

func ExampleRedisLocker_ObtainLockDeadline() {
	lock, cancel, ok, err := locker.ObtainLockDeadline(
		context.Background(),
		"redis:key",
		time.Second*5,
		time.Second,
		time.Millisecond*500,
	)
	if err != nil {
		// redis 操作错误
		log.Fatal(err)
	}

	if !ok {
		// 没有获取到锁
		log.Fatal("lock not obtained")
	}

	// 取消获取锁
	defer cancel()

	// 释放锁
	lock.Release()
}

func ExampleRedisLocker_ObtainLockRetry() {
	lock, ok, err := locker.ObtainLockRetry(
		context.Background(),
		"redis:key",
		time.Second*10,
		5,
		time.Second,
	)
	if err != nil {
		// redis 操作错误
		log.Fatal(err)
	}

	if !ok {
		// 没有获取到锁
		log.Fatal("lock not obtained")
	}

	// 释放锁
	lock.Release()
}
