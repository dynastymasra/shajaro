package web

import (
	"net/http"
	"shajaro/actor/config"
	"shajaro/actor/helper"
	"shajaro/actor/infrastructure/web/middleware"
	"shajaro/actor/infrastructure/web/route"

	"github.com/gin-gonic/gin"

	"reflect"
	"runtime"

	log "github.com/dynastymasra/gochill"
)

func Run() {
	pack := runtime.FuncForPC(reflect.ValueOf(Run).Pointer()).Name()

	log.Info(log.Msg("Start run web application"), log.O("package", pack),
		log.O("version", config.Version), log.O("project", config.ProjectName))

	gin.SetMode(config.GinMode)
	router := gin.Default()

	router.Use(middleware.RequestKey())

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, helper.FailResponse(config.ErrEndpointNotFound))
		return
	})

	v1 := router.Group("/v1")
	{
		route.ControllerRouter(v1)
		route.ActorRouter(v1)
	}

	router.Run(config.Address)
}
