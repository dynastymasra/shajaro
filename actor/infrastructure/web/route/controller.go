package route

import (
	"shajaro/actor/infrastructure/web/controller"

	"github.com/gin-gonic/gin"
)

// ControllerRouter group router
func ControllerRouter(router *gin.RouterGroup) {
	router.Any("/healthz", controller.HealthzController)
}
