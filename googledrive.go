package GoogleDrive

import (
	"net/http"

	bigquerytools "github.com/leapforce-libraries/go_bigquerytools"
	errortools "github.com/leapforce-libraries/go_errortools"
	oauth2 "github.com/leapforce-libraries/go_oauth2"
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
func NewGoogleDrive(clientID string, clientSecret string, scope string, bigQuery *bigquerytools.BigQuery) *GoogleDrive {
	gd := GoogleDrive{}
	config := oauth2.OAuth2Config{
		ApiName:         apiName,
		ClientID:        clientID,
		ClientSecret:    clientSecret,
		Scope:           scope,
		RedirectURL:     redirectURL,
		AuthURL:         authURL,
		TokenURL:        tokenURL,
		TokenHTTPMethod: tokenHTTPMethod,
	}
	gd.oAuth2 = oauth2.NewOAuth(config, bigQuery)
	return &gd
}

func (gd *GoogleDrive) ValidateToken() (*oauth2.Token, *errortools.Error) {
	return gd.oAuth2.ValidateToken()
}

func (gd *GoogleDrive) InitToken() *errortools.Error {
	return gd.oAuth2.InitToken()
}

func (gd *GoogleDrive) Get(url string, model interface{}) (*http.Response, *errortools.Error) {
	_, res, e := gd.oAuth2.Get(url, model, nil)

	if e != nil {
		return nil, e
	}

	return res, nil
}

func (gd *GoogleDrive) Patch(url string, model interface{}) (*http.Response, *errortools.Error) {
	_, res, e := gd.oAuth2.Patch(url, nil, model, nil)

	if e != nil {
		return nil, e
	}

	return res, nil
}
