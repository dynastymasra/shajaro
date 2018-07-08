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
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"gopkg.in/go-playground/validator.v9"
)

func RegisterController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user actor.Actor

	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	reqBody := buf.String()
	traceKey := r.Context().Value(config.TraceKey)
	pack := runtime.FuncForPC(reflect.ValueOf(RegisterController).Pointer()).Name()

	log.Info(log.Msg("Request create user", reqBody), log.O("version", config.Version),
		log.O("project", config.ProjectName), log.O("package", pack),
		log.O(config.TraceKey, traceKey))

	if err := json.Unmarshal([]byte(reqBody), &user); err != nil {
		log.Error(log.Msg("Failed unmarshal body", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, traceKey),
			log.O("package", pack), log.O("body", reqBody))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, helper.FailResponse(err.Error()).Stringify())
		return
	}

	validate := validator.New()
	if err := validate.Struct(user); err != nil {
		log.Error(log.Msg("Failed validate request body", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, traceKey),
			log.O("package", pack), log.O("body", helper.Stringify(user)))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, helper.FailResponse(err.Error()).Stringify())
		return
	}

	gender, err := domain.GenderValidation(user.Gender)
	if err != nil {
		log.Error(log.Msg("Failed validation gender", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, traceKey),
			log.O("package", pack), log.O("body", helper.Stringify(user)))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, helper.FailResponse(err.Error()).Stringify())
		return
	}
	user.Gender = gender

	db, err := provider.ConnectSQL()
	if err != nil {
		log.Error(log.Msg("Failed connect database", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, traceKey),
			log.O("package", pack), log.O("body", helper.Stringify(user)))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, helper.FailResponse(err.Error()).Stringify())
		return
	}

	userRepository := sql.NewUserRepository(r.Context(), db)

	notExists := userRepository.CheckEmailNotExist(user.Email)
	if !notExists {
		log.Warn(log.Msg("Email already exists", user.Email), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, traceKey),
			log.O("package", pack), log.O("body", helper.Stringify(user)))
		w.WriteHeader(http.StatusConflict)
		fmt.Fprintf(w, helper.FailResponse(fmt.Sprintf("email %v already exists", user.Email)).Stringify())
		return
	}

	password, err := actor.HashPassword(user.Password)
	if err != nil {
		log.Error(log.Msg("Failed to hash password", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, traceKey),
			log.O("package", pack), log.O("body", helper.Stringify(user)))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, helper.FailResponse(err.Error()).Stringify())
		return
	}
	user.ID = uuid.NewV4().String()
	user.Password = password

	consumerRepository := api.NewConsumerRepository(r.Context())
	consumer := kong.Kong{
		CustomID: user.ID,
		Username: user.Email,
	}

	consumerResp, status, err := consumerRepository.CreateConsumer(consumer)
	if err != nil {
		log.Error(log.Msg("Failed create kong consumer", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, traceKey),
			log.O("package", pack), log.O("status_code", status), log.O("body", helper.Stringify(consumer)))
		w.WriteHeader(status)
		fmt.Fprintf(w, helper.FailResponse(err.Error()).Stringify())
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
				log.O("project", config.ProjectName), log.O(config.TraceKey, ctx.Value(config.TraceKey)),
				log.O("package", naming), log.O("consumer_id", consumerID),
				log.O("body", helper.Stringify(oauth)))

			oauthRepository := api.NewOauthRepository(ctx)
			res, status, err := oauthRepository.CreateOauth(consumerID, oauth)
			if err != nil {
				log.Error(log.Msg("Failed create kong oauth", err.Error()), log.O("version", config.Version),
					log.O("project", config.ProjectName), log.O(config.TraceKey, ctx.Value(config.TraceKey)),
					log.O("package", naming), log.O("body", helper.Stringify(oauth)),
					log.O("status_code", status))
				return
			}

			log.Info(log.Msg("Success create kong oauth", helper.Stringify(res)), log.O("version", config.Version),
				log.O("project", config.ProjectName), log.O(config.TraceKey, ctx.Value(config.TraceKey)),
				log.O("package", naming), log.O("consumer_id", consumerID), log.O("status_code", status),
				log.O("body", helper.Stringify(oauth)))
		}(r.Context(), pack, consumerResp.ID, auth)
	}

	if err := userRepository.Create(user); err != nil {
		log.Error(log.Msg("Failed create new user", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, traceKey),
			log.O("package", pack), log.O("body", helper.Stringify(user)))
		w.WriteHeader(status)
		fmt.Fprintf(w, helper.FailResponse(err.Error()).Stringify())
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, helper.ObjectResponse(user).Stringify())
}

