package main

import (
	"os"
	"reflect"
	"runtime"

	"os/signal"
	"syscall"

	"fmt"

	"github.com/dynastymasra/shajaro/actor/config"
	"github.com/dynastymasra/shajaro/actor/infrastructure/web"

	"github.com/dynastymasra/shajaro/actor/infrastructure/provider"

	"github.com/dynastymasra/shajaro/actor/infrastructure/instrumentation"

	log "github.com/dynastymasra/gochill"
	"gopkg.in/tylerb/graceful.v1"
)

func init() {
	os.Setenv("SERVICE_NAME", config.AppName)

	config.InitConfig()
}

func main() {
	pack := runtime.FuncForPC(reflect.ValueOf(main).Pointer()).Name()

	log.Info(log.Msg("Prepare run application"), log.O("package", pack), log.O("version", config.Version),
		log.O("project", config.ProjectName))

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	log.Info(log.Msg("Prepare to run database migration"), log.O("package", pack), log.O("version", config.Version),
		log.O("project", config.ProjectName))

	if err := instrumentation.InitiateStatsD(config.StatsDHost, config.StatsDPort, config.AppName, config.StatsDEnable); err != nil {
		log.Alert(log.Msg("Failed setup statsd client", err.Error()), log.O("package", pack),
			log.O("version", config.Version), log.O("project", config.ProjectName))
		os.Exit(1)
	}
	defer instrumentation.CloseStatsDClient()

	db, err := provider.ConnectSQL()
	if err != nil {
		log.Alert(log.Msg("Failed connecting to database", err.Error()), log.O("package", pack),
			log.O("version", config.Version), log.O("project", config.ProjectName))
		os.Exit(1)
	}

	if err := provider.Migration(db); err != nil {
		log.Alert(log.Msg("Failed run database migration", err.Error()), log.O("package", pack),
			log.O("version", config.Version), log.O("project", config.ProjectName))
		os.Exit(1)
	}

	log.Info(log.Msg("Running migration success"), log.O("package", pack), log.O("version", config.Version),
		log.O("project", config.ProjectName))

	server := &graceful.Server{
		Timeout: 0,
	}
	go web.Run(server)

	log.Info(log.Msg("Application start running"), log.O("package", pack), log.O("version", config.Version),
		log.O("project", config.ProjectName))

	select {
	case sig := <-stop:
		provider.CloseDB(db)
		<-server.StopChan()
		log.Warn(log.Msg("Shutdown the service", fmt.Sprintf("%+v", sig)),
			log.O("package", pack), log.O("project", config.ProjectName), log.O("version", config.Version))
		os.Exit(0)
	}
}
