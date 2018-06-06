package actor

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"reflect"
	"runtime"

	"shajaro/actor/config"
	"shajaro/actor/domain"
	"shajaro/actor/helper"

	log "github.com/dynastymasra/gochill"
)

type (
	User struct {
		ID         string        `json:"id,omitempty" validate:"omitempty,uuid4" sql:"type:uuid" gorm:"column:id;primary_key"`
		ConsumerID string        `json:"consumer_id,omitempty" validate:"omitempty,uuid4" sql:"type:uuid" gorm:"column:consumer_id;unique"`
		FirstName  string        `json:"first_name" validate:"required" gorm:"not null;column:first_name"`
		MiddleName string        `json:"middle_name,omitempty" validate:"omitempty" gorm:"column:middle_name"`
		LastName   string        `json:"last_name" validate:"required" gorm:"not null;column:last_name"`
		AvatarURL  string        `json:"avatar_url" validate:"required,url" gorm:"not null;column:avatar_url"`
		Phone      Phone         `json:"phone" validate:"required" sql:"type:jsonb" gorm:"not null;column:phone"`
		Address    Address       `json:"address" validate:"required"  sql:"type:jsonb" gorm:"not null;column:address"`
		Gender     domain.Gender `json:"gender" validate:"required" gorm:"not null;column:gender"`
		Email      string        `json:"email" validate:"email,required" gorm:"unique;not null;column:email"`
		Password   string        `json:"password" validate:"required" gorm:"not null;column:password"`
		BirthDate  string        `json:"birth_date" validate:"required" gorm:"not null;column:birth_date"`
	}

	Login struct {
		Email     string           `json:"email" validate:"email,required"`
		Password  string           `json:"password" validate:"required"`
		OauthName domain.OauthName `json:"type" validate:"required"`
	}

	Phone struct {
		CallingCode string `json:"calling_code" validate:"required"`
		Number      string `json:"number" validate:"required"`
	}

	Address struct {
		Street     string `json:"street" validate:"required"`
		City       string `json:"city" validate:"required"`
		State      string `json:"state" validate:"required"`
		PostalCode string `json:"postal_code" validate:"required"`
		Country    string `json:"country" validate:"required" `
	}

	Userer interface {
		CheckEmailNotExist(string) bool
		UserLogin(string, string) (*User, error)
		CreateUser(User) error
	}
)

func (User) TableName() string {
	return "users"
}

// Value implement value interface
func (u Phone) Value() (driver.Value, error) {
	value, err := json.Marshal(u)
	if err != nil {
		log.Error(log.Msg("Failed marshal phone", err.Error()), log.O("version", config.Version),
			log.O("package", runtime.FuncForPC(reflect.ValueOf(u.Value).Pointer()).Name()),
			log.O("project", config.ProjectName), log.O("data", helper.Stringify(u)))
		return nil, err
	}
	return value, nil
}

func (u Address) Value() (driver.Value, error) {
	value, err := json.Marshal(u)
	if err != nil {
		log.Error(log.Msg("Failed marshal address", err.Error()), log.O("version", config.Version),
			log.O("package", runtime.FuncForPC(reflect.ValueOf(u.Value).Pointer()).Name()),
			log.O("project", config.ProjectName), log.O("data", helper.Stringify(u)))
		return nil, err
	}
	return value, nil
}

func (u *Phone) Scan(value interface{}) error {
	source, ok := value.([]byte)
	if !ok {
		log.Error(log.Msg("Failed casting data phone", config.ErrCastingData), log.O("version", config.Version),
			log.O("package", runtime.FuncForPC(reflect.ValueOf(u.Scan).Pointer()).Name()),
			log.O("project", config.ProjectName), log.O("data", helper.Stringify(u)))
		return errors.New(config.ErrCastingData)
	}

	if err := json.Unmarshal(source, &u); err != nil {
		log.Error(log.Msg("Failed unmarshal phone", err.Error()), log.O("version", config.Version),
			log.O("package", runtime.FuncForPC(reflect.ValueOf(u.Scan).Pointer()).Name()),
			log.O("project", config.ProjectName), log.O("data", string(source)))
		return err
	}
	return nil
}

func (u *Address) Scan(value interface{}) error {
	source, ok := value.([]byte)
	if !ok {
		log.Error(log.Msg("Failed casting data address", config.ErrCastingData), log.O("version", config.Version),
			log.O("package", runtime.FuncForPC(reflect.ValueOf(u.Scan).Pointer()).Name()),
			log.O("project", config.ProjectName), log.O("data", helper.Stringify(u)))
		return errors.New(config.ErrCastingData)
	}

	if err := json.Unmarshal(source, &u); err != nil {
		log.Error(log.Msg("Failed unmarshal address", err.Error()), log.O("version", config.Version),
			log.O("package", runtime.FuncForPC(reflect.ValueOf(u.Scan).Pointer()).Name()),
			log.O("project", config.ProjectName), log.O("data", string(source)))
		return err
	}
	return nil
}
