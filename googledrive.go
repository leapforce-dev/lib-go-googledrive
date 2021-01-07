package GoogleDrive

import (
	google "github.com/leapforce-libraries/go_google"
)

const (
	apiName string = "GoogleDrive"
	apiURL  string = "https://www.googleapis.com/drive/v3"
)

// GoogleDrive stores GoogleDrive configuration
//
type GoogleDrive struct {
	Client *google.GoogleClient
}

// methods
//
func NewGoogleDrive(clientID string, clientSecret string, scope string, bigQuery *google.BigQuery) *GoogleDrive {
	config := google.GoogleClientConfig{
		APIName:      apiName,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scope:        scope,
	}

	googleClient := google.NewGoogleClient(config, bigQuery)

	return &GoogleDrive{googleClient}
}
