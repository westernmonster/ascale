package dlock

import (
	"ascale/pkg/cache/redis"
	"ascale/pkg/log"
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"strconv"
	"sync"
	"time"
)

var (
	luaRefresh = redis.NewScript(1, `if redis.call("get", KEYS[1]) == ARGV[1] then return redis.call("pexpire", KEYS[1], ARGV[2]) else return 0 end`)
	luaRelease = redis.NewScript(1, `if redis.call("get", KEYS[1]) == ARGV[1] then return redis.call("del", KEYS[1]) else return 0 end`)
	luaPTTL    = redis.NewScript(1, `if redis.call("get", KEYS[1]) == ARGV[1] then return redis.call("pttl", KEYS[1]) else return -3 end`)
)

var (
	// ErrNotObtained is returned when a lock cannot be obtained.
	ErrNotObtained = errors.New("redislock: not obtained")

	// ErrLockNotHeld is returned when trying to release an inactive lock.
	ErrLockNotHeld = errors.New("redislock: lock not held")
)

type Client struct {
	client *redis.Pool
	tmp    []byte
	tmpMu  sync.Mutex
}

func New(pool *redis.Pool) *Client {
	return &Client{client: pool}
}

func (c *Client) Obtain(ctx context.Context, key string, ttl time.Duration, opt *Options) (*Lock, error) {
	// Create a random token
	token, err := c.randomToken()
	if err != nil {
		return nil, err
	}

	value := token + opt.getMetadata()
	retry := opt.getRetryStrategy()

	var cancel context.CancelFunc = func() {}
	if _, ok := ctx.Deadline(); !ok {
		ctx, cancel = context.WithDeadline(ctx, time.Now().Add(ttl))
	}
	defer cancel()

	var timer *time.Timer
	for {
		ok, err := c.obtain(ctx, key, value, ttl)
		if err != nil {
			return nil, err
		} else if ok {
			return &Lock{client: c, key: key, value: value}, nil
		}

		backoff := retry.NextBackoff()
		if backoff < 1 {
			return nil, ErrNotObtained
		}

		if timer == nil {
			timer = time.NewTimer(backoff)
			defer timer.Stop()
		} else {
			timer.Reset(backoff)
		}

		select {
		case <-ctx.Done():
			return nil, ErrNotObtained
		case <-timer.C:
		}
	}
}

func (c *Client) obtain(ctx context.Context, key, value string, ttl time.Duration) (ret bool, err error) {
	var conn redis.Conn
	if conn, err = c.client.GetContext(ctx); err != nil {
		log.For(ctx).Errorf("dlock.obtain(),  err(%+v)", err)
		return
	}
	defer conn.Close()

	var result string
	if result, err = redis.String(conn.Do("SET", key, value, "NX", "EX", int64(ttl.Seconds()))); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.For(ctx).Errorf("dlock.obtain() SETNX(%s) error(%v)", key, err)
			return
		}
	}

	if result == "OK" {
		ret = true
	}

	return
}

func (c *Client) randomToken() (string, error) {
	c.tmpMu.Lock()
	defer c.tmpMu.Unlock()

	if len(c.tmp) == 0 {
		c.tmp = make([]byte, 16)
	}

	if _, err := io.ReadFull(rand.Reader, c.tmp); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(c.tmp), nil
}

// --------------------------------------------------------------------

// Lock represents an obtained, distributed lock.
type Lock struct {
	client *Client
	key    string
	value  string
}

// Obtain is a short-cut for New(...).Obtain(...).
func Obtain(c context.Context, client *redis.Pool, key string, ttl time.Duration, opt *Options) (*Lock, error) {
	return New(client).Obtain(c, key, ttl, opt)
}

// Key returns the redis key used by the lock.
func (l *Lock) Key() string {
	return l.key
}

// Token returns the token value set by the lock.
func (l *Lock) Token() string {
	return l.value[:22]
}

// Metadata returns the metadata of the lock.
func (l *Lock) Metadata() string {
	return l.value[22:]
}

