package api

import (
	"sirius/actor/config"
	"time"

	"github.com/cenkalti/backoff"
)

func BackOffRetry() *backoff.ExponentialBackOff {
	back := backoff.NewExponentialBackOff()
	back.MaxElapsedTime = time.Duration(config.RetryDuration) * time.Second
	back.MaxInterval = 10 * time.Second

	return back
}
