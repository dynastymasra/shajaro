package api

import (
	"context"

	"reflect"
	"runtime"

	"errors"
	"net/http"

	"fmt"

	"shajaro/actor/config"
	"shajaro/actor/domain/kong"
	"shajaro/actor/helper"

	"shajaro/actor/domain"

	"github.com/cenkalti/backoff"
	log "github.com/dynastymasra/gochill"
	"github.com/parnurzeal/gorequest"
)

type KongRepository struct {
	Ctx context.Context
}

func NewKongRepository(ctx context.Context) KongRepository {
	return KongRepository{
		Ctx: ctx,
	}
}

func (kr KongRepository) CreateConsumer(consumer kong.Consumer) (*kong.Consumer, int, error) {
	var status int

	pack := runtime.FuncForPC(reflect.ValueOf(kr.CreateConsumer).Pointer()).Name()
	response := &kong.Consumer{}
	url := config.KongURL + "/consumers"
	backOff := BackOffRetry()

	operation := func() error {
		res, body, errs := gorequest.New().Post(url).Send(consumer).EndStruct(response)
		if len(errs) > 0 {
			log.Error(log.Msg("Error create kong consumer", errs[0].Error()), log.O("version", config.Version),
				log.O("package", pack), log.O("project", config.ProjectName),
				log.O("body", helper.Stringify(consumer)), log.O("url", url),
				log.O(config.TraceKey, kr.Ctx.Value(config.TraceKey)), log.O("retry_in", backOff.NextBackOff()))
			status = http.StatusInternalServerError
			return errs[0]
		}

		if res.StatusCode >= http.StatusBadRequest {
			log.Error(log.Msg("Failed create kong consumer", string(body)), log.O("version", config.Version),
				log.O("package", pack), log.O("project", config.ProjectName),
				log.O("body", helper.Stringify(consumer)), log.O("url", url),
				log.O(config.TraceKey, kr.Ctx.Value(config.TraceKey)), log.O("status_code", res.StatusCode),
				log.O("retry_in", backOff.NextBackOff()))
			status = res.StatusCode
			return errors.New(string(body))
		}

		log.Info(log.Msg("Success create new kong consumer", helper.Stringify(consumer)),
			log.O("version", config.Version), log.O("package", pack), log.O("url", url),
			log.O("project", config.ProjectName), log.O(config.TraceKey, kr.Ctx.Value(config.TraceKey)),
			log.O("status_code", res.StatusCode), log.O("elapsed_time", backOff.GetElapsedTime()))
		status = res.StatusCode

		return nil
	}

	if err := backoff.Retry(operation, backOff); err != nil {
		log.Error(log.Msg("Failed retry create kong consumer", err.Error()), log.O("version", config.Version),
			log.O("package", pack), log.O("project", config.ProjectName),
			log.O("body", helper.Stringify(consumer)), log.O("url", url),
			log.O(config.TraceKey, kr.Ctx.Value(config.TraceKey)), log.O("elapsed_time", backOff.GetElapsedTime()))
		return nil, status, err
	}

	return response, status, nil
}

func (kr KongRepository) CreateOauth(consumerID string, oauth kong.Oauth) (*kong.Oauth, int, error) {
	var status int

	pack := runtime.FuncForPC(reflect.ValueOf(kr.CreateOauth).Pointer()).Name()
	response := &kong.Oauth{}
	url := config.KongURL + fmt.Sprintf("/consumers/%v/oauth2", consumerID)
	backOff := BackOffRetry()

	operation := func() error {
		res, body, errs := gorequest.New().Post(url).Send(oauth).EndStruct(response)
		if len(errs) > 0 {
			log.Error(log.Msg("Error create kong oauth", errs[0].Error()), log.O("version", config.Version),
				log.O("package", pack), log.O("project", config.ProjectName),
				log.O("body", helper.Stringify(oauth)), log.O("url", url),
				log.O(config.TraceKey, kr.Ctx.Value(config.TraceKey)), log.O("consumer_id", consumerID),
				log.O("retry_in", backOff.NextBackOff()))
			status = http.StatusInternalServerError
			return errs[0]
		}

		if res.StatusCode >= http.StatusBadRequest {
			log.Error(log.Msg("Failed create kong oauth", string(body)), log.O("version", config.Version),
				log.O("package", pack), log.O("project", config.ProjectName),
				log.O("body", helper.Stringify(oauth)), log.O("url", url),
				log.O(config.TraceKey, kr.Ctx.Value(config.TraceKey)), log.O("status_code", res.StatusCode),
				log.O("consumer_id", consumerID), log.O("retry_in", backOff.NextBackOff()))
			status = res.StatusCode
			return errors.New(string(body))
		}

		log.Info(log.Msg("Success create new kong consumer", helper.Stringify(oauth)),
			log.O("version", config.Version), log.O("package", pack), log.O("url", url),
			log.O("project", config.ProjectName), log.O(config.TraceKey, kr.Ctx.Value(config.TraceKey)),
			log.O("status_code", res.StatusCode), log.O("elapsed_time", backOff.GetElapsedTime()),
			log.O("consumer_id", consumerID))
		status = res.StatusCode

		return nil
	}

	if err := backoff.Retry(operation, backOff); err != nil {
		log.Error(log.Msg("Failed retry create kong oauth", err.Error()), log.O("version", config.Version),
			log.O("package", pack), log.O("project", config.ProjectName),
			log.O("body", helper.Stringify(oauth)), log.O("url", url), log.O("consumer_id", consumerID),
			log.O(config.TraceKey, kr.Ctx.Value(config.TraceKey)), log.O("elapsed_time", backOff.GetElapsedTime()))
		return nil, status, err
	}

	return response, status, nil
}

func (kr KongRepository) GetOauthByName(consumer string, name domain.OauthName) ([]kong.Oauth, int, error) {
	return nil, 0, nil
}
