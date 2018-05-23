package actor

import "sirius/actor/domain"

type (
	Address struct {
		ID         string       `json:"id,omitempty" validate:"omitempty,uuid4" sql:"type:uuid" gorm:"column:id;primary_key"`
		UserID     string       `json:"user_id" validate:"omitempty,uuid4" sql:"type:uuid" gorm:"column:user_id;unique;not null"`
		Street     string       `json:"street" validate:"required" gorm:"column:street;unique;not null"`
		City       string       `json:"city" validate:"required" gorm:"column:city;unique;not null"`
		State      string       `json:"state" validate:"required" gorm:"column:state;unique;not null"`
		PostalCode string       `json:"postal_code" validate:"required" gorm:"column:postal_code;unique;not null"`
		Country    domain.JSONB `json:"country" validate:"required" sql:"type:jsonb" gorm:"column:country;unique;not null"`
	}
)

func (Address) TableName() string {
	return "address"
}
