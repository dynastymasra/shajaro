package actor

import (
	"bytes"
	"reflect"
	"runtime"
	"shajaro/actor/config"
	"shajaro/actor/domain/actor"

	"encoding/json"
	"net/http"
	"shajaro/actor/helper"

	"fmt"
	"shajaro/actor/infrastructure/provider"
	"shajaro/actor/infrastructure/repository/sql"

	"errors"

	"shajaro/actor/domain"

	"shajaro/actor/infrastructure/repository/api"

	"shajaro/actor/service"

	log "github.com/dynastymasra/gochill"
	"gopkg.in/go-playground/validator.v9"
)

func LoginController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var login actor.Login

	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	reqBody := buf.String()
	traceKey := r.Context().Value(config.TraceKey)
	pack := runtime.FuncForPC(reflect.ValueOf(LoginController).Pointer()).Name()

	log.Info(log.Msg("Request login user", reqBody), log.O("version", config.Version),
		log.O("project", config.ProjectName), log.O("package", pack),
		log.O(config.TraceKey, traceKey))

	if err := json.Unmarshal([]byte(reqBody), &login); err != nil {
		log.Error(log.Msg("Failed unmarshal body", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, traceKey),
			log.O("package", pack), log.O("body", reqBody))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, helper.FailResponse(err.Error()).Stringify())
		return
	}

	validate := validator.New()
	if err := validate.Struct(login); err != nil {
		log.Error(log.Msg("Failed validate request body", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, traceKey),
			log.O("package", pack), log.O("body", helper.Stringify(login)))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, helper.FailResponse(err.Error()).Stringify())
		return
	}

	oauthName, err := domain.OauthNameValidation(login.OauthName)
	if err != nil {
		log.Error(log.Msg("Failed oauth validation", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, traceKey),
			log.O("package", pack), log.O("body", helper.Stringify(login)))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, helper.FailResponse(err.Error()).Stringify())
		return
	}
	login.OauthName = oauthName

	db, err := provider.ConnectSQL()
	if err != nil {
		log.Error(log.Msg("Failed connect database", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, traceKey),
			log.O("package", pack), log.O("body", helper.Stringify(login)))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, helper.FailResponse(err.Error()).Stringify())
		return
	}

	actorRepository := sql.NewUserRepository(r.Context(), db)

	notExists := actorRepository.CheckEmailNotExist(login.Email)
	if notExists {
		log.Warn(log.Msg("Email is not exists", login.Email), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, traceKey),
			log.O("package", pack), log.O("body", helper.Stringify(login)))
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, helper.FailResponse(fmt.Sprintf("email %v is not found", login.Email)).Stringify())
		return
	}

	user, err := actorRepository.Login(login.Email)
	if err != nil {
		log.Error(log.Msg("Failed login user", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, traceKey),
			log.O("package", pack), log.O("body", helper.Stringify(login)))
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, helper.FailResponse(err.Error()).Stringify())
		return
	}

	if match := actor.CheckPasswordHash(login.Password, user.Password); !match {
		log.Error(log.Msg("Password is not match"), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, traceKey),
			log.O("package", pack), log.O("body", helper.Stringify(login)),
			log.O("user", helper.Stringify(user)))
		err := errors.New("invalid password")
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, helper.FailResponse(err.Error()).Stringify())
		return
	}

	oauthRepository := api.NewOauthRepository(r.Context())
	oauthService := service.NewOauthService(r.Context(), oauthRepository)

	oauth, status, err := oauthService.LoginService(user.ConsumerID, user.ID, login.OauthName)
	if err != nil {
		log.Error(log.Msg("Failed get access token", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, traceKey),
			log.O("package", pack), log.O("body", helper.Stringify(login)))
		w.WriteHeader(status)
		fmt.Fprintf(w, helper.FailResponse(err.Error()).Stringify())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, helper.ObjectResponse(oauth).Stringify())
}
