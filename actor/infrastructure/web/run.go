package web

import (
	"net/http"
	"sirius/actor/config"
	"sirius/actor/helper"
	"sirius/actor/infrastructure/web/route"

	"github.com/gin-gonic/gin"

	"reflect"
	"runtime"

	"sirius/actor/infrastructure/web/middleware"

	log "github.com/dynastymasra/gochill"
)

func Run() {
	pack := runtime.FuncForPC(reflect.ValueOf(Run).Pointer()).Name()

	log.Alert(log.Msg("Prepare run web application"), log.O("package", pack),
		log.O("version", config.Version))

	gin.SetMode(config.GinMode)
	router := gin.Default()

	router.Use(middleware.Header())

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, helper.FailResponse(config.ErrEndpointNotFound))
		return
	})

	v2 := router.Group("/v1")
	{
		route.ControllerRouter(v2)
	}

	router.Run(config.Address)
}
