package kong

import "github.com/dynastymasra/shajaro/actor/domain"

type (
	Oauth struct {
		ID           string           `json:"id,omitempty"`
		Name         domain.OauthName `json:"name"`
		ClientID     string           `json:"client_id,omitempty"`
		ClientSecret string           `json:"client_secret,omitempty"`
		RedirectURI  []string         `json:"redirect_uri"`
	}

	AccessToken struct {
		RefreshToken string `json:"refresh_token"`
		TokenType    string `json:"token_type"`
		AccessToken  string `json:"access_token"`
		ExpiresIn    int64  `json:"expires_in"`
	}

	Oauther interface {
		CreateOauth(string, Oauth) (*Oauth, int, error)
		DeleteOauth(string, string) (int, error)
		GetAccessToken(clientID, clientSecret, scope, userID string) (*AccessToken, int, error)
		GetAllOauth(string) ([]Oauth, int, error)
		GetOauthByName(string, domain.OauthName) ([]Oauth, int, error)
	}
)