// TTL returns the remaining time-to-live. Returns 0 if the lock has expired.
func (l *Lock) TTL(c context.Context) (dt time.Duration, err error) {
	var conn redis.Conn
	if conn, err = l.client.client.GetContext(c); err != nil {
		log.For(c).Errorf("dlock.TTL(),  err(%+v)", err)
		return
	}
	defer conn.Close()

	res, err := redis.Int64(luaPTTL.Do(conn, l.key, l.value))
	if err == redis.ErrNil {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	if res > 0 {
		return time.Duration(res) * time.Millisecond, nil
	}
	return 0, nil
}

// Refresh extends the lock with a new TTL.
// May return ErrNotObtained if refresh is unsuccessful.
func (l *Lock) Refresh(c context.Context, ttl time.Duration, opt *Options) (err error) {
	var conn redis.Conn
	if conn, err = l.client.client.GetContext(c); err != nil {
		log.For(c).Errorf("dlock.Refresh(),  err(%+v)", err)
		return
	}
	defer conn.Close()

	ttlVal := strconv.FormatInt(int64(ttl/time.Millisecond), 10)

	status, err := redis.Int64(luaRefresh.Do(conn, l.key, l.value, ttlVal))
	if err != nil {
		return err
	} else if status == int64(1) {
		return nil
	}
	return ErrNotObtained
}

// Release manually releases the lock.
// May return ErrLockNotHeld.
func (l *Lock) Release(c context.Context) (err error) {
	var conn redis.Conn
	if conn, err = l.client.client.GetContext(c); err != nil {
		log.For(c).Errorf("dlock.Release(),  err(%+v)", err)
		return
	}
	defer conn.Close()

	res, err := redis.Int64(luaRelease.Do(conn, l.key, l.value))
	if err == redis.ErrNil {
		return ErrLockNotHeld
	} else if err != nil {
		return err
	}

	if res != 1 {
		return ErrLockNotHeld
	}

	return nil
}

// --------------------------------------------------------------------

// Options describe the options for the lock
type Options struct {
	// RetryStrategy allows to customise the lock retry strategy.
	// Default: do not retry
	RetryStrategy RetryStrategy

	// Metadata string is appended to the lock token.
	Metadata string

	// Context provides an optional context for timeout and cancellation control.
	// If requested, Obtain will by default retry until the TTL exires. This
	// behaviour can be tweaked with a custom context deadline.
	Context context.Context
}

func (o *Options) getMetadata() string {
	if o != nil {
		return o.Metadata
	}
	return ""
}

func (o *Options) getContext() context.Context {
	if o != nil && o.Context != nil {
		return o.Context
	}
	return context.Background()
}

func (o *Options) getRetryStrategy() RetryStrategy {
	if o != nil && o.RetryStrategy != nil {
		return o.RetryStrategy
	}
	return NoRetry()
}

// --------------------------------------------------------------------

// RetryStrategy allows to customise the lock retry strategy.
type RetryStrategy interface {
	// NextBackoff returns the next backoff duration.
	NextBackoff() time.Duration
}

type linearBackoff time.Duration

// LinearBackoff allows retries regularly with customized intervals
func LinearBackoff(backoff time.Duration) RetryStrategy {
	return linearBackoff(backoff)
}

// NoRetry acquire the lock only once.
func NoRetry() RetryStrategy {
	return linearBackoff(0)
}

func (r linearBackoff) NextBackoff() time.Duration {
	return time.Duration(r)
}

type limitedRetry struct {
	s RetryStrategy

	cnt, max int
}

// LimitRetry limits the number of retries to max attempts.
func LimitRetry(s RetryStrategy, max int) RetryStrategy {
	return &limitedRetry{s: s, max: max}
}

func (r *limitedRetry) NextBackoff() time.Duration {
	if r.cnt >= r.max {
		return 0
	}
	r.cnt++
	return r.s.NextBackoff()
}

type exponentialBackoff struct {
	cnt uint

	min, max time.Duration
}

// ExponentialBackoff strategy is an optimization strategy with a retry time of 2**n milliseconds (n means number of times).
// You can set a minimum and maximum value, the recommended minimum value is not less than 16ms.
func ExponentialBackoff(min, max time.Duration) RetryStrategy {
	return &exponentialBackoff{min: min, max: max}
}

func (r *exponentialBackoff) NextBackoff() time.Duration {
	r.cnt++

	ms := 2 << 25
	if r.cnt < 25 {
		ms = 2 << r.cnt
	}

	if d := time.Duration(ms) * time.Millisecond; d < r.min {
		return r.min
	} else if r.max != 0 && d > r.max {
		return r.max
	} else {
		return d
	}
}
