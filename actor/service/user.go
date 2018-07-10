package service

import (
	"context"
	"reflect"
	"runtime"

	"github.com/dynastymasra/shajaro/actor/config"
	"github.com/dynastymasra/shajaro/actor/domain/actor"
	"github.com/dynastymasra/shajaro/actor/helper"

	log "github.com/dynastymasra/gochill"
)

type UserService struct {
	Ctx  context.Context
	User actor.User
}

func NewUserService(ctx context.Context, user actor.User) UserService {
	return UserService{
		Ctx:  ctx,
		User: user,
	}
}

func (us *UserService) UpdateService(data actor.Actor) (*actor.Actor, error) {
	pack := runtime.FuncForPC(reflect.ValueOf(us.UpdateService).Pointer()).Name()

	user, err := us.User.GetUserByID(data.ID)
	if err != nil {
		log.Error(log.Msg("Failed get user by id", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, us.Ctx.Value(config.TraceKey)),
			log.O("package", pack), log.O("data", helper.Stringify(data)))
		return nil, err
	}

	data.Password = user.Password
	result, err := us.User.Update(data)
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

	user, err := us.User.GetUserByID(id)
	if err != nil {
		log.Error(log.Msg("Failed get user by id", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, us.Ctx.Value(config.TraceKey)),
			log.O("package", pack), log.O("id", id))
		return err
	}

	if err := us.User.Delete(*user); err != nil {
		log.Error(log.Msg("Failed delete user", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, us.Ctx.Value(config.TraceKey)),
			log.O("package", pack), log.O("user", helper.Stringify(user)))
		return err
	}
	return nil
}
