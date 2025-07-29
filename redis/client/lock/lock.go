package lock

import (
	"context"
	"errors"
	"fmt"
	"github.com/bsm/redislock"
	"time"
)

var ErrLockNotExists = errors.New("lock not exists")

type Lock struct {
	lock *redislock.Lock
}

func newLock(l *redislock.Lock) *Lock {
	return &Lock{lock: l}
}

// TTL 返回剩余的生存时间，如果锁过期，返回 0
func (l *Lock) TTL() (time.Duration, error) {
	if l.lock == nil {
		return 0, ErrLockNotExists
	}

	return l.lock.TTL(context.Background())
}

// Refresh 使用新的 TTL 刷新延长锁的时间
func (l *Lock) Refresh(ctx context.Context, ttl time.Duration) (ok bool, err error) {
	if l.lock == nil {
		return false, ErrLockNotExists
	}

	err = l.lock.Refresh(ctx, ttl, nil)
	if errors.Is(err, redislock.ErrNotObtained) {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("refresh: %w", err)
	}

	return true, nil
}

// Release 手动释放锁，不要忘记释放锁
func (l *Lock) Release() {
	if l.lock == nil {
		return
	}

	// 有可能返回 ErrLockNotHeld
	_ = l.lock.Release(context.Background())
}

// RedisLocker redis locker
type RedisLocker struct {
	locker *redislock.Client
}

// ObtainLockOnce
// 直接获取锁，只获取一次，redis 处理失败或者获取不到锁（锁被占用）直接返回
func (r *RedisLocker) ObtainLockOnce(ctx context.Context, key string, ttl time.Duration) (lock *Lock, ok bool, err error) {
	key = r.genKey(key)

	l, err := r.locker.Obtain(ctx, key, ttl, nil)
	if errors.Is(err, redislock.ErrNotObtained) {
		return
	} else if err != nil {
		err = fmt.Errorf("obtain: %w", err)
		return
	}

	return newLock(l), true, nil
}

// ObtainLockDeadline
// 每 backoff 时间重试一次获取锁，直到 wait 时间，如果还未获取到锁则返回
//
// example: ttl is 1s，Retry every 500ms, for up-to a minute
// ObtainLockDeadline(context.Background(), "key", time.Second, time.Minute, 500 * time.Millisecond)
func (r *RedisLocker) ObtainLockDeadline(ctx context.Context, key string, ttl, wait, backoff time.Duration) (lock *Lock, cancel func(), ok bool, err error) {
	key = r.genKey(key)

	var lockCtx context.Context
	lockCtx, cancel = context.WithDeadline(ctx, time.Now().Add(wait))

	var l *redislock.Lock
	l, err = r.locker.Obtain(lockCtx, key, ttl, &redislock.Options{
		RetryStrategy: redislock.LinearBackoff(backoff),
	})

	if errors.Is(err, redislock.ErrNotObtained) {
		return
	} else if err != nil {
		cancel()
		err = fmt.Errorf("obtain: %w", err)
		return
	}

	return newLock(l), cancel, true, nil
}

// ObtainLockRetry
// 每 backoff 时间重试一次获取锁，最多重试 retryNum 次数，如果还未获取到锁则返回
//
// example: Retry every 100ms, for up-to 3x
// ObtainLockRetry(context.Background(), "key", time.Second, 3, 100 * time.Millisecond)
func (r *RedisLocker) ObtainLockRetry(ctx context.Context, key string, ttl time.Duration, retryNum int, backoff time.Duration) (lock *Lock, ok bool, err error) {
	key = r.genKey(key)

	var l *redislock.Lock
	l, err = r.locker.Obtain(ctx, key, ttl, &redislock.Options{
		RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(backoff), retryNum),
	})

	if errors.Is(err, redislock.ErrNotObtained) {
		return
	} else if err != nil {
		err = fmt.Errorf("obtain: %w", err)
		return
	}

	return newLock(l), true, nil
}

func (r *RedisLocker) genKey(key string) string {
	return key + ":lock_"
}
