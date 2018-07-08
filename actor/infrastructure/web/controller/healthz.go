package controller

import (
	"net/http"

	"reflect"
	"runtime"

	"shajaro/actor/config"
	"shajaro/actor/helper"
	"shajaro/actor/infrastructure/provider"

	"fmt"

	log "github.com/dynastymasra/gochill"
)

// HealthzController to check service is ok
func HealthzController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	traceKey := r.Context().Value(config.TraceKey)
	pack := runtime.FuncForPC(reflect.ValueOf(HealthzController).Pointer()).Name()

	log.Info(log.Msg("Request healthz"), log.O("version", config.Version), log.O("package", pack),
		log.O("project", config.ProjectName), log.O(config.TraceKey, traceKey))

	db, err := provider.ConnectSQL()
	if err != nil {
		log.Error(log.Msg("Failed connect to database", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, traceKey), log.O("package", pack))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, helper.FailResponse(err.Error()).Stringify())
		return
	}

	if err := provider.SQLPing(db); err != nil {
		log.Error(log.Msg("Failed ping database", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, traceKey), log.O("package", pack))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, helper.FailResponse(err.Error()).Stringify())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, helper.SuccessResponse().Stringify())
}
