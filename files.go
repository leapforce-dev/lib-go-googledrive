package GoogleDrive

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
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

func (gd *GoogleDrive) GetFiles(driveID *string, mimeType *string) (*[]File, error) {
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

	_, err := gd.Get(url, &filesReponse)
	if err != nil {
		return nil, err
	}

	return &filesReponse.Files, nil
}

func (gd *GoogleDrive) DownloadFile(fileID string) (*http.Response, error) {
	url := fmt.Sprintf("%s/files/%s?alt=media", apiURL, fileID)
	//fmt.Println(url)

	res, err := gd.Get(url, nil)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (gd *GoogleDrive) MoveFile(fileID string, fromDriveID string, toDriveID string) (*http.Response, error) {
	url := fmt.Sprintf("%s/files/%s?uploadType=media&addParents=%s&removeParents=%s", apiURL, fileID, toDriveID, fromDriveID)
	//fmt.Println(url)

	res, err := gd.Patch(url, nil)
	if err != nil {
		return nil, err
	}

	return res, nil
}
