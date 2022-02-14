/*
 * Usage:
 *
 * initialInterval := 1 * time.Second
 * maxInterval := 3 * time.Second
 * // maxRetries refers to only the number of retry attempts, and don't
 * // include the first attempt.
 * maxRetries := 3
 *
 * eb := NewExponentialBackoff(initialInterval, maxInterval, maxRetries)
 * for {
 *    res, err := serviceCall(arg, timeout=1)
 *    if err == nil || err != TIMEOUT_ERROR {
 *       break
 *    }
 *
 *    duration := eb.Backoff()
 *		// Reach maximum number of retries.
 *    if duration == Stop {
 *      break
 *    }
 * }
 *
 */

package misc

import (
	"time"
	rd "math/rand"
)

const Stop time.Duration = -1
const UnlimitedRetry int = -1
const UnlimitedTimeout time.Duration = -1
const NextDelayUnavailable time.Duration = -2

type Backoff interface {
	Backoff() time.Duration
	Reset()
}

type ConstantBackoff struct {
	Interval       time.Duration
	MaxRetries     int
	currentRetries int
}

func (b *ConstantBackoff) Backoff() time.Duration {
	b.currentRetries++
	// Over limit of retry amount.
	if b.MaxRetries != UnlimitedRetry && b.currentRetries > b.MaxRetries {
		return Stop
	}
	time.Sleep(b.Interval)
	return b.Interval
}

func (b *ConstantBackoff) Reset() {
	b.currentRetries = 0
}

// One of MaxRetries or RetryTimeout needs to be provided.
type ExponentialBackoff struct {
	InitialInterval time.Duration
	MaxInterval     time.Duration
	MaxRetries      int
	currentRetries  int
	hitMaxInterval  bool
	startTime       time.Time
	RetryTimeout    time.Duration
	randomBackoff   bool
}

func (b *ExponentialBackoff) Backoff() time.Duration {
	b.currentRetries++
	// Over limit of retry amount.
	if b.MaxRetries != UnlimitedRetry && b.currentRetries > b.MaxRetries {
		return Stop
	}
	// Over limit of retry timeout.
	curTime := time.Now()
	if b.RetryTimeout != UnlimitedTimeout &&
		curTime.Sub(b.startTime) >= b.RetryTimeout {
		return Stop
	}
	var delay time.Duration
	if b.hitMaxInterval {
		delay = b.MaxInterval
	} else {
		delay = b.InitialInterval << uint64((b.currentRetries - 1))
		if delay > b.MaxInterval {
			delay = b.MaxInterval
			b.hitMaxInterval = true
		} else if b.randomBackoff && b.currentRetries != 1 {
			// we take a random value from 2^(retries-1) to 2^(retries)
			prevDelay := delay >> 1
			delay = time.Duration(int64(prevDelay) + rd.Int63n(int64(delay) - int64(prevDelay)))
		}
	}
	if b.RetryTimeout != UnlimitedTimeout &&
		curTime.Sub(b.startTime)+delay > b.RetryTimeout {
		delay = b.RetryTimeout - curTime.Sub(b.startTime)
	}
	time.Sleep(delay)
	return delay
}

func (b *ExponentialBackoff) NextRetryDelay() time.Duration {
	if b.randomBackoff {
		// next retry delay cannot be calculated if randomized
		// backoff is used
		return NextDelayUnavailable
	}
	if b.hitMaxInterval {
		return b.MaxInterval
	}
	return b.InitialInterval << uint64(b.currentRetries)
}

func (b *ExponentialBackoff) SetRandomBackoff() {
	rd.Seed(time.Now().UnixNano())
	b.randomBackoff = true
}

func (b *ExponentialBackoff) Reset() {
	b.hitMaxInterval = false
	b.currentRetries = 0
	b.startTime = time.Now()
}

func NewConstantBackoff(interval time.Duration, retries int) *ConstantBackoff {
	return &ConstantBackoff{
		Interval:   interval,
		MaxRetries: retries}
}

/*
initialInterval - Init Interval between retry.
maxInterval - Max Interval between retry.
maxRetries - Max number of retries. It refers to only the number of retry
attempts, and don't include the first attempt.
*/
func NewExponentialBackoff(initialInterval, maxInterval time.Duration,
	maxRetries int) *ExponentialBackoff {
	return &ExponentialBackoff{
		InitialInterval: initialInterval,
		MaxInterval:     maxInterval,
		MaxRetries:      maxRetries,
		startTime:       time.Now(),
		RetryTimeout:    UnlimitedTimeout}
}

/*
initialInterval - Init Interval between retry.
maxInterval - Max Interval between retry.
retryTimeout - Max duration for retries.
*/
func NewExponentialBackoffWithTimeout(initialInterval,
	maxInterval time.Duration, retryTimeout time.Duration) *ExponentialBackoff {
	return &ExponentialBackoff{
		InitialInterval: initialInterval,
		MaxInterval:     maxInterval,
		MaxRetries:      UnlimitedRetry,
		startTime:       time.Now(),
		RetryTimeout:    retryTimeout}
}
