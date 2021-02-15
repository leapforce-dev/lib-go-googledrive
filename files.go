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
		URL:           service.url(fmt.Sprintf("files?q=%s", url.QueryEscape(q))),
		ResponseModel: &filesReponse,
	}
	_, _, e := service.googleService.Get(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &filesReponse.Files, nil
}

func (service *Service) GetFile(fileID string) (*File, *errortools.Error) {
	file := File{}

	requestConfig := go_http.RequestConfig{
		URL:           service.url(fmt.Sprintf("files/%s", fileID)),
		ResponseModel: &file,
	}
	_, _, e := service.googleService.Get(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &file, nil
}

func (service *Service) DownloadFile(fileID string) (*http.Response, *errortools.Error) {
	requestConfig := go_http.RequestConfig{
		URL: service.url(fmt.Sprintf("files/%s?alt=media", fileID)),
	}
	_, res, e := service.googleService.Get(&requestConfig)
	if e != nil {
		return nil, e
	}

	return res, nil
}

func (service *Service) MoveFile(fileID string, fromDriveID string, toDriveID string) (*http.Response, *errortools.Error) {
	requestConfig := go_http.RequestConfig{
		URL: service.url(fmt.Sprintf("files/%s?uploadType=media&addParents=%s&removeParents=%s", fileID, toDriveID, fromDriveID)),
	}
	_, res, e := service.googleService.Patch(&requestConfig)
	if e != nil {
		return nil, e
	}

	return res, nil
}

func (service *Service) ExportFile(fileID string, mimeType string) (*http.Response, *errortools.Error) {
	requestConfig := go_http.RequestConfig{
		URL: service.url(fmt.Sprintf("files/%s/export?mimeType=%s", fileID, mimeType)),
	}

	_, res, e := service.googleService.Get(&requestConfig)
	if e != nil {
		return nil, e
	}

	return res, nil
}
