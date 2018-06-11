package service

import (
	"context"
	"reflect"
	"runtime"

	"shajaro/actor/config"
	"shajaro/actor/domain"
	"shajaro/actor/domain/kong"
	"shajaro/actor/helper"

	log "github.com/dynastymasra/gochill"
)

type KongService struct {
	Ctx     context.Context
	Konger  kong.Konger
	Oauther kong.Oauther
}

func NewKongService(ctx context.Context, konger kong.Konger, oauther kong.Oauther) KongService {
	return KongService{
		Ctx:     ctx,
		Konger:  konger,
		Oauther: oauther,
	}
}

func (ks *KongService) RegisterNewConsumer(consumer kong.Consumer) (*kong.Consumer, int, error) {
	pack := runtime.FuncForPC(reflect.ValueOf(ks.RegisterNewConsumer).Pointer()).Name()

	cons, status, err := ks.Konger.CreateConsumer(consumer)
	if err != nil {
		log.Error(log.Msg("Failed create kong consumer", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, ks.Ctx.Value(config.TraceKey)),
			log.O("package", pack), log.O("body", helper.Stringify(consumer)))
		return nil, status, err
	}

	names := []domain.OauthName{domain.OauthNameAndroid, domain.OauthNameWeb, domain.OauthNameIOS, domain.OauthNameDesktop}
	for _, name := range names {
		auth := kong.Oauth{
			Name:        name,
			RedirectURI: []string{config.RedirectURI},
		}

		go func(srv *KongService, naming, consumerID string, oauth kong.Oauth) {
			log.Info(log.Msg("Prepare create kong oauth"), log.O("version", config.Version),
				log.O("project", config.ProjectName), log.O(config.TraceKey, ks.Ctx.Value(config.TraceKey)),
				log.O("package", naming), log.O("consumer_id", consumerID),
				log.O("body", helper.Stringify(oauth)))

			res, status, err := srv.Oauther.CreateOauth(consumerID, oauth)
			if err != nil {
				log.Error(log.Msg("Failed create kong oauth", err.Error()), log.O("version", config.Version),
					log.O("project", config.ProjectName), log.O(config.TraceKey, ks.Ctx.Value(config.TraceKey)),
					log.O("package", naming), log.O("body", helper.Stringify(oauth)),
					log.O("status_code", status))
				return
			}

			log.Info(log.Msg("Success create kong oauth", helper.Stringify(res)), log.O("version", config.Version),
				log.O("project", config.ProjectName), log.O(config.TraceKey, ks.Ctx.Value(config.TraceKey)),
				log.O("package", naming), log.O("consumer_id", consumerID), log.O("status_code", status),
				log.O("body", helper.Stringify(oauth)))

			return
		}(ks, pack, cons.ID, auth)
	}

	log.Info(log.Msg("Success create kong consumer", helper.Stringify(cons)), log.O("version", config.Version),
		log.O("project", config.ProjectName), log.O(config.TraceKey, ks.Ctx.Value(config.TraceKey)),
		log.O("package", pack), log.O("body", helper.Stringify(consumer)), log.O("status_code", status))

	return cons, status, nil
}

func (ks *KongService) LoginService(consumerID string, oauthName domain.OauthName) (*kong.Oauth, int, error) {
	pack := runtime.FuncForPC(reflect.ValueOf(ks.LoginService).Pointer()).Name()

	res, status, err := ks.Oauther.GetOauthByName(consumerID, oauthName)
	if err != nil {
		log.Error(log.Msg("Failed get kong oauth", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, ks.Ctx.Value(config.TraceKey)),
			log.O("package", pack), log.O("consumer_id", consumerID), log.O("oauth_name", oauthName))
		return nil, status, err
	}

	if len(res) < 1 {
		auth := kong.Oauth{
			Name:        oauthName,
			RedirectURI: []string{config.RedirectURI},
		}

		log.Info(log.Msg("Kong oauth not found prepare to create"), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, ks.Ctx.Value(config.TraceKey)),
			log.O("package", pack), log.O("consumer_id", consumerID), log.O("body", helper.Stringify(auth)))

		oauth, status, err := ks.Oauther.CreateOauth(consumerID, auth)
		if err != nil {
			log.Error(log.Msg("Failed create kong oauth", err.Error()), log.O("version", config.Version),
				log.O("project", config.ProjectName), log.O(config.TraceKey, ks.Ctx.Value(config.TraceKey)),
				log.O("package", pack), log.O("body", helper.Stringify(oauth)),
				log.O("status_code", status))
			return nil, status, err
		}

		log.Info(log.Msg("Success create kong oauth", helper.Stringify(oauth)),
			log.O("version", config.Version), log.O("project", config.ProjectName),
			log.O(config.TraceKey, ks.Ctx.Value(config.TraceKey)), log.O("package", pack),
			log.O("consumer_id", consumerID), log.O("body", helper.Stringify(auth)),
			log.O("status_code", status))

		return oauth, status, nil
	}

	log.Info(log.Msg("Success get kong oauth", helper.Stringify(res)), log.O("version", config.Version),
		log.O("project", config.ProjectName), log.O(config.TraceKey, ks.Ctx.Value(config.TraceKey)),
		log.O("package", pack), log.O("consumer_id", consumerID), log.O("type", oauthName),
		log.O("status_code", status))

	return &res[0], status, nil
}
