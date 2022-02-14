package googledrive

import (
	"fmt"

	errortools "github.com/leapforce-libraries/go_errortools"
	google "github.com/leapforce-libraries/go_google"
	"github.com/leapforce-libraries/go_oauth2/tokensource"
)

const (
	apiName string = "GoogleDrive"
	apiURL  string = "https://www.googleapis.com/drive/v3"
)

type Service struct {
	clientID      string
	googleService *google.Service
}

type ServiceConfig struct {
	ClientID     string
	ClientSecret string
	TokenSource  tokensource.TokenSource
}

func NewService(serviceConfig *ServiceConfig) (*Service, *errortools.Error) {
	if serviceConfig == nil {
		return nil, errortools.ErrorMessage("ServiceConfig must not be a nil pointer")
	}

	if serviceConfig.ClientID == "" {
		return nil, errortools.ErrorMessage("ClientID not provided")
	}

	if serviceConfig.ClientSecret == "" {
		return nil, errortools.ErrorMessage("ClientSecret not provided")
	}

	googleServiceConfig := google.ServiceConfig{
		APIName:      apiName,
		ClientID:     serviceConfig.ClientID,
		ClientSecret: serviceConfig.ClientSecret,
		TokenSource:  serviceConfig.TokenSource,
	}

	googleService, e := google.NewService(&googleServiceConfig)
	if e != nil {
		return nil, e
	}

	return &Service{
		clientID:      serviceConfig.ClientID,
		googleService: googleService,
	}, nil
}

func (service *Service) url(path string) string {
	return fmt.Sprintf("%s/%s", apiURL, path)
}

func (service *Service) InitToken(scope string, accessType *string, prompt *string, state *string) *errortools.Error {
	return service.googleService.InitToken(scope, accessType, prompt, state)
}

func (service *Service) APIName() string {
	return apiName
}

func (service *Service) APIKey() string {
	return service.clientID
}

func (service *Service) APICallCount() int64 {
	return service.googleService.APICallCount()
}

func (service *Service) APIReset() {
	service.googleService.APIReset()
}
