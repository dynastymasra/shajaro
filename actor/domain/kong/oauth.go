package kong

import "shajaro/actor/domain"

type (
	Oauth struct {
		ID           string           `json:"id,omitempty"`
		Name         domain.OauthName `json:"name"`
		ClientID     string           `json:"client_id,omitempty"`
		ClientSecret string           `json:"client_secret,omitempty"`
		RedirectURI  []string         `json:"redirect_uri"`
	}

	Oauther interface {
		CreateOauth(string, Oauth) (*Oauth, int, error)
		GetOauthByName(string, domain.OauthName) ([]Oauth, int, error)
	}
)
