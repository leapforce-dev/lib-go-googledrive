package GoogleDrive

import (
	"net/http"

	bigquerytools "github.com/Leapforce-nl/go_bigquerytools"

	oauth2 "github.com/Leapforce-nl/go_oauth2"
)

const (
	apiName         string = "GoogleDrive"
	apiURL          string = "https://www.googleapis.com/drive/v3"
	authURL         string = "https://accounts.google.com/o/oauth2/v2/auth"
	tokenURL        string = "https://oauth2.googleapis.com/token"
	tokenHTTPMethod string = http.MethodPost
	redirectURL     string = "http://localhost:8080/oauth/redirect"
)

// GoogleDrive stores GoogleDrive configuration
//
type GoogleDrive struct {
	oAuth2 *oauth2.OAuth2
}

// methods
//
func NewGoogleDrive(clientID string, clientSecret string, scope string, bigQuery *bigquerytools.BigQuery, isLive bool) (*GoogleDrive, error) {
	gd := GoogleDrive{}
	gd.oAuth2 = oauth2.NewOAuth(apiName, clientID, clientSecret, scope, redirectURL, authURL, tokenURL, tokenHTTPMethod, bigQuery, isLive)
	return &gd, nil
}

func (gd *GoogleDrive) ValidateToken() error {
	return gd.oAuth2.ValidateToken()
}

func (gd *GoogleDrive) InitToken() error {
	return gd.oAuth2.InitToken()
}

func (gd *GoogleDrive) Get(url string, model interface{}) (*http.Response, error) {
	res, err := gd.oAuth2.Get(url, model)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (gd *GoogleDrive) Patch(url string, model interface{}) (*http.Response, error) {
	res, err := gd.oAuth2.Patch(url, nil, model)

	if err != nil {
		return nil, err
	}

	return res, nil
}
