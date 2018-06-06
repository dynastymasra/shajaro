package controller

import (
	"net/http"

	"reflect"
	"runtime"

	"shajaro/actor/config"
	"shajaro/actor/helper"
	"shajaro/actor/infrastructure/provider"

	log "github.com/dynastymasra/gochill"
	"github.com/gin-gonic/gin"
)

// HealthzController to check service is ok
func HealthzController(c *gin.Context) {
	log.Info(log.Msg("Request healthz"), log.O("version", config.Version), log.O("project", config.ProjectName),
		log.O("package", runtime.FuncForPC(reflect.ValueOf(HealthzController).Pointer()).Name()),
		log.O(config.TraceKey, c.GetString(config.TraceKey)))

	c.Header("Content-Type", "application/json")

	db, err := provider.ConnectSQL()
	if err != nil {
		log.Info(log.Msg("Failed connect to database", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, c.GetString(config.TraceKey)),
			log.O("package", runtime.FuncForPC(reflect.ValueOf(HealthzController).Pointer()).Name()))
		c.Error(err)
		c.JSON(http.StatusInternalServerError, helper.FailResponse(config.ErrDatabaseConnectFail))
		return
	}

	if err := provider.SQLPing(db); err != nil {
		log.Info(log.Msg("Failed ping database", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, c.GetString(config.TraceKey)),
			log.O("package", runtime.FuncForPC(reflect.ValueOf(HealthzController).Pointer()).Name()))
		c.Error(err)
		c.JSON(http.StatusInternalServerError, helper.FailResponse(config.ErrPingDatabaseFail))
		return
	}

	c.JSON(http.StatusOK, helper.SuccessResponse())
}
