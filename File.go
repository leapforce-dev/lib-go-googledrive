package googledrive

import (
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"net/http"
	"net/url"
	"strings"
)

type FilesResponse struct {
	Kind             string `json:"kind"`
	NextPageToken    string `json:"nextPageToken"`
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
	Owners            []User   `json:"owners"`
	WebViewLink       string   `json:"webViewLink"`
}

type GetFilesConfig struct {
	DriveId                   *string
	Fields                    *string
	MimeType                  *string
	Trashed                   *bool
	IncludeItemsFromAllDrives *bool
	SupportsAllDrives         *bool
	SupportsTeamDrives        *bool
}

func (service *Service) GetFiles(config *GetFilesConfig) (*[]File, *errortools.Error) {
	values := url.Values{}

	filters := []string{}

	if config != nil {

		if config.DriveId != nil {
			filters = append(filters, fmt.Sprintf("\"%s\" in parents", *config.DriveId))
		}

		if config.Fields != nil {
			values.Set("fields", *config.Fields)
		}

		if config.MimeType != nil {
			filters = append(filters, fmt.Sprintf("mimeType = '%s'", *config.MimeType))
		}

		if config.Trashed != nil {
			filters = append(filters, fmt.Sprintf("trashed = %s", fmt.Sprintf("%v", *config.Trashed)))
		}

		if config.IncludeItemsFromAllDrives != nil {
			values.Set("includeItemsFromAllDrives", fmt.Sprintf("%v", *config.IncludeItemsFromAllDrives))
		}

		if config.SupportsAllDrives != nil {
			values.Set("supportsAllDrives", fmt.Sprintf("%v", *config.SupportsAllDrives))
		}

		if config.SupportsTeamDrives != nil {
			values.Set("supportsTeamDrives", fmt.Sprintf("%v", *config.SupportsTeamDrives))
		}
	}

	if len(filters) > 0 {
		values.Set("q", strings.Join(filters, " and "))
	}

	var files []File

	for {
		var filesResponse FilesResponse

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           service.url("files"),
			Parameters:    &values,
			ResponseModel: &filesResponse,
		}
		_, _, e := service.googleService().HttpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		files = append(files, filesResponse.Files...)

		if filesResponse.NextPageToken == "" {
			break
		}

		values.Set("pageToken", filesResponse.NextPageToken)
	}

	return &files, nil
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

type CreateFileConfig struct {
	SupportsAllDrives *bool
}

func (service *Service) CreateFile(parentId string, name string, mimeType string, config *CreateFileConfig) (*File, *errortools.Error) {
	data := struct {
		MimeType string   `json:"mimeType"`
		Name     string   `json:"name"`
		Parents  []string `json:"parents"`
	}{
		mimeType,
		name,
		[]string{parentId},
	}

	values := url.Values{}

	if config != nil {
		if config.SupportsAllDrives != nil {
			values.Set("supportsAllDrives", fmt.Sprintf("%v", *config.SupportsAllDrives))
		}
	}

	file := File{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPost,
		Url:           service.url("files"),
		Parameters:    &values,
		BodyModel:     data,
		ResponseModel: &file,
	}

	_, _, e := service.googleService().HttpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &file, nil
}

type UpdateFileConfig struct {
	SupportsAllDrives *bool
}

func (service *Service) UpdateFile(fileId string, mimeType string, content *[]byte, config *UpdateFileConfig) (*File, *errortools.Error) {
	file := File{}

	header := http.Header{}
	header.Set("Content-Type", mimeType)

	values := url.Values{}

	values.Set("uploadType", "media")

	if config != nil {
		if config.SupportsAllDrives != nil {
			values.Set("supportsAllDrives", fmt.Sprintf("%v", *config.SupportsAllDrives))
		}
	}

	requestConfig := go_http.RequestConfig{
		Method:            http.MethodPatch,
		Url:               fmt.Sprintf("https://www.googleapis.com/upload/drive/v3/files/%s", fileId),
		Parameters:        &values,
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
