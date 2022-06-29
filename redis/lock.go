package redis

import (
	"context"
	"errors"
	"github.com/bsm/redislock"
	"time"
)

var ErrLockNotExists = errors.New("lock not exists")

// LockError 获取锁返回的错误
// Err() 表示获取锁未成功（包括 redis 操作失败、锁被占用的情况下获取不到锁等）
//
// 如需要判断锁被占用的情况下，获取不到锁，使用 ErrNotObtained()
//
//	if lockErr.ErrNotObtained() {
//		锁被占用，获取不到锁
//	} else if lockErr.Err() != nil {
// 		redis 操作失败
//	}
//
type LockError struct {
	err error
}

// ErrNotObtained 是否未成功获取到锁（锁被占用）
// true： 未获取到锁
//
func (le *LockError) ErrNotObtained() bool {
	return le.err == redislock.ErrNotObtained
}

// Err 返回错误，包括 redis 操作错误 和 未获取到锁错误
func (le *LockError) Err() error {
	return le.err
}

func (le *LockError) Error() string {
	return le.err.Error()
}

type Lock struct {
	lock *redislock.Lock
	ctx  context.Context
}

func newLock(l *redislock.Lock, ctx context.Context) *Lock {
	return &Lock{lock: l, ctx: ctx}
}

// TTL 返回剩余的生存时间，如果锁过期，返回 0
func (l *Lock) TTL() (time.Duration, error) {
	if l.lock == nil {
		return 0, ErrLockNotExists
	}

	return l.lock.TTL(l.ctx)
}

// Refresh 使用新的 TTL 刷新延长锁的时间
// 如果刷新没有成功，有可能返回 ErrNotObtained
func (l *Lock) Refresh(ctx context.Context, ttl time.Duration) *LockError {
	if l.lock == nil {
		return &LockError{err: ErrLockNotExists}
	}

	return &LockError{
		err: l.lock.Refresh(ctx, ttl, nil),
	}
}

// Release 手动释放锁，不要忘记释放锁
func (l *Lock) Release() {
	if l.lock == nil {
		return
	}

	// 有可能返回 ErrLockNotHeld
	_ = l.lock.Release(l.ctx)
}

// ObtainLockOnce
// 直接获取锁，只获取一次，redis 处理失败或者获取不到锁（锁被占用）直接返回
//
func (r *Redis) ObtainLockOnce(ctx context.Context, key string, ttl time.Duration) (*Lock, *LockError) {
	l, err := r.locker.Obtain(ctx, key, ttl, nil)
	if err != nil {
		return nil, &LockError{err: err}
	}

	return newLock(l, ctx), &LockError{}
}

// ObtainLockDeadline
// 每 backoff 时间重试一次获取锁，直到 wait 时间，如果还未获取到锁则返回
//
// example: ttl is 1s，Retry every 500ms, for up-to a minute
// ObtainLockDeadline(context.Background(), "key", time.Second, time.Minute, 500 * time.Millisecond)
//
func (r *Redis) ObtainLockDeadline(ctx context.Context, key string, ttl, wait, backoff time.Duration) (lock *Lock, cancel func(), lockError *LockError) {
	lockCtx, cancel := context.WithDeadline(ctx, time.Now().Add(wait))

	l, err := r.locker.Obtain(lockCtx, key, ttl, &redislock.Options{
		RetryStrategy: redislock.LinearBackoff(backoff),
	})

	if err != nil {
		cancel()
		return nil, nil, &LockError{err: err}
	}

	return newLock(l, ctx), cancel, &LockError{}
}

// ObtainLockRetry
// 每 backoff 时间重试一次获取锁，最多重试 retryNum 次数，如果还未获取到锁则返回
//
// example: Retry every 100ms, for up-to 3x
// ObtainLockRetry(context.Background(), "key", time.Second, 3, 100 * time.Millisecond)
//
func (r *Redis) ObtainLockRetry(ctx context.Context, key string, ttl time.Duration, retryNum int, backoff time.Duration) (*Lock, *LockError) {
	l, err := r.locker.Obtain(ctx, key, ttl, &redislock.Options{
		RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(backoff), retryNum),
	})

	if err != nil {
		return nil, &LockError{err: err}
	}

	return newLock(l, ctx), &LockError{}
}
