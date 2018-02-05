package controller

import (
	"net/http"
	"sirius/actor/helper"

	"github.com/gin-gonic/gin"
)

// PingController to check service is ok
func PingController(c *gin.Context) {
	c.JSON(http.StatusOK, helper.SuccessResponse())
}
