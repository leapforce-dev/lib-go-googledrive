package googledrive

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	errortools "github.com/leapforce-libraries/go_errortools"
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

	url := fmt.Sprintf("%s/files?q=%s", APIURL, url.QueryEscape(q))
	//fmt.Println(url)

	filesReponse := FilesResponse{}

	_, _, e := service.googleService.Get(url, &filesReponse)
	if e != nil {
		return nil, e
	}

	return &filesReponse.Files, nil
}

func (service *Service) GetFile(fileID string) (*File, *errortools.Error) {
	url := fmt.Sprintf("%s/files/%s", APIURL, fileID)
	//fmt.Println(url)

	file := File{}

	_, _, e := service.googleService.Get(url, &file)
	if e != nil {
		return nil, e
	}

	return &file, nil
}

func (service *Service) DownloadFile(fileID string) (*http.Response, *errortools.Error) {
	url := fmt.Sprintf("%s/files/%s?alt=media", APIURL, fileID)
	//fmt.Println(url)

	_, res, e := service.googleService.Get(url, nil)
	if e != nil {
		return nil, e
	}

	return res, nil
}

func (service *Service) MoveFile(fileID string, fromDriveID string, toDriveID string) (*http.Response, *errortools.Error) {
	url := fmt.Sprintf("%s/files/%s?uploadType=media&addParents=%s&removeParents=%s", APIURL, fileID, toDriveID, fromDriveID)
	//fmt.Println(url)

	_, res, e := service.googleService.Patch(url, nil, nil)
	if e != nil {
		return nil, e
	}

	return res, nil
}

func (service *Service) ExportFile(fileID string, mimeType string) (*http.Response, *errortools.Error) {
	url := fmt.Sprintf("%s/files/%s/export?mimeType=%s", APIURL, fileID, mimeType)
	//fmt.Println(url)

	_, res, e := service.googleService.Get(url, nil)
	if e != nil {
		return nil, e
	}

	return res, nil
}
