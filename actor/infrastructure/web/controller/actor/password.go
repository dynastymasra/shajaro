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

	log "github.com/dynastymasra/gochill"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"
)

type password struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}

func UpdatePasswordController(c *gin.Context) {
	var updatePassword password

	c.Header("Content-Type", "application/json")

	buf := new(bytes.Buffer)
	buf.ReadFrom(c.Request.Body)
	reqBody := buf.String()
	id := c.GetHeader(config.AuthUserIDHeader)
	pack := runtime.FuncForPC(reflect.ValueOf(UpdatePasswordController).Pointer()).Name()

	log.Info(log.Msg("Request update password", reqBody), log.O("version", config.Version),
		log.O("project", config.ProjectName), log.O("package", pack),
		log.O(config.TraceKey, c.GetString(config.TraceKey)))

	if err := json.Unmarshal([]byte(reqBody), &updatePassword); err != nil {
		log.Error(log.Msg("Failed unmarshal body", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, c.GetString(config.TraceKey)),
			log.O("package", pack), log.O("body", reqBody))
		c.Error(err)
		c.JSON(http.StatusBadRequest, helper.FailResponse(err.Error()))
		return
	}

	validate := validator.New()
	if err := validate.Struct(updatePassword); err != nil {
		log.Error(log.Msg("Failed validate request body", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, c.GetString(config.TraceKey)),
			log.O("package", pack), log.O("body", helper.Stringify(updatePassword)))
		c.Error(err)
		c.JSON(http.StatusBadRequest, helper.FailResponse(err.Error()))
		return
	}

	password, err := actor.HashPassword(updatePassword.NewPassword)
	if err != nil {
		log.Error(log.Msg("Failed to hash password", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, c.GetString(config.TraceKey)),
			log.O("package", pack), log.O("body", helper.Stringify(updatePassword)))
		c.Error(err)
		c.JSON(http.StatusBadRequest, helper.FailResponse(err.Error()))
		return
	}

	db, err := provider.ConnectSQL()
	if err != nil {
		log.Error(log.Msg("Failed connect database", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, c.GetString(config.TraceKey)),
			log.O("package", pack), log.O("body", helper.Stringify(updatePassword)))
		c.Error(err)
		c.JSON(http.StatusInternalServerError, helper.FailResponse(err.Error()))
		return
	}

	userRepository := sql.NewUserRepository(c, db)

	user, err := userRepository.GetUserByID(id)
	if err != nil {
		log.Error(log.Msg("Failed get user by id", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, c.GetString(config.TraceKey)),
			log.O("package", pack), log.O("id", id))
		c.Error(err)
		c.JSON(http.StatusInternalServerError, helper.FailResponse(err.Error()))
		return
	}

	if match := actor.CheckPasswordHash(updatePassword.OldPassword, user.Password); !match {
		err := errors.New("invalid password")
		log.Error(log.Msg("Password is not match", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, c.GetString(config.TraceKey)),
			log.O("package", pack), log.O("old_password", updatePassword.OldPassword),
			log.O("user", helper.Stringify(user)))
		c.Error(err)
		c.JSON(http.StatusBadRequest, helper.FailResponse(err.Error()))
		return
	}

	if err := userRepository.UpdatePassword(*user, password); err != nil {
		log.Error(log.Msg("Failed update password", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, c.GetString(config.TraceKey)),
			log.O("package", pack), log.O("new_password", password),
			log.O("user", helper.Stringify(user)))
		c.Error(err)
		c.JSON(http.StatusInternalServerError, helper.FailResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, helper.SuccessResponse())
}
