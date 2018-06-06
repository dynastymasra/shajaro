package sql

import (
	"context"

	"shajaro/actor/domain/actor"

	"reflect"
	"runtime"
	"shajaro/actor/config"

	log "github.com/dynastymasra/gochill"
	"github.com/jinzhu/gorm"
)

type ActorRepository struct {
	DB  *gorm.DB
	Ctx context.Context
}

func NewUserRepository(ctx context.Context, db *gorm.DB) ActorRepository {
	return ActorRepository{
		Ctx: ctx,
		DB:  db,
	}
}

func (ar ActorRepository) CreateUser(user actor.User) error {
	return ar.DB.Create(&user).Error
}

func (ar ActorRepository) UserLogin(email, password string) (*actor.User, error) {
	user := &actor.User{}

	if err := ar.DB.Where("email = ? AND password = ?", email, password).First(user).Error; err != nil {
		log.Error(log.Msg("Failed login user", err.Error()), log.O("version", config.Version),
			log.O("project", config.ProjectName), log.O(config.TraceKey, ar.Ctx.Value(config.TraceKey)),
			log.O("package", runtime.FuncForPC(reflect.ValueOf(ar.UserLogin).Pointer()).Name()),
			log.O("email", email), log.O("password", password))
		return nil, err
	}

	return user, nil
}

func (ar ActorRepository) CheckEmailNotExist(email string) bool {
	var user actor.User
	return ar.DB.Where("email = ?", email).First(&user).RecordNotFound()
}