func GetUserByIDController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	traceKey := r.Context().Value(config.TraceKey)
	pack := runtime.FuncForPC(reflect.ValueOf(GetUserByIDController).Pointer()).Name()
	id := r.Header.Get(config.AuthUserIDHeader)

	log.Info(log.Msg("Get user by id"), log.O("version", config.Version), log.O("id", id),
		log.O("project", config.ProjectName), log.O("package", pack),
		log.O(config.TraceKey, traceKey))

	db, err := provider.ConnectSQL()
	if err != nil {
		log.Error(log.Msg("Failed connect database", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, traceKey),
			log.O("package", pack))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, helper.FailResponse(err.Error()).Stringify())
		return
	}

	userRepository := sql.NewUserRepository(r.Context(), db)
	user, err := userRepository.GetUserByID(id)
	if err == gorm.ErrRecordNotFound {
		log.Error(log.Msg("User not found", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, traceKey),
			log.O("package", pack), log.O("id", id))
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, helper.FailResponse(err.Error()).Stringify())
		return
	}

	if err != nil {
		log.Error(log.Msg("Failed get user by id", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, traceKey),
			log.O("package", pack), log.O("id", id))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, helper.FailResponse(err.Error()).Stringify())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, helper.ObjectResponse(user).Stringify())
}

func UpdateUserController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user actor.Actor

	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	reqBody := buf.String()
	traceKey := r.Context().Value(config.TraceKey)
	pack := runtime.FuncForPC(reflect.ValueOf(UpdateUserController).Pointer()).Name()
	id := r.Header.Get(config.AuthUserIDHeader)
	consumerId := r.Header.Get(config.ConsumerIDHeader)

	log.Info(log.Msg("Request to update user", reqBody), log.O("version", config.Version),
		log.O("project", config.ProjectName), log.O("package", pack),
		log.O(config.TraceKey, traceKey))

	if err := json.Unmarshal([]byte(reqBody), &user); err != nil {
		log.Error(log.Msg("Failed unmarshal body", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, traceKey),
			log.O("package", pack), log.O("body", reqBody))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, helper.FailResponse(err.Error()).Stringify())
		return
	}

	validate := validator.New()
	if err := validate.Struct(user); err != nil {
		log.Error(log.Msg("Failed validate request body", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, traceKey),
			log.O("package", pack), log.O("body", helper.Stringify(user)))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, helper.FailResponse(err.Error()).Stringify())
		return
	}

	gender, err := domain.GenderValidation(user.Gender)
	if err != nil {
		log.Error(log.Msg("Failed validation gender", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, traceKey),
			log.O("package", pack), log.O("body", helper.Stringify(user)))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, helper.FailResponse(err.Error()).Stringify())
		return
	}
	user.Gender = gender

	db, err := provider.ConnectSQL()
	if err != nil {
		log.Error(log.Msg("Failed connect database", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, traceKey),
			log.O("package", pack), log.O("body", helper.Stringify(user)))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, helper.FailResponse(err.Error()).Stringify())
		return
	}

	user.ID = id
	user.ConsumerID = consumerId

	userRepository := sql.NewUserRepository(r.Context(), db)
	userService := service.NewUserService(r.Context(), userRepository)
	result, err := userService.UpdateService(user)
	if err == gorm.ErrRecordNotFound {
		log.Error(log.Msg("Data not found", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, traceKey),
			log.O("package", pack), log.O("body", helper.Stringify(user)))
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, helper.FailResponse(err.Error()).Stringify())
		return
	}

	if err != nil {
		log.Error(log.Msg("Failed update data user", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, traceKey),
			log.O("package", pack), log.O("body", helper.Stringify(user)))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, helper.FailResponse(err.Error()).Stringify())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, helper.ObjectResponse(result).Stringify())
}

func DeleteUserController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	pack := runtime.FuncForPC(reflect.ValueOf(DeleteUserController).Pointer()).Name()
	id := r.Header.Get(config.AuthUserIDHeader)
	consumerId := r.Header.Get(config.ConsumerIDHeader)
	traceKey := r.Header.Get(config.TraceKey)

	log.Info(log.Msg("Delete user by id"), log.O("version", config.Version), log.O("id", id),
		log.O("project", config.ProjectName), log.O("package", pack), log.O("consumer_id", consumerId),
		log.O(config.TraceKey, traceKey))

	db, err := provider.ConnectSQL()
	if err != nil {
		log.Error(log.Msg("Failed connect database", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, traceKey),
			log.O("package", pack), log.O("consumer_id", consumerId), log.O("id", id))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, helper.FailResponse(err.Error()).Stringify())
		return
	}

	userRepository := sql.NewUserRepository(r.Context(), db)
	userService := service.NewUserService(r.Context(), userRepository)
	err = userService.DeleteService(id)
	if err == gorm.ErrRecordNotFound {
		log.Error(log.Msg("User not found", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, traceKey),
			log.O("package", pack), log.O("id", id), log.O("consumer_id", consumerId))
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, helper.FailResponse(err.Error()).Stringify())
		return
	}

	if err != nil {
		log.Error(log.Msg("Failed delete user", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, traceKey),
			log.O("package", pack), log.O("id", id), log.O("consumer_id", consumerId))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, helper.FailResponse(err.Error()).Stringify())
		return
	}

	consumerRepository := api.NewConsumerRepository(r.Context())
	status, err := consumerRepository.DeleteConsumer(consumerId)
	if err != nil {
		log.Error(log.Msg("Failed delete consumer", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, traceKey),
			log.O("package", pack), log.O("id", id), log.O("consumer_id", consumerId))
		w.WriteHeader(status)
		fmt.Fprintf(w, helper.FailResponse(err.Error()).Stringify())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, helper.SuccessResponse().Stringify())
}
