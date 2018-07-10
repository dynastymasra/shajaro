package api

import (
	"time"

	"github.com/dynastymasra/shajaro/actor/config"

	"github.com/cenkalti/backoff"
)

func BackOffRetry() *backoff.ExponentialBackOff {
	back := backoff.NewExponentialBackOff()
	back.MaxElapsedTime = time.Duration(config.RetryDuration) * time.Second
	back.MaxInterval = time.Duration(config.MaxRetryInterval) * time.Second

	return back
}
