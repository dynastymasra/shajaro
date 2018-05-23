package domain

import (
	"database/sql/driver"
	"encoding/json"
	"reflect"
	"runtime"
	"sirius/actor/config"
	"sirius/actor/helper"

	"errors"

	log "github.com/dynastymasra/gochill"
)

type JSONB map[string]interface{}

// Value implement value interface
func (j JSONB) Value() (driver.Value, error) {
	value, err := json.Marshal(j)
	if err != nil {
		log.Error(log.Msg("Failed marshal JSON", err.Error()), log.O("version", config.Version),
			log.O("package", runtime.FuncForPC(reflect.ValueOf(j.Value).Pointer()).Name()),
			log.O("project", config.ProjectName), log.O("data", helper.Stringify(j)))
		return nil, err
	}
	return value, nil
}

func (j *JSONB) Scan(value interface{}) error {
	source, ok := value.([]byte)
	if !ok {
		log.Error(log.Msg("Failed casting data", config.ErrCastingData), log.O("version", config.Version),
			log.O("package", runtime.FuncForPC(reflect.ValueOf(j.Scan).Pointer()).Name()),
			log.O("project", config.ProjectName), log.O("data", helper.Stringify(j)))
		return errors.New(config.ErrCastingData)
	}

	if err := json.Unmarshal(source, &j); err != nil {
		log.Error(log.Msg("Failed unmarshal JSON", err.Error()), log.O("version", config.Version),
			log.O("package", runtime.FuncForPC(reflect.ValueOf(j.Scan).Pointer()).Name()),
			log.O("project", config.ProjectName), log.O("data", string(source)))
		return err
	}
	return nil
}
