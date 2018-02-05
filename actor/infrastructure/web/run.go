package web

import (
	"net/http"
	"sirius/actor/config"
	"sirius/actor/utility"

	"sirius/actor/infrastructure/web/route"

	"github.com/gin-gonic/gin"
)

func Run() {
	gin.SetMode(config.GinMode)
	router := gin.Default()

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, utility.FailResponse("Endpoint your requested not found"))
	})

	v2 := router.Group("/v1")
	{
		route.ControllerRouter(v2)
	}

	router.Run(config.Address)
}
