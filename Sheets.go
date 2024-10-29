package googledrive

import (
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"net/http"
)

type ValueRange struct {
	Range          string     `json:"range"`
	MajorDimension string     `json:"majorDimension"`
	Values         [][]string `json:"values"`
}

type ReadSheetConfig struct {
	SpreadsheetId        string
	Range                string
	MajorDimension       *string
	ValueRenderOption    *string
	DateTimeRenderOption *string
}

func (service *Service) ReadSheet(config *ReadSheetConfig) (*ValueRange, *errortools.Error) {
	if config == nil {
		return nil, errortools.ErrorMessage("config must not be nil")
	}

	//values := url.Values{}

	var valueRange ValueRange

	requestConfig := go_http.RequestConfig{
		Method: http.MethodGet,
		Url:    service.sheetsUrl(fmt.Sprintf("%s/values/%s", config.SpreadsheetId, config.Range)),
		//Parameters:    &values,
		ResponseModel: &valueRange,
	}
	_, _, e := service.googleService().HttpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &valueRange, nil
}
