package googledrive

import (
	"fmt"

	google "github.com/leapforce-libraries/go_google"
)

const (
	apiName string = "GoogleDrive"
	apiURL  string = "https://www.googleapis.com/drive/v3"
)

type Service google.Service

func (service *Service) url(path string) string {
	return fmt.Sprintf("%s/%s", apiURL, path)
}

func (service *Service) APIName() string {
	return apiName
}

func (service *Service) APIKey() string {
	return service.googleService().APIKey()
}

func (service *Service) APICallCount() int64 {
	return service.googleService().APICallCount()
}

func (service *Service) APIReset() {
	service.googleService().APIReset()
}

func (service *Service) googleService() *google.Service {
	googleService := google.Service(*service)
	return &googleService
}
