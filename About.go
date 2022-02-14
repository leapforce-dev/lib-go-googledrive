package googledrive

import (
	"fmt"
	"net/http"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type About struct {
	Kind string `json:"kind"`
	User struct {
		Kind         string `json:"kind"`
		DisplayName  string `json:"displayName"`
		PhotoLink    string `json:"photoLink"`
		Me           bool   `json:"me"`
		PermissionID string `json:"permissionId"`
		EmailAddress string `json:"emailAddress"`
	} `json:"user"`
}

func (service *Service) GetAbout(fields string) (*About, *errortools.Error) {
	about := About{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		URL:           service.url(fmt.Sprintf("about?fields=%s", fields)),
		ResponseModel: &about,
	}
	_, _, e := service.googleService.HTTPRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &about, nil
}
