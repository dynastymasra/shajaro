package route

import (
	"shajaro/actor/infrastructure/web/controller/actor"

	"shajaro/actor/config"
	"shajaro/actor/infrastructure/web/middleware"

	"github.com/gin-gonic/gin"
)

func ActorRouter(router *gin.RouterGroup) {
	router.POST("/register", actor.RegisterController)
	router.POST("/login", actor.LoginController)

	router.GET("/countries", actor.CountryListController)
	router.GET("/actor", middleware.ValidateScope(config.ActorRead), actor.GetUserByIDController)

	router.PUT("/actor", middleware.ValidateScope(config.ActorUpdate), actor.UpdateUserController)

	router.DELETE("/actor", middleware.ValidateScope(config.ActorDelete), actor.DeleteUserController)
}
