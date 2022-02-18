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
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	MimeType          string   `json:"mimeType"`
	Starred           bool     `json:"starred"`
	Trashed           bool     `json:"trashed"`
	ExplicitlyTrashed bool     `json:"explicitlyTrashed"`
	Parents           []string `json:"parents"`
}

func (service *Service) GetFiles(driveID *string, mimeType *string) (*[]File, *errortools.Error) {
	q := ""
	filters := []string{}

	if driveID != nil {
		filters = append(filters, fmt.Sprintf("'%s' in parents", *driveID))
	}

	if mimeType != nil {
		filters = append(filters, fmt.Sprintf("mimeType = '%s'", *mimeType))
	}

	if len(filters) > 0 {
		q = strings.Join(filters, " and ")
	}

	filesReponse := FilesResponse{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		URL:           service.url(fmt.Sprintf("files?q=%s", url.QueryEscape(q))),
		ResponseModel: &filesReponse,
	}
	_, _, e := service.googleService().HttpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &filesReponse.Files, nil
}

func (service *Service) GetFile(fileID string) (*File, *errortools.Error) {
	file := File{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		URL:           service.url(fmt.Sprintf("files/%s", fileID)),
		ResponseModel: &file,
	}
	_, _, e := service.googleService().HttpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &file, nil
}

func (service *Service) DownloadFile(fileID string) (*http.Response, *errortools.Error) {
	requestConfig := go_http.RequestConfig{
		Method: http.MethodGet,
		URL:    service.url(fmt.Sprintf("files/%s?alt=media", fileID)),
	}
	_, res, e := service.googleService().HttpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return res, nil
}

func (service *Service) MoveFile(fileID string, fromDriveID string, toDriveID string) (*http.Response, *errortools.Error) {
	requestConfig := go_http.RequestConfig{
		Method: http.MethodPatch,
		URL:    service.url(fmt.Sprintf("files/%s?uploadType=media&addParents=%s&removeParents=%s", fileID, toDriveID, fromDriveID)),
	}
	_, res, e := service.googleService().HttpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return res, nil
}

func (service *Service) ExportFile(fileID string, mimeType string) (*http.Response, *errortools.Error) {
	requestConfig := go_http.RequestConfig{
		Method: http.MethodGet,
		URL:    service.url(fmt.Sprintf("files/%s/export?mimeType=%s", fileID, mimeType)),
	}

	_, res, e := service.googleService().HttpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return res, nil
}

func (service *Service) CreateFile(parentID string, name string, mimeType string) (*File, *errortools.Error) {
	data := struct {
		MimeType string   `json:"mimeType"`
		Name     string   `json:"name"`
		Parents  []string `json:"parents"`
	}{
		mimeType,
		name,
		[]string{parentID},
	}

	file := File{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPost,
		URL:           service.url("files"),
		BodyModel:     data,
		ResponseModel: &file,
	}

	_, _, e := service.googleService().HttpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &file, nil
}

func (service *Service) UpdateFile(fileID string, mimeType string, content *[]byte) (*File, *errortools.Error) {
	file := File{}

	header := http.Header{}
	header.Set("Content-Type", mimeType)

	requestConfig := go_http.RequestConfig{
		Method:            http.MethodPatch,
		URL:               fmt.Sprintf("https://www.googleapis.com/upload/drive/v3/files/%s", fileID),
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
