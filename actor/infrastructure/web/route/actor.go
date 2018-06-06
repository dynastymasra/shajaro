package route

import (
	"shajaro/actor/infrastructure/web/controller/actor"

	"github.com/gin-gonic/gin"
)

func ActorRouter(router *gin.RouterGroup) {
	router.POST("/register", actor.RegisterController)

	router.GET("/countries", actor.CountryListController)
}
