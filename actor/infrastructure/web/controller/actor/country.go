package actor

import (
	"io/ioutil"
	"reflect"
	"runtime"

	"encoding/json"
	"net/http"

	"github.com/dynastymasra/shajaro/actor/config"
	"github.com/dynastymasra/shajaro/actor/domain/actor"
	"github.com/dynastymasra/shajaro/actor/helper"

	"fmt"

	log "github.com/dynastymasra/gochill"
)

func CountryListController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	traceKey := r.Context().Value(config.TraceKey)
	pack := runtime.FuncForPC(reflect.ValueOf(CountryListController).Pointer()).Name()

	log.Info(log.Msg("Request list countries"), log.O("version", config.Version),
		log.O("project", config.ProjectName), log.O("package", pack),
		log.O(config.TraceKey, traceKey))

	raw, err := ioutil.ReadFile(config.CountryJSON)
	if err != nil {
		log.Error(log.Msg("Failed read file", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O("package", pack),
			log.O(config.TraceKey, traceKey), log.O("file", config.CountryJSON))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, helper.FailResponse(err.Error()).Stringify())
		return
	}

	var countries []actor.Country
	if err := json.Unmarshal(raw, &countries); err != nil {
		log.Error(log.Msg("Failed unmarshal byte", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O("package", pack),
			log.O(config.TraceKey, traceKey), log.O("data", string(raw)))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, helper.FailResponse(err.Error()).Stringify())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, helper.ObjectResponse(countries).Stringify())
}
