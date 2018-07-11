package web

import (
	"github.com/dynastymasra/shajaro/actor/config"

	"reflect"
	"runtime"

	"github.com/urfave/negroni"

	"net/http"
	"os"

	"github.com/dynastymasra/shajaro/actor/infrastructure/web/middleware"

	nrgorilla "github.com/newrelic/go-agent/_integrations/nrgorilla/v1"

	log "github.com/dynastymasra/gochill"
	"github.com/dynastymasra/shajaro/actor/infrastructure/instrumentation"
	"gopkg.in/tylerb/graceful.v1"
)

func Run(server *graceful.Server) {
	pack := runtime.FuncForPC(reflect.ValueOf(Run).Pointer()).Name()

	log.Info(log.Msg("Start run web application"), log.O("package", pack),
		log.O("version", config.Version), log.O("project", config.ProjectName))

	newRelicApp := instrumentation.NewRelicApp()
	muxRouter := Router()
	router := nrgorilla.InstrumentRoutes(muxRouter, newRelicApp)

	n := negroni.New(negroni.NewRecovery())

	n.Use(middleware.RequestType())
	n.Use(middleware.TraceKey())
	n.Use(middleware.HTTPStatLogger())

	if config.StatsDEnable {
		log.Info(log.Msg("Service used instrumentation statsd"), log.O("package", pack),
			log.O("version", config.Version), log.O("project", config.ProjectName))

		n.Use(middleware.StatsDMiddlewareLogger())
	}

	if instrumentation.NewRelicConfig().Enabled {
		log.Info(log.Msg("Service used instrumentation newrelic"), log.O("package", pack),
			log.O("version", config.Version), log.O("project", config.ProjectName))

		n.Use(middleware.NewrelicMiddlewareHandler())
	}

	n.UseHandlerFunc(router.ServeHTTP)

	server.Server = &http.Server{
		Addr:    config.Address,
		Handler: n,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Alert(log.Msg("Failed to start server", err.Error()), log.O("package", pack),
			log.O("version", config.Version), log.O("project", config.ProjectName))
		os.Exit(1)
	}
}
