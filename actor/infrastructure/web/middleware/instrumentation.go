package middleware

import (
	"net/http"
	"runtime"
	"strings"

	"github.com/dynastymasra/shajaro/actor/infrastructure/instrumentation"

	"context"

	"github.com/dynastymasra/shajaro/actor/config"
	"github.com/newrelic/go-agent"
	"github.com/urfave/negroni"
)

func StatsDMiddlewareLogger() negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		t := instrumentation.NewTimingStatsD()
		totalGoroutine := runtime.NumGoroutine()

		next(w, r)

		key := strings.Replace(r.URL.Path, "/", ".", len(r.URL.Path))

		instrumentation.TimingSend(key+".time", t)
		instrumentation.StatsDIncrement(key + ".calls")
		instrumentation.StatsDGauge(key+".goroutines", totalGoroutine)
	}
}

func NewrelicMiddlewareHandler() negroni.HandlerFunc {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		_, nextHandler := newrelic.WrapHandleFunc(instrumentation.NewRelicApp(), r.URL.Path, next)
		nextHandler(w, r)
	})
}

func NewrelicInstrumentation() negroni.HandlerFunc {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		newRelicTrx, ok := w.(newrelic.Transaction)
		if !ok {
			next(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), config.NewrelicTxnKey, newRelicTrx)
		next(w, r.WithContext(ctx))
	})
}
