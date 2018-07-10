package middleware

import (
	"net/http"
	"runtime"
	"strings"

	"github.com/dynastymasra/shajaro/actor/infrastructure/instrumentation"

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
