package main

import (
	"os"
	"sirius/actor/config"
	"sirius/actor/infrastructure/web"
)

func init() {
	os.Setenv("SERVICE_NAME", config.AppName)

	config.InitConfig()
}

func main() {
	web.Run()
}
