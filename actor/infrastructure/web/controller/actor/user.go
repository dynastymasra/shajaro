package actor

import (
	"reflect"
	"runtime"

	"bytes"
	"encoding/json"
	"net/http"

	"fmt"

	"crypto/sha512"

	"shajaro/actor/config"
	"shajaro/actor/domain"
	"shajaro/actor/domain/actor"
	"shajaro/actor/domain/kong"
	"shajaro/actor/helper"
	"shajaro/actor/infrastructure/provider"
	"shajaro/actor/infrastructure/repository/api"
	"shajaro/actor/infrastructure/repository/sql"
	"shajaro/actor/service"

	log "github.com/dynastymasra/gochill"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"gopkg.in/go-playground/validator.v9"
)

func RegisterController(c *gin.Context) {
	var user actor.User

	c.Header("Content-Type", "application/json")

	buf := new(bytes.Buffer)
	buf.ReadFrom(c.Request.Body)
	reqBody := buf.String()
	pack := runtime.FuncForPC(reflect.ValueOf(RegisterController).Pointer()).Name()

	log.Info(log.Msg("Request create user", reqBody), log.O("version", config.Version),
		log.O("project", config.ProjectName), log.O("package", pack),
		log.O(config.TraceKey, c.GetString(config.TraceKey)))

	if err := json.Unmarshal([]byte(reqBody), &user); err != nil {
		log.Error(log.Msg("Failed unmarshal body", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, c.GetString(config.TraceKey)),
			log.O("package", pack), log.O("body", reqBody))
		c.Error(err)
		c.JSON(http.StatusBadRequest, helper.FailResponse(err.Error()))
		return
	}

	validate := validator.New()
	if err := validate.Struct(user); err != nil {
		log.Error(log.Msg("Failed validate request body", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, c.GetString(config.TraceKey)),
			log.O("package", pack), log.O("body", helper.Stringify(user)))
		c.Error(err)
		c.JSON(http.StatusBadRequest, helper.FailResponse(err.Error()))
		return
	}

	gender, err := domain.GenderValidation(user.Gender)
	if err != nil {
		log.Error(log.Msg("Failed validation gender", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, c.GetString(config.TraceKey)),
			log.O("package", pack), log.O("body", helper.Stringify(user)))
		c.Error(err)
		c.JSON(http.StatusBadRequest, helper.FailResponse(err.Error()))
		return
	}
	user.Gender = gender

	db, err := provider.ConnectSQL()
	if err != nil {
		log.Error(log.Msg("Failed connect database", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, c.GetString(config.TraceKey)),
			log.O("package", pack), log.O("body", helper.Stringify(user)))
		c.Error(err)
		c.JSON(http.StatusInternalServerError, helper.FailResponse(config.ErrDatabaseConnectFail))
		return
	}

	actorRepository := sql.NewUserRepository(c, db)

	exits := actorRepository.CheckEmailExist(user.Email)
	if !exits {
		log.Warn(log.Msg("Email already exists", user.Email), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, c.GetString(config.TraceKey)),
			log.O("package", pack), log.O("body", helper.Stringify(user)))
		errDesc := errors.New(fmt.Sprintf("email %v already exists", user.Email))
		c.Error(errDesc)
		c.JSON(http.StatusConflict, helper.FailResponse(errDesc.Error()))
		return
	}

	h := sha512.New512_256()
	user.Password = fmt.Sprintf("%x", h.Sum([]byte(user.Password)))
	user.ID = uuid.NewV4().String()

	kongRepository := api.NewKongRepository(c)
	kongService := service.NewKongService(c, kongRepository, kongRepository)

	consumer := kong.Consumer{
		CustomID: user.ID,
		Username: user.Email,
	}

	consumerResp, status, err := kongService.RegisterNewConsumer(consumer)
	if err != nil {
		log.Error(log.Msg("Failed create kong consumer", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, c.GetString(config.TraceKey)),
			log.O("package", pack), log.O("status_code", status), log.O("body", helper.Stringify(consumer)))
		c.Error(err)
		c.JSON(status, helper.FailResponse(err.Error()))
		return
	}
	user.ConsumerID = consumerResp.ID

	if err := actorRepository.CreateUser(user); err != nil {
		log.Error(log.Msg("Failed create new user", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, c.GetString(config.TraceKey)),
			log.O("package", pack), log.O("body", helper.Stringify(user)))
		c.Error(err)
		c.JSON(status, helper.FailResponse(config.ErrFailedSaveNewUser))
		return
	}

	c.JSON(http.StatusCreated, helper.ObjectResponse(user))
}
