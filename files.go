package googledrive

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type FilesResponse struct {
	Kind             string `json:"kind"`
	IncompleteSearch bool   `json:"incompleteSearch"`
	Files            []File `json:"files"`
}

type File struct {
	Kind              string   `json:"kind"`
	Id                string   `json:"id"`
	Name              string   `json:"name"`
	MimeType          string   `json:"mimeType"`
	Starred           bool     `json:"starred"`
	Trashed           bool     `json:"trashed"`
	ExplicitlyTrashed bool     `json:"explicitlyTrashed"`
	Parents           []string `json:"parents"`
}

func (service *Service) GetFiles(driveId *string, mimeType *string) (*[]File, *errortools.Error) {
	q := ""
	filters := []string{}

	if driveId != nil {
		filters = append(filters, fmt.Sprintf("'%s' in parents", *driveId))
	}

	if mimeType != nil {
		filters = append(filters, fmt.Sprintf("mimeType = '%s'", *mimeType))
	}

	if len(filters) > 0 {
		q = strings.Join(filters, " and ")
	}

	values := url.Values{}
	values.Set("q", q)

	filesReponse := FilesResponse{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		Url:           service.url("files"),
		Parameters:    &values,
		ResponseModel: &filesReponse,
	}
	_, _, e := service.googleService().HttpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &filesReponse.Files, nil
}

func (service *Service) GetFile(fileId string) (*File, *errortools.Error) {
	file := File{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		Url:           service.url(fmt.Sprintf("files/%s", fileId)),
		ResponseModel: &file,
	}
	_, _, e := service.googleService().HttpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &file, nil
}

func (service *Service) DownloadFile(fileId string) (*http.Response, *errortools.Error) {
	values := url.Values{}
	values.Set("alt", "media")

	requestConfig := go_http.RequestConfig{
		Method:     http.MethodGet,
		Url:        service.url(fmt.Sprintf("files/%s", fileId)),
		Parameters: &values,
	}
	_, res, e := service.googleService().HttpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return res, nil
}

func (service *Service) MoveFile(fileId string, fromDriveId string, toDriveId string) (*http.Response, *errortools.Error) {
	values := url.Values{}
	values.Set("uploadType", "media")
	values.Set("addParents", toDriveId)
	values.Set("removeParents", fromDriveId)

	requestConfig := go_http.RequestConfig{
		Method:     http.MethodPatch,
		Url:        service.url(fmt.Sprintf("files/%s", fileId)),
		Parameters: &values,
	}
	_, res, e := service.googleService().HttpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return res, nil
}

func (service *Service) ExportFile(fileId string, mimeType string) (*http.Response, *errortools.Error) {
	values := url.Values{}
	values.Set("mimeType", mimeType)

	requestConfig := go_http.RequestConfig{
		Method:     http.MethodGet,
		Url:        service.url(fmt.Sprintf("files/%s/export", fileId)),
		Parameters: &values,
	}

	_, res, e := service.googleService().HttpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return res, nil
}

func (service *Service) CreateFile(parentId string, name string, mimeType string) (*File, *errortools.Error) {
	data := struct {
		MimeType string   `json:"mimeType"`
		Name     string   `json:"name"`
		Parents  []string `json:"parents"`
	}{
		mimeType,
		name,
		[]string{parentId},
	}

	file := File{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPost,
		Url:           service.url("files"),
		BodyModel:     data,
		ResponseModel: &file,
	}

	_, _, e := service.googleService().HttpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &file, nil
}

func (service *Service) UpdateFile(fileId string, mimeType string, content *[]byte) (*File, *errortools.Error) {
	file := File{}

	header := http.Header{}
	header.Set("Content-Type", mimeType)

	requestConfig := go_http.RequestConfig{
		Method:            http.MethodPatch,
		Url:               fmt.Sprintf("https://www.googleapis.com/upload/drive/v3/files/%s", fileId),
		BodyRaw:           content,
		ResponseModel:     &file,
		NonDefaultHeaders: &header,
	}

	_, _, e := service.googleService().HttpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &file, nil
}
