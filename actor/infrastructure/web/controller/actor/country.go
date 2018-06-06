package actor

import (
	"io/ioutil"
	"reflect"
	"runtime"

	"encoding/json"
	"net/http"

	"shajaro/actor/config"
	"shajaro/actor/domain/actor"
	"shajaro/actor/helper"

	log "github.com/dynastymasra/gochill"
	"github.com/gin-gonic/gin"
)

func CountryListController(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	pack := runtime.FuncForPC(reflect.ValueOf(CountryListController).Pointer()).Name()

	log.Info(log.Msg("Request list countries"), log.O("version", config.Version),
		log.O("project", config.ProjectName), log.O("package", pack),
		log.O(config.TraceKey, c.GetString(config.TraceKey)))

	raw, err := ioutil.ReadFile(config.CountryJSON)
	if err != nil {
		log.Error(log.Msg("Failed read file", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O("package", pack),
			log.O(config.TraceKey, c.GetString(config.TraceKey)), log.O("file", config.CountryJSON))
		c.Error(err)
		c.JSON(http.StatusInternalServerError, helper.FailResponse(err.Error()))
		return
	}

	var countries []actor.Country
	if err := json.Unmarshal(raw, &countries); err != nil {
		log.Error(log.Msg("Failed unmarshal byte", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O("package", pack),
			log.O(config.TraceKey, c.GetString(config.TraceKey)), log.O("data", string(raw)))
		c.Error(err)
		c.JSON(http.StatusInternalServerError, helper.FailResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, helper.ObjectResponse(countries))
}
