package controller

import (
	"net/http"
	"sirius/actor/helper"

	"sirius/actor/config"

	"reflect"
	"runtime"

	log "github.com/dynastymasra/gochill"
	"github.com/gin-gonic/gin"
)

// HealthzController to check service is ok
func HealthzController(c *gin.Context) {
	log.Info(log.Msg("Request healthz"), log.O("version", config.Version), log.O("project", config.ProjectName),
		log.O("package", runtime.FuncForPC(reflect.ValueOf(HealthzController).Pointer()).Name()),
		log.O(config.TraceKey, c.GetString(config.TraceKey)))

	c.Header("Content-Type", "application/json")

	c.JSON(http.StatusOK, helper.SuccessResponse())
}
