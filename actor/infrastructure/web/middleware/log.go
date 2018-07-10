package middleware

import (
	"net/http"
	"reflect"
	"runtime"
	"time"

	"github.com/dynastymasra/shajaro/actor/config"

	log "github.com/dynastymasra/gochill"
	"github.com/urfave/negroni"
)

func HTTPStatLogger() negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		startTime := time.Now()

		next(w, r)

		responseTime := time.Now()
		deltaTime := responseTime.Sub(startTime).Seconds() * 1000

		log.Info(log.Msg("Actor HTTP request log"), log.O("version", config.Version),
			log.O("package", runtime.FuncForPC(reflect.ValueOf(HTTPStatLogger).Pointer()).Name()),
			log.O("project", config.ProjectName), log.O("request_time", startTime.Format(time.RFC3339)),
			log.O("response_time", responseTime.Format(time.RFC3339)), log.O("delta_time", deltaTime),
			log.O("url", r.URL.Path), log.O("method", r.Method), log.O("request_proxy", r.RemoteAddr),
			log.O("request_source", r.Header.Get("X-FORWARDED-FOR")),
			log.O("trace_id", r.Context().Value(config.TraceKey)))
	}
}
