package sql

import (
	"context"

	"shajaro/actor/domain/actor"

	"reflect"
	"runtime"
	"shajaro/actor/config"

	"shajaro/actor/helper"

	log "github.com/dynastymasra/gochill"
	"github.com/jinzhu/gorm"
)

type UserRepository struct {
	DB  *gorm.DB
	Ctx context.Context
}

func NewUserRepository(ctx context.Context, db *gorm.DB) UserRepository {
	return UserRepository{
		Ctx: ctx,
		DB:  db,
	}
}

func (ar UserRepository) Create(user actor.Actor) error {
	return ar.DB.Create(&user).Error
}

func (ar UserRepository) Login(email string) (*actor.Actor, error) {
	user := &actor.Actor{}

	if err := ar.DB.Where("email = ?", email).First(user).Error; err != nil {
		log.Error(log.Msg("Failed login user", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, ar.Ctx.Value(config.TraceKey)),
			log.O("package", runtime.FuncForPC(reflect.ValueOf(ar.Login).Pointer()).Name()),
			log.O("email", email))
		return nil, err
	}
	return user, nil
}

func (ar UserRepository) CheckEmailNotExist(email string) bool {
	var user actor.Actor
	return ar.DB.Where("email = ?", email).First(&user).RecordNotFound()
}

func (ar UserRepository) GetUserByID(id string) (*actor.Actor, error) {
	user := &actor.Actor{}

	if err := ar.DB.Where("id = ?", id).First(user).Error; err != nil {
		log.Error(log.Msg("Failed get user by id", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, ar.Ctx.Value(config.TraceKey)),
			log.O("package", runtime.FuncForPC(reflect.ValueOf(ar.GetUserByID).Pointer()).Name()),
			log.O("id", id))
		return nil, err
	}
	return user, nil
}

func (ar UserRepository) Update(user actor.Actor) (*actor.Actor, error) {
	if err := ar.DB.Save(&user).Error; err != nil {
		log.Error(log.Msg("Failed update user", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, ar.Ctx.Value(config.TraceKey)),
			log.O("package", runtime.FuncForPC(reflect.ValueOf(ar.Update).Pointer()).Name()),
			log.O("data", helper.Stringify(user)))
		return nil, err
	}
	return &user, nil
}

func (ar UserRepository) UpdatePassword(user actor.Actor, password string) error {
	if err := ar.DB.Model(&user).Update("password", password).Error; err != nil {
		log.Error(log.Msg("Failed update password", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, ar.Ctx.Value(config.TraceKey)),
			log.O("package", runtime.FuncForPC(reflect.ValueOf(ar.UpdatePassword).Pointer()).Name()),
			log.O("body", helper.Stringify(user)), log.O("password", password))
		return err
	}
	return nil
}

func (ar UserRepository) Delete(user actor.Actor) error {
	if err := ar.DB.Delete(&user).Error; err != nil {
		log.Error(log.Msg("Failed delete user", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, ar.Ctx.Value(config.TraceKey)),
			log.O("package", runtime.FuncForPC(reflect.ValueOf(ar.Delete).Pointer()).Name()),
			log.O("user", helper.Stringify(user)))
		return err
	}
	return nil
}
