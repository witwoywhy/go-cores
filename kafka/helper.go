package kafka

import (
	"math/rand/v2"
	"time"
)

func jitteredBackoff(retries int) time.Duration {
	minBackoff := 250 * time.Millisecond
	maxBackoff := 5 * time.Second

	backoff := minBackoff * (1 << uint(retries))
	if backoff > maxBackoff || backoff <= 0 {
		backoff = maxBackoff
	}

	jitter := time.Duration(float64(backoff) * 0.2 * (2*rand.Float64() - 1))
	return backoff + jitter
}
