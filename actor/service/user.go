package service

import (
	"context"
	"reflect"
	"runtime"
	"shajaro/actor/config"
	"shajaro/actor/domain/actor"
	"shajaro/actor/helper"

	log "github.com/dynastymasra/gochill"
)

type UserService struct {
	Ctx    context.Context
	Userer actor.Userer
}

func NewUserService(ctx context.Context, userer actor.Userer) UserService {
	return UserService{
		Ctx:    ctx,
		Userer: userer,
	}
}

func (us *UserService) UpdateService(data actor.User) (*actor.User, error) {
	pack := runtime.FuncForPC(reflect.ValueOf(us.UpdateService).Pointer()).Name()

	_, err := us.Userer.GetUserByID(data.ID)
	if err != nil {
		log.Error(log.Msg("Failed get user by id", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, us.Ctx.Value(config.TraceKey)),
			log.O("package", pack), log.O("data", helper.Stringify(data)))
		return nil, err
	}

	result, err := us.Userer.Update(data)
	if err != nil {
		log.Error(log.Msg("Failed update user", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, us.Ctx.Value(config.TraceKey)),
			log.O("package", pack), log.O("data", helper.Stringify(data)))
		return nil, err
	}

	return result, nil
}

func (us *UserService) DeleteService(id string) error {
	pack := runtime.FuncForPC(reflect.ValueOf(us.DeleteService).Pointer()).Name()

	user, err := us.Userer.GetUserByID(id)
	if err != nil {
		log.Error(log.Msg("Failed get user by id", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, us.Ctx.Value(config.TraceKey)),
			log.O("package", pack), log.O("id", id))
		return err
	}

	if err := us.Userer.Delete(*user); err != nil {
		log.Error(log.Msg("Failed delete user", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, us.Ctx.Value(config.TraceKey)),
			log.O("package", pack), log.O("user", helper.Stringify(user)))
		return err
	}
	return nil
}
