package service

import (
	"context"
	"reflect"
	"runtime"

	"github.com/dynastymasra/shajaro/actor/config"
	"github.com/dynastymasra/shajaro/actor/domain"
	"github.com/dynastymasra/shajaro/actor/domain/kong"
	"github.com/dynastymasra/shajaro/actor/helper"

	log "github.com/dynastymasra/gochill"
)

type OauthService struct {
	Ctx     context.Context
	Oauther kong.Oauther
}

func NewOauthService(ctx context.Context, oauther kong.Oauther) OauthService {
	return OauthService{
		Ctx:     ctx,
		Oauther: oauther,
	}
}

func (os *OauthService) LoginService(consumerID, userID string, oauthName domain.OauthName) (*kong.AccessToken, int, error) {
	pack := runtime.FuncForPC(reflect.ValueOf(os.LoginService).Pointer()).Name()

	res, status, err := os.Oauther.GetOauthByName(consumerID, oauthName)
	if err != nil {
		log.Error(log.Msg("Failed get kong oauth", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, os.Ctx.Value(config.TraceKey)),
			log.O("package", pack), log.O("consumer_id", consumerID), log.O("oauth_name", oauthName))
		return nil, status, err
	}

	if len(res) < 1 {
		auth := kong.Oauth{
			Name:        oauthName,
			RedirectURI: []string{config.RedirectURI},
		}

		log.Info(log.Msg("Kong oauth not found prepare to create"), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, os.Ctx.Value(config.TraceKey)),
			log.O("package", pack), log.O("consumer_id", consumerID), log.O("body", helper.Stringify(auth)))

		oauth, status, err := os.Oauther.CreateOauth(consumerID, auth)
		if err != nil {
			log.Error(log.Msg("Failed create kong oauth", err.Error()), log.O("version", config.Version),
				log.O("project", config.ProjectName), log.O(config.TraceKey, os.Ctx.Value(config.TraceKey)),
				log.O("package", pack), log.O("body", helper.Stringify(oauth)),
				log.O("status_code", status), log.O("consumer_id", consumerID))
			return nil, status, err
		}

		log.Info(log.Msg("Success create kong oauth", helper.Stringify(oauth)),
			log.O("version", config.Version), log.O("project", config.ProjectName),
			log.O(config.TraceKey, os.Ctx.Value(config.TraceKey)), log.O("package", pack),
			log.O("consumer_id", consumerID), log.O("body", helper.Stringify(auth)),
			log.O("status_code", status))

		accessToken, status, err := os.Oauther.GetAccessToken(oauth.ClientID, oauth.ClientSecret, config.ActorScopes, userID)
		if err != nil {
			log.Error(log.Msg("Failed get access token", err.Error()), log.O("version", config.Version),
				log.O("project", config.ProjectName), log.O(config.TraceKey, os.Ctx.Value(config.TraceKey)),
				log.O("package", pack), log.O("body", helper.Stringify(oauth)), log.O("user_id", userID),
				log.O("status_code", status), log.O("consumer_id", consumerID),
				log.O("oauth", helper.Stringify(oauth)))
			return nil, status, err
		}

		return accessToken, status, nil
	}

	accessToken, status, err := os.Oauther.GetAccessToken(res[0].ClientID, res[0].ClientSecret, config.ActorScopes, userID)
	if err != nil {
		log.Error(log.Msg("Failed get access token", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, os.Ctx.Value(config.TraceKey)),
			log.O("package", pack), log.O("user_id", userID), log.O("status_code", status),
			log.O("consumer_id", consumerID), log.O("oauth", helper.Stringify(res[0])))
		return nil, status, err
	}

	log.Info(log.Msg("Success get access token", helper.Stringify(res)), log.O("version", config.Version),
		log.O("project", config.ProjectName), log.O(config.TraceKey, os.Ctx.Value(config.TraceKey)),
		log.O("package", pack), log.O("consumer_id", consumerID), log.O("type", oauthName),
		log.O("status_code", status))

	return accessToken, status, nil
}
