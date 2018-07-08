package actor

import (
	"bytes"
	"encoding/json"
	"net/http"
	"reflect"
	"runtime"
	"shajaro/actor/config"
	"shajaro/actor/helper"
	"shajaro/actor/infrastructure/provider"

	"shajaro/actor/domain/actor"

	"errors"
	"shajaro/actor/infrastructure/repository/sql"

	"fmt"

	log "github.com/dynastymasra/gochill"
	"gopkg.in/go-playground/validator.v9"
)

type password struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}

func UpdatePasswordController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var updatePassword password

	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	reqBody := buf.String()
	id := r.Header.Get(config.AuthUserIDHeader)
	traceKey := r.Context().Value(config.TraceKey)
	pack := runtime.FuncForPC(reflect.ValueOf(UpdatePasswordController).Pointer()).Name()

	log.Info(log.Msg("Request update password", reqBody), log.O("version", config.Version),
		log.O("project", config.ProjectName), log.O("package", pack), log.O(config.TraceKey, traceKey))

	if err := json.Unmarshal([]byte(reqBody), &updatePassword); err != nil {
		log.Error(log.Msg("Failed unmarshal body", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, traceKey),
			log.O("package", pack), log.O("body", reqBody))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, helper.FailResponse(err.Error()).Stringify())
		return
	}

	validate := validator.New()
	if err := validate.Struct(updatePassword); err != nil {
		log.Error(log.Msg("Failed validate request body", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, traceKey),
			log.O("package", pack), log.O("body", helper.Stringify(updatePassword)))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, helper.FailResponse(err.Error()).Stringify())
		return
	}

	password, err := actor.HashPassword(updatePassword.NewPassword)
	if err != nil {
		log.Error(log.Msg("Failed to hash password", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, traceKey),
			log.O("package", pack), log.O("body", helper.Stringify(updatePassword)))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, helper.FailResponse(err.Error()).Stringify())
		return
	}

	db, err := provider.ConnectSQL()
	if err != nil {
		log.Error(log.Msg("Failed connect database", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, traceKey),
			log.O("package", pack), log.O("body", helper.Stringify(updatePassword)))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, helper.FailResponse(err.Error()).Stringify())
		return
	}

	userRepository := sql.NewUserRepository(r.Context(), db)

	user, err := userRepository.GetUserByID(id)
	if err != nil {
		log.Error(log.Msg("Failed get user by id", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, traceKey),
			log.O("package", pack), log.O("id", id))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, helper.FailResponse(err.Error()).Stringify())
		return
	}

	if match := actor.CheckPasswordHash(updatePassword.OldPassword, user.Password); !match {
		err := errors.New("invalid password")
		log.Error(log.Msg("Password is not match", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, traceKey),
			log.O("package", pack), log.O("old_password", updatePassword.OldPassword),
			log.O("user", helper.Stringify(user)))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, helper.FailResponse(err.Error()).Stringify())
		return
	}

	if err := userRepository.UpdatePassword(*user, password); err != nil {
		log.Error(log.Msg("Failed update password", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, traceKey),
			log.O("package", pack), log.O("new_password", password),
			log.O("user", helper.Stringify(user)))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, helper.FailResponse(err.Error()).Stringify())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, helper.SuccessResponse().Stringify())
}
