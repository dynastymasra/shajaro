package kong

const (
	OauthTypeAndroid OauthType = "Android"
	OauthTypeIOS               = "iOS"
	OauthTypeWeb               = "Web"
	OauthTypeDesktop           = "Desktop"
)

type (
	OauthType string

	Oauth struct {
		ID           string    `json:"id,omitempty" validate:"omitempty"  sql:"type:uuid" gorm:"not null;column:id"`
		Name         OauthType `json:"name,omitempty" validate:"omitempty" gorm:"not null;column:name"`
		ClientID     string    `json:"client_id,omitempty" validate:"omitempty" gorm:"not null;column:client_id"`
		ClientSecret string    `json:"client_secret,omitempty" validate:"omitempty" gorm:"not null;column:client_secret"`
		ConsumerID   string    `json:"consumer_id,omitempty" validate:"omitempty" sql:"type:uuid" gorm:"not null;column:consumer_id"`
		RedirectURI  string    `json:"redirect_uri,omitempty" validate:"omitempty,url" gorm:"-"`
	}
)

func (Oauth) TableName() string {
	return "oauth"
}
