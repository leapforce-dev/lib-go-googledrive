package GoogleDrive

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

func (gd *GoogleDrive) GetFiles(driveID *string, mimeType *string) (*[]File, *errortools.Error) {
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

	url := fmt.Sprintf("%s/files?q=%s", apiURL, url.QueryEscape(q))
	//fmt.Println(url)

	filesReponse := FilesResponse{}

	_, _, e := gd.Client.Get(url, &filesReponse)
	if e != nil {
		return nil, e
	}

	return &filesReponse.Files, nil
}

func (gd *GoogleDrive) GetFile(fileID string) (*File, *errortools.Error) {
	url := fmt.Sprintf("%s/files/%s", apiURL, fileID)
	//fmt.Println(url)

	file := File{}

	_, _, e := gd.Client.Get(url, &file)
	if e != nil {
		return nil, e
	}

	return &file, nil
}

func (gd *GoogleDrive) DownloadFile(fileID string) (*http.Response, *errortools.Error) {
	url := fmt.Sprintf("%s/files/%s?alt=media", apiURL, fileID)
	//fmt.Println(url)

	_, res, e := gd.Client.Get(url, nil)
	if e != nil {
		return nil, e
	}

	return res, nil
}

func (gd *GoogleDrive) MoveFile(fileID string, fromDriveID string, toDriveID string) (*http.Response, *errortools.Error) {
	url := fmt.Sprintf("%s/files/%s?uploadType=media&addParents=%s&removeParents=%s", apiURL, fileID, toDriveID, fromDriveID)
	//fmt.Println(url)

	_, res, e := gd.Client.Patch(url, nil, nil)
	if e != nil {
		return nil, e
	}

	return res, nil
}

func (gd *GoogleDrive) ExportFile(fileID string, mimeType string) (*http.Response, *errortools.Error) {
	url := fmt.Sprintf("%s/files/%s/export?mimeType=%s", apiURL, fileID, mimeType)
	//fmt.Println(url)

	_, res, e := gd.Client.Get(url, nil)
	if e != nil {
		return nil, e
	}

	return res, nil
}
