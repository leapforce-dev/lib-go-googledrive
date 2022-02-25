package googledrive

import (
	"net/http"
	"net/url"

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
		PermissionId string `json:"permissionId"`
		EmailAddress string `json:"emailAddress"`
	} `json:"user"`
}

func (service *Service) GetAbout(fields string) (*About, *errortools.Error) {
	values := url.Values{}
	values.Set("fields", fields)

	about := About{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		Url:           service.url("about"),
		Parameters:    &values,
		ResponseModel: &about,
	}
	_, _, e := service.googleService().HttpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &about, nil
}
