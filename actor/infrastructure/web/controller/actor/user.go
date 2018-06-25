package actor

import (
	"reflect"
	"runtime"

	"bytes"
	"encoding/json"
	"net/http"

	"fmt"

	"shajaro/actor/config"
	"shajaro/actor/domain"
	"shajaro/actor/domain/actor"
	"shajaro/actor/domain/kong"
	"shajaro/actor/helper"
	"shajaro/actor/infrastructure/provider"
	"shajaro/actor/infrastructure/repository/api"
	"shajaro/actor/infrastructure/repository/sql"
	"shajaro/actor/service"

	"context"

	log "github.com/dynastymasra/gochill"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"gopkg.in/go-playground/validator.v9"
)

func RegisterController(c *gin.Context) {
	var user actor.Actor

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
		c.JSON(http.StatusInternalServerError, helper.FailResponse(err.Error()))
		return
	}

	userRepository := sql.NewUserRepository(c, db)

	notExists := userRepository.CheckEmailNotExist(user.Email)
	if !notExists {
		log.Warn(log.Msg("Email already exists", user.Email), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, c.GetString(config.TraceKey)),
			log.O("package", pack), log.O("body", helper.Stringify(user)))
		errDesc := errors.New(fmt.Sprintf("email %v already exists", user.Email))
		c.Error(errDesc)
		c.JSON(http.StatusConflict, helper.FailResponse(errDesc.Error()))
		return
	}

	password, err := actor.HashPassword(user.Password)
	if err != nil {
		log.Error(log.Msg("Failed to hash password", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, c.GetString(config.TraceKey)),
			log.O("package", pack), log.O("body", helper.Stringify(user)))
		c.Error(err)
		c.JSON(http.StatusBadRequest, helper.FailResponse(err.Error()))
		return
	}
	user.ID = uuid.NewV4().String()
	user.Password = password

	consumerRepository := api.NewConsumerRepository(c)
	consumer := kong.Kong{
		CustomID: user.ID,
		Username: user.Email,
	}

	consumerResp, status, err := consumerRepository.CreateConsumer(consumer)
	if err != nil {
		log.Error(log.Msg("Failed create kong consumer", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, c.GetString(config.TraceKey)),
			log.O("package", pack), log.O("status_code", status), log.O("body", helper.Stringify(consumer)))
		c.Error(err)
		c.JSON(status, helper.FailResponse(err.Error()))
		return
	}
	user.ConsumerID = consumerResp.ID

	for _, name := range domain.OauthNames {
		auth := kong.Oauth{
			Name:        name,
			RedirectURI: []string{config.RedirectURI},
		}

		go func(ctx context.Context, naming, consumerID string, oauth kong.Oauth) {
			log.Info(log.Msg("Prepare create kong oauth"), log.O("version", config.Version),
				log.O("project", config.ProjectName), log.O(config.TraceKey, c.GetString(config.TraceKey)),
				log.O("package", naming), log.O("consumer_id", consumerID),
				log.O("body", helper.Stringify(oauth)))

			oauthRepository := api.NewOauthRepository(ctx)
			res, status, err := oauthRepository.CreateOauth(consumerID, oauth)
			if err != nil {
				log.Error(log.Msg("Failed create kong oauth", err.Error()), log.O("version", config.Version),
					log.O("project", config.ProjectName), log.O(config.TraceKey, c.GetString(config.TraceKey)),
					log.O("package", naming), log.O("body", helper.Stringify(oauth)),
					log.O("status_code", status))
				return
			}

			log.Info(log.Msg("Success create kong oauth", helper.Stringify(res)), log.O("version", config.Version),
				log.O("project", config.ProjectName), log.O(config.TraceKey, c.GetString(config.TraceKey)),
				log.O("package", naming), log.O("consumer_id", consumerID), log.O("status_code", status),
				log.O("body", helper.Stringify(oauth)))
		}(c, pack, consumerResp.ID, auth)
	}

	if err := userRepository.Create(user); err != nil {
		log.Error(log.Msg("Failed create new user", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, c.GetString(config.TraceKey)),
			log.O("package", pack), log.O("body", helper.Stringify(user)))
		c.Error(err)
		c.JSON(status, helper.FailResponse(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, helper.ObjectResponse(user))
}

func GetUserByIDController(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	pack := runtime.FuncForPC(reflect.ValueOf(GetUserByIDController).Pointer()).Name()
	id := c.GetHeader(config.AuthUserIDHeader)

	log.Info(log.Msg("Get user by id"), log.O("version", config.Version),
		log.O("project", config.ProjectName), log.O("package", pack),
		log.O(config.TraceKey, c.GetString(config.TraceKey)), log.O("id", id))

	db, err := provider.ConnectSQL()
	if err != nil {
		log.Error(log.Msg("Failed connect database", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, c.GetString(config.TraceKey)),
			log.O("package", pack))
		c.Error(err)
		c.JSON(http.StatusInternalServerError, helper.FailResponse(err.Error()))
		return
	}

	userRepository := sql.NewUserRepository(c, db)
	user, err := userRepository.GetUserByID(id)
	if err == gorm.ErrRecordNotFound {
		log.Error(log.Msg("User not found", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, c.GetString(config.TraceKey)),
			log.O("package", pack), log.O("id", id))
		c.Error(err)
		c.JSON(http.StatusNotFound, helper.FailResponse(err.Error()))
		return
	}

	if err != nil {
		log.Error(log.Msg("Failed get user by id", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, c.GetString(config.TraceKey)),
			log.O("package", pack), log.O("id", id))
		c.Error(err)
		c.JSON(http.StatusInternalServerError, helper.FailResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, helper.ObjectResponse(user))
}

func UpdateUserController(c *gin.Context) {
	var user actor.Actor

	c.Header("Content-Type", "application/json")

	buf := new(bytes.Buffer)
	buf.ReadFrom(c.Request.Body)
	reqBody := buf.String()
	pack := runtime.FuncForPC(reflect.ValueOf(UpdateUserController).Pointer()).Name()
	id := c.GetHeader(config.AuthUserIDHeader)
	consumerId := c.GetHeader(config.ConsumerIDHeader)

	log.Info(log.Msg("Request to update user", reqBody), log.O("version", config.Version),
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
		c.JSON(http.StatusInternalServerError, helper.FailResponse(err.Error()))
		return
	}

	user.ID = id
	user.ConsumerID = consumerId

	userRepository := sql.NewUserRepository(c, db)
	userService := service.NewUserService(c, userRepository)
	result, err := userService.UpdateService(user)
	if err == gorm.ErrRecordNotFound {
		log.Error(log.Msg("Data not found", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, c.GetString(config.TraceKey)),
			log.O("package", pack), log.O("body", helper.Stringify(user)))
		c.Error(err)
		c.JSON(http.StatusNotFound, helper.FailResponse(err.Error()))
		return
	}

	if err != nil {
		log.Error(log.Msg("Failed update data user", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, c.GetString(config.TraceKey)),
			log.O("package", pack), log.O("body", helper.Stringify(user)))
		c.Error(err)
		c.JSON(http.StatusInternalServerError, helper.FailResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, helper.ObjectResponse(result))
}

func DeleteUserController(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	pack := runtime.FuncForPC(reflect.ValueOf(DeleteUserController).Pointer()).Name()
	id := c.GetHeader(config.AuthUserIDHeader)
	consumerId := c.GetHeader(config.ConsumerIDHeader)

	log.Info(log.Msg("Delete user by id"), log.O("version", config.Version),
		log.O("project", config.ProjectName), log.O("package", pack), log.O("consumer_id", consumerId),
		log.O(config.TraceKey, c.GetString(config.TraceKey)), log.O("id", id))

	db, err := provider.ConnectSQL()
	if err != nil {
		log.Error(log.Msg("Failed connect database", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, c.GetString(config.TraceKey)),
			log.O("package", pack), log.O("consumer_id", consumerId), log.O("id", id))
		c.Error(err)
		c.JSON(http.StatusInternalServerError, helper.FailResponse(err.Error()))
		return
	}

	userRepository := sql.NewUserRepository(c, db)
	userService := service.NewUserService(c, userRepository)
	err = userService.DeleteService(id)
	if err == gorm.ErrRecordNotFound {
		log.Error(log.Msg("User not found", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, c.GetString(config.TraceKey)),
			log.O("package", pack), log.O("id", id), log.O("consumer_id", consumerId))
		c.Error(err)
		c.JSON(http.StatusNotFound, helper.FailResponse(err.Error()))
		return
	}

	if err != nil {
		log.Error(log.Msg("Failed delete user", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, c.GetString(config.TraceKey)),
			log.O("package", pack), log.O("id", id), log.O("consumer_id", consumerId))
		c.Error(err)
		c.JSON(http.StatusInternalServerError, helper.FailResponse(err.Error()))
		return
	}

	consumerRepository := api.NewConsumerRepository(c)
	status, err := consumerRepository.DeleteConsumer(consumerId)
	if err != nil {
		log.Error(log.Msg("Failed delete consumer", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, c.GetString(config.TraceKey)),
			log.O("package", pack), log.O("id", id), log.O("consumer_id", consumerId))
		c.Error(err)
		c.JSON(status, helper.FailResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, helper.SuccessResponse())
}
