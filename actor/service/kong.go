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
			_, status, err := srv.Oauther.CreateOauth(consumerID, oauth)
			if err != nil {
				log.Error(log.Msg("Failed create kong oauth", err.Error()), log.O("version", config.Version),
					log.O("project", config.ProjectName), log.O(config.TraceKey, ks.Ctx.Value(config.TraceKey)),
					log.O("package", naming), log.O("body", helper.Stringify(oauth)),
					log.O("status_code", status))
				return
			}
			return
		}(ks, pack, cons.ID, auth)
	}

	return cons, status, nil
}
