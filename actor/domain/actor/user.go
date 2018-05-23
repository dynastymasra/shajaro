package actor

import (
	"sirius/actor/domain"
	"sirius/actor/domain/kong"
	"time"
)

const (
	GenderMale   Gender = "Male"
	GenderFemale        = "Female"
)

type (
	Gender string

	User struct {
		ID         string       `json:"id,omitempty" validate:"omitempty,uuid4" sql:"type:uuid" gorm:"column:id;primary_key"`
		FirstName  string       `json:"first_name" validate:"required" gorm:"not null;column:first_name"`
		MiddleName string       `json:"middle_name,omitempty" validate:"omitempty" gorm:"column:middle_name"`
		LastName   string       `json:"last_name" validate:"required" gorm:"not null;column:last_name"`
		Phone      domain.JSONB `json:"phone" validate:"required" sql:"type:jsonb" gorm:"not null;column:phone"`
		Address    Address      `json:"address" validate:"required"  gorm:"foreignkey:UserID"`
		Gender     Gender       `json:"gender" validate:"required" gorm:"not null;column:gender"`
		Email      string       `json:"email" validate:"email,required" gorm:"not null;column:email"`
		BirthDate  time.Time    `json:"birth_date" validate:"required" gorm:"not null;column:birth_date"`
		Oauth      []kong.Oauth `json:"oauth,omitempty" validate:"omitempty" gorm:"foreignkey:OrderID"`
	}
)

func (User) TableName() string {
	return "users"
}
