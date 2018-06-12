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
	url := config.KongAdminURL + "/consumers"
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
	url := config.KongAdminURL + fmt.Sprintf("/consumers/%v/oauth2", consumerID)
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

		log.Info(log.Msg("Success create kong oauth", helper.Stringify(oauth)),
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

func (kr KongRepository) GetOauthByName(consumerID string, name domain.OauthName) ([]kong.Oauth, int, error) {
	var status int
	response := struct {
		Data []kong.Oauth `json:"data"`
	}{}

	pack := runtime.FuncForPC(reflect.ValueOf(kr.GetOauthByName).Pointer()).Name()
	url := config.KongAdminURL + fmt.Sprintf("/consumers/%v/oauth2", consumerID)
	backOff := BackOffRetry()

	operation := func() error {
		res, body, errs := gorequest.New().Get(url).Param("name", string(name)).EndStruct(&response)
		if len(errs) > 0 {
			log.Error(log.Msg("Error get oauth by name", errs[0].Error()), log.O("version", config.Version),
				log.O("package", pack), log.O("project", config.ProjectName),
				log.O("url", url), log.O(config.TraceKey, kr.Ctx.Value(config.TraceKey)),
				log.O("consumer_id", consumerID), log.O("retry_in", backOff.NextBackOff()))
			status = http.StatusInternalServerError
			return errs[0]
		}

		if res.StatusCode >= http.StatusBadRequest {
			log.Error(log.Msg("Failed get kong oauth by name", string(body)), log.O("version", config.Version),
				log.O("package", pack), log.O("project", config.ProjectName),
				log.O("url", url), log.O(config.TraceKey, kr.Ctx.Value(config.TraceKey)),
				log.O("status_code", res.StatusCode), log.O("consumer_id", consumerID),
				log.O("retry_in", backOff.NextBackOff()))
			status = res.StatusCode
			return errors.New(string(body))
		}

		log.Info(log.Msg("Success get kong oauth", helper.Stringify(response)),
			log.O("version", config.Version), log.O("package", pack), log.O("url", url),
			log.O("project", config.ProjectName), log.O(config.TraceKey, kr.Ctx.Value(config.TraceKey)),
			log.O("status_code", res.StatusCode), log.O("elapsed_time", backOff.GetElapsedTime()),
			log.O("consumer_id", consumerID))
		status = res.StatusCode

		return nil
	}

	if err := backoff.Retry(operation, backOff); err != nil {
		log.Error(log.Msg("Failed retry get kong oauth", err.Error()), log.O("version", config.Version),
			log.O("package", pack), log.O("project", config.ProjectName), log.O("url", url),
			log.O("consumer_id", consumerID), log.O(config.TraceKey, kr.Ctx.Value(config.TraceKey)),
			log.O("elapsed_time", backOff.GetElapsedTime()))
		return nil, status, err
	}

	return response.Data, status, nil
}

func (kr KongRepository) GetAccessToken(clientID, clientSecret, scope, userID string) (*kong.AccessToken, int, error) {
	var status int
	body := struct {
		ClientID            string `json:"client_id"`
		ClientSecret        string `json:"client_secret"`
		GrantType           string `json:"grant_type"`
		Scope               string `json:"scope"`
		ProvisionKey        string `json:"provision_key"`
		AuthenticatedUserID string `json:"authenticated_userid"`
	}{
		ClientID:            clientID,
		ClientSecret:        clientSecret,
		GrantType:           "password",
		Scope:               scope,
		ProvisionKey:        config.ProvisionKey,
		AuthenticatedUserID: userID,
	}

	pack := runtime.FuncForPC(reflect.ValueOf(kr.GetAccessToken).Pointer()).Name()
	response := &kong.AccessToken{}
	url := config.KongAuthURL + "/oauth2/token"
	backOff := BackOffRetry()

	operation := func() error {
		res, body, errs := gorequest.New().Post(url).Send(body).EndStruct(response)
		if len(errs) > 0 {
			log.Error(log.Msg("Error get access token", errs[0].Error()), log.O("version", config.Version),
				log.O("package", pack), log.O("project", config.ProjectName), log.O("scope", scope),
				log.O("url", url), log.O(config.TraceKey, kr.Ctx.Value(config.TraceKey)),
				log.O("retry_in", backOff.NextBackOff()), log.O("body", helper.Stringify(body)))
			status = http.StatusInternalServerError
			return errs[0]
		}

		if res.StatusCode >= http.StatusBadRequest {
			log.Error(log.Msg("Failed access token", string(body)), log.O("version", config.Version),
				log.O("package", pack), log.O("project", config.ProjectName), log.O("scope", scope),
				log.O("url", url), log.O(config.TraceKey, kr.Ctx.Value(config.TraceKey)),
				log.O("status_code", res.StatusCode), log.O("retry_in", backOff.NextBackOff()),
				log.O("body", helper.Stringify(body)))
			status = res.StatusCode
			return errors.New(string(body))
		}

		log.Info(log.Msg("Success get access token", helper.Stringify(response)), log.O("body", helper.Stringify(body)),
			log.O("version", config.Version), log.O("package", pack), log.O("url", url),
			log.O("project", config.ProjectName), log.O(config.TraceKey, kr.Ctx.Value(config.TraceKey)),
			log.O("status_code", res.StatusCode), log.O("elapsed_time", backOff.GetElapsedTime()))
		status = res.StatusCode

		return nil
	}

	if err := backoff.Retry(operation, backOff); err != nil {
		log.Error(log.Msg("Failed retry get access token", err.Error()), log.O("version", config.Version),
			log.O("package", pack), log.O("project", config.ProjectName), log.O("url", url),
			log.O(config.TraceKey, kr.Ctx.Value(config.TraceKey)), log.O("elapsed_time", backOff.GetElapsedTime()),
			log.O("body", helper.Stringify(body)))
		return nil, status, err
	}

	return response, status, nil
}
