package sql

import (
	"context"

	"shajaro/actor/domain/actor"

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

func (ar ActorRepository) CheckEmailExist(email string) bool {
	var user actor.User
	return ar.DB.Where("email = ?", email).First(&user).RecordNotFound()
}
