package api

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"runtime"
	"shajaro/actor/config"
	"shajaro/actor/domain/kong"
	"shajaro/actor/helper"

	"fmt"

	"github.com/cenkalti/backoff"
	log "github.com/dynastymasra/gochill"
	"github.com/parnurzeal/gorequest"
)

type ConsumerRepository struct {
	Ctx context.Context
}

func NewConsumerRepository(ctx context.Context) ConsumerRepository {
	return ConsumerRepository{
		Ctx: ctx,
	}
}

func (cr ConsumerRepository) CreateConsumer(consumer kong.Kong) (*kong.Kong, int, error) {
	var status int

	pack := runtime.FuncForPC(reflect.ValueOf(cr.CreateConsumer).Pointer()).Name()
	response := &kong.Kong{}
	url := config.KongAdminURL + "/consumers"
	backOff := BackOffRetry()

	operation := func() error {
		res, body, errs := gorequest.New().Post(url).Send(consumer).EndStruct(response)
		if len(errs) > 0 {
			log.Error(log.Msg("Error create kong consumer", errs[0].Error()), log.O("version", config.Version),
				log.O("package", pack), log.O("project", config.ProjectName),
				log.O("body", helper.Stringify(consumer)), log.O("url", url),
				log.O(config.TraceKey, cr.Ctx.Value(config.TraceKey)), log.O("retry_in", backOff.NextBackOff()))
			status = http.StatusInternalServerError
			return errs[0]
		}

		if res.StatusCode >= http.StatusBadRequest {
			log.Error(log.Msg("Failed create kong consumer", string(body)), log.O("version", config.Version),
				log.O("package", pack), log.O("project", config.ProjectName),
				log.O("body", helper.Stringify(consumer)), log.O("url", url),
				log.O(config.TraceKey, cr.Ctx.Value(config.TraceKey)), log.O("status_code", res.StatusCode),
				log.O("retry_in", backOff.NextBackOff()))
			status = res.StatusCode
			return errors.New(string(body))
		}

		log.Info(log.Msg("Success create new kong consumer", helper.Stringify(consumer)),
			log.O("version", config.Version), log.O("package", pack), log.O("url", url),
			log.O("project", config.ProjectName), log.O(config.TraceKey, cr.Ctx.Value(config.TraceKey)),
			log.O("status_code", res.StatusCode), log.O("elapsed_time", backOff.GetElapsedTime()))
		status = res.StatusCode

		return nil
	}

	if err := backoff.Retry(operation, backOff); err != nil {
		log.Error(log.Msg("Failed retry create kong consumer", err.Error()), log.O("version", config.Version),
			log.O("package", pack), log.O("project", config.ProjectName),
			log.O("body", helper.Stringify(consumer)), log.O("url", url),
			log.O(config.TraceKey, cr.Ctx.Value(config.TraceKey)), log.O("elapsed_time", backOff.GetElapsedTime()))
		return nil, status, err
	}

	return response, status, nil
}

func (cr ConsumerRepository) DeleteConsumer(consumerID string) (int, error) {
	var status int

	pack := runtime.FuncForPC(reflect.ValueOf(cr.DeleteConsumer).Pointer()).Name()
	url := config.KongAdminURL + fmt.Sprintf("/consumers/%v", consumerID)
	backOff := BackOffRetry()

	operation := func() error {
		res, body, errs := gorequest.New().Delete(url).End()
		if len(errs) > 0 {
			log.Error(log.Msg("Error delete kong consumer", errs[0].Error()), log.O("version", config.Version),
				log.O("package", pack), log.O("project", config.ProjectName),
				log.O("consumer_id", consumerID), log.O("url", url),
				log.O(config.TraceKey, cr.Ctx.Value(config.TraceKey)), log.O("retry_in", backOff.NextBackOff()))
			status = http.StatusInternalServerError
			return errs[0]
		}

		if res.StatusCode >= http.StatusBadRequest {
			log.Error(log.Msg("Failed delete kong consumer", body), log.O("version", config.Version),
				log.O("package", pack), log.O("project", config.ProjectName),
				log.O("consumer_id", consumerID), log.O("url", url),
				log.O(config.TraceKey, cr.Ctx.Value(config.TraceKey)), log.O("status_code", res.StatusCode),
				log.O("retry_in", backOff.NextBackOff()))
			status = res.StatusCode
			return errors.New(body)
		}

		log.Info(log.Msg("Success delete kong consumer", body),
			log.O("version", config.Version), log.O("package", pack), log.O("url", url),
			log.O("project", config.ProjectName), log.O(config.TraceKey, cr.Ctx.Value(config.TraceKey)),
			log.O("status_code", res.StatusCode), log.O("elapsed_time", backOff.GetElapsedTime()))
		status = res.StatusCode

		return nil
	}

	if err := backoff.Retry(operation, backOff); err != nil {
		log.Error(log.Msg("Failed retry delete kong consumer", err.Error()), log.O("version", config.Version),
			log.O("package", pack), log.O("project", config.ProjectName),
			log.O("consumer_id", consumerID), log.O("url", url),
			log.O(config.TraceKey, cr.Ctx.Value(config.TraceKey)), log.O("elapsed_time", backOff.GetElapsedTime()))
		return status, err
	}

	return status, nil
}
