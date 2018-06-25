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
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"
)

func LoginController(c *gin.Context) {
	var login actor.Login

	c.Header("Content-Type", "application/json")

	buf := new(bytes.Buffer)
	buf.ReadFrom(c.Request.Body)
	reqBody := buf.String()
	pack := runtime.FuncForPC(reflect.ValueOf(LoginController).Pointer()).Name()

	log.Info(log.Msg("Request login user", reqBody), log.O("version", config.Version),
		log.O("project", config.ProjectName), log.O("package", pack),
		log.O(config.TraceKey, c.GetString(config.TraceKey)))

	if err := json.Unmarshal([]byte(reqBody), &login); err != nil {
		log.Error(log.Msg("Failed unmarshal body", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, c.GetString(config.TraceKey)),
			log.O("package", pack), log.O("body", reqBody))
		c.Error(err)
		c.JSON(http.StatusBadRequest, helper.FailResponse(err.Error()))
		return
	}

	validate := validator.New()
	if err := validate.Struct(login); err != nil {
		log.Error(log.Msg("Failed validate request body", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, c.GetString(config.TraceKey)),
			log.O("package", pack), log.O("body", helper.Stringify(login)))
		c.Error(err)
		c.JSON(http.StatusBadRequest, helper.FailResponse(err.Error()))
		return
	}

	oauthName, err := domain.OauthNameValidation(login.OauthName)
	if err != nil {
		log.Error(log.Msg("Failed oauth validation", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, c.GetString(config.TraceKey)),
			log.O("package", pack), log.O("body", helper.Stringify(login)))
		c.Error(err)
		c.JSON(http.StatusBadRequest, helper.FailResponse(err.Error()))
		return
	}
	login.OauthName = oauthName

	db, err := provider.ConnectSQL()
	if err != nil {
		log.Error(log.Msg("Failed connect database", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, c.GetString(config.TraceKey)),
			log.O("package", pack), log.O("body", helper.Stringify(login)))
		c.Error(err)
		c.JSON(http.StatusInternalServerError, helper.FailResponse(err.Error()))
		return
	}

	actorRepository := sql.NewUserRepository(c, db)

	notExists := actorRepository.CheckEmailNotExist(login.Email)
	if notExists {
		log.Warn(log.Msg("Email is not exists", login.Email), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, c.GetString(config.TraceKey)),
			log.O("package", pack), log.O("body", helper.Stringify(login)))
		errDesc := errors.New(fmt.Sprintf("email %v is not found", login.Email))
		c.Error(errDesc)
		c.JSON(http.StatusNotFound, helper.FailResponse(errDesc.Error()))
		return
	}

	user, err := actorRepository.Login(login.Email)
	if err != nil {
		log.Error(log.Msg("Failed login user", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, c.GetString(config.TraceKey)),
			log.O("package", pack), log.O("body", helper.Stringify(login)))
		c.Error(err)
		c.JSON(http.StatusUnauthorized, helper.FailResponse(err.Error()))
		return
	}

	if match := actor.CheckPasswordHash(login.Password, user.Password); !match {
		log.Error(log.Msg("Password is not match"), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, c.GetString(config.TraceKey)),
			log.O("package", pack), log.O("body", helper.Stringify(login)),
			log.O("user", helper.Stringify(user)))
		err := errors.New("invalid password")
		c.Error(err)
		c.JSON(http.StatusForbidden, helper.FailResponse(err.Error()))
		return
	}

	oauthRepository := api.NewOauthRepository(c)
	oauthService := service.NewOauthService(c, oauthRepository)

	oauth, status, err := oauthService.LoginService(user.ConsumerID, user.ID, login.OauthName)
	if err != nil {
		log.Error(log.Msg("Failed get access token", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, c.GetString(config.TraceKey)),
			log.O("package", pack), log.O("body", helper.Stringify(login)))
		c.Error(err)
		c.JSON(status, helper.FailResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, helper.ObjectResponse(oauth))
}
