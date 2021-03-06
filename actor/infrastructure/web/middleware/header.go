package middleware

import (
	"github.com/dynastymasra/shajaro/actor/config"

	"net/http"

	"github.com/dynastymasra/shajaro/actor/helper"

	"strings"

	"context"

	"fmt"

	"github.com/labstack/gommon/random"
	"github.com/urfave/negroni"
)

func TraceKey() negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		traceKey := r.Header.Get(config.HeaderTraceID)
		if len(traceKey) < 1 {
			traceKey = random.String(11)
		}
		ctx := context.WithValue(r.Context(), config.TraceKey, traceKey)
		next(w, r.WithContext(ctx))
	}
}

func RequestType() negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		contentType := r.Header.Get("Content-Type")
		httpMethod := r.Method

		if httpMethod == http.MethodPost || httpMethod == http.MethodPut {
			if !strings.Contains(contentType, "application/json") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnsupportedMediaType)
				fmt.Fprintf(w, helper.FailResponse("content type is not supported").Stringify())
				return
			}
		}
		next(w, r)
	}
}
