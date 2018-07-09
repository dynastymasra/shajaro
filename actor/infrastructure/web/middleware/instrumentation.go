package middleware

import (
	"net/http"
	"runtime"
	"shajaro/actor/config"
	"shajaro/actor/infrastructure/instrumentation"
	"strings"

	"github.com/urfave/negroni"
)

func StatsDMiddlewareLogger() negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		t := instrumentation.NewTimingStatsD(config.StatsDEnable)
		totalGoroutine := runtime.NumGoroutine()

		next(w, r)

		key := strings.Replace(r.URL.Path, "/", ".", len(r.URL.Path))

		instrumentation.TimingSend(key+".time", t, config.StatsDEnable)
		instrumentation.StatsDIncrement(key+".calls", config.StatsDEnable)
		instrumentation.StatsDGauge(key+".goroutines", totalGoroutine, config.StatsDEnable)
	}
}
