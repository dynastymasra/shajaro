package main

import (
	"os"
	"reflect"
	"runtime"

	"os/signal"
	"syscall"

	"fmt"

	"shajaro/actor/config"
	"shajaro/actor/infrastructure/provider"
	"shajaro/actor/infrastructure/web"

	log "github.com/dynastymasra/gochill"
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

	go web.Run()

	log.Info(log.Msg("Application start running"), log.O("package", pack), log.O("version", config.Version),
		log.O("project", config.ProjectName))

	select {
	case sig := <-stop:
		log.Warn(log.Msg("Application prepare to shutdown", fmt.Sprintf("%+v", sig)),
			log.O("package", pack), log.O("project", config.ProjectName), log.O("version", config.Version))
		provider.CloseDB(db)
		os.Exit(0)
	}
}
