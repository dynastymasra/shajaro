package helper

import (
	"encoding/json"
	"reflect"
	"runtime"

	"github.com/dynastymasra/shajaro/actor/config"

	log "github.com/dynastymasra/gochill"
)

//Jsend used to format JSON with jsend rules
type Jsend struct {
	Status  string      `json:"status" binding:"required"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// FailResponse is used to return response with JSON format if failure
func FailResponse(msg string) Jsend {
	return Jsend{Status: "failed", Message: msg}
}

// SuccessResponse used to return response with JSON format success
func SuccessResponse() Jsend {
	return Jsend{Status: "success"}
}

// ObjectResponse used to return response JSON format if have data value
func ObjectResponse(data interface{}) Jsend {
	return Jsend{Status: "success", Data: data}
}

// Stringify used to stringify json object
func (j Jsend) Stringify() string {
	toJSON, err := json.Marshal(j)
	if err != nil {
		log.Error(log.Msg("Unable to stringify JSON", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName),
			log.O("package", runtime.FuncForPC(reflect.ValueOf(j.Stringify).Pointer()).Name()))
		return ""
	}
	return string(toJSON)
}
