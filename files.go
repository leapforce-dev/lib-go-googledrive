package GoogleDrive

import (
	"fmt"
	"net/http"
	"net/url"
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

func (gd *GoogleDrive) GetFiles(q string) (*[]File, error) {
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
	fmt.Println(url)

	res, err := gd.Get(url, nil)
	if err != nil {
		return nil, err
	}

	return res, nil
}
