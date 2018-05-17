package main

import (
	"os"
	"reflect"
	"runtime"
	"sirius/actor/config"
	"sirius/actor/infrastructure/web"

	"os/signal"
	"syscall"

	"fmt"

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

	go web.Run()

	log.Info(log.Msg("Application start running"), log.O("package", pack), log.O("version", config.Version),
		log.O("project", config.ProjectName))

	select {
	case sig := <-stop:
		log.Warn(log.Msg("Application prepare to shutdown", fmt.Sprintf("%+v", sig)),
			log.O("package", pack), log.O("project", config.ProjectName), log.O("version", config.Version))
		os.Exit(0)
	}
}
