package client

import (
	"context"
	redis2 "github.com/go-redis/redis/v8"
	"log"
	"time"
)

func ExampleRedis_ObtainLockOnce() {
	var err error
	redisClient, err = NewRedisClient(&redis2.Options{
		Network:  "tcp",
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	if err != nil {
		panic(err)
	}

	lock, lockErr := redisClient.ObtainLockOnce(context.Background(), "redis:key", time.Millisecond*500)
	if lockErr.ErrNotObtained() {
		// 没有获取到锁
	} else if lockErr.Err() != nil {
		// 没有获取到锁，或者 redis 操作错误
		log.Fatal(lockErr.Error())
	}

	// 释放锁
	lock.Release()

	// 返回剩余的生存时间
	_, _ = lock.TTL()
}

func ExampleRedis_ObtainLockDeadline() {
	lock, cancel, lockErr := redisClient.ObtainLockDeadline(
		context.Background(),
		"redis:key",
		time.Second*5,
		time.Second,
		time.Millisecond*500,
	)
	if lockErr.ErrNotObtained() {
		// 没有获取到锁
	} else if lockErr.Err() != nil {
		// 没有获取到锁，或者 redis 操作错误
		log.Fatal(lockErr.Error())
	}

	// 取消获取锁
	defer cancel()

	// 释放锁
	lock.Release()
}

func ExampleRedis_ObtainLockRetry() {
	lock, lockErr := redisClient.ObtainLockRetry(
		context.Background(),
		"redis:key",
		time.Second*10,
		5,
		time.Second,
	)
	if lockErr.ErrNotObtained() {
		// 没有获取到锁
	} else if lockErr.Err() != nil {
		// 没有获取到锁，或者 redis 操作错误
		log.Fatal(lockErr.Error())
	}

	// 释放锁
	lock.Release()
}
