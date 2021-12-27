package model

import (
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

// CommonHTTPResponse ...
type CommonHTTPResponse struct {
	ResultCode string `json:"resultCode"`
	ResultMsg  string `json:"resultMsg"`
	TrID       string `json:"trID"`
}

// CountHTTPResponse ...
type CountHTTPResponse struct {
	CommonHTTPResponse
	Count uint `json:"count"`
}

type TBackLogListResponse struct {
	CommonHTTPResponse
	Logs []*cloudwatchlogs.OutputLogEvent `json:"logs"`
}

// SetResult ...
func (c *CommonHTTPResponse) SetResult(resultCode, resultMsg, trID string) {
	c.ResultCode = resultCode
	c.ResultMsg = resultMsg
	c.TrID = trID
}

// CmcListResponse ...
type CmcListResponse struct {
	CommonHTTPResponse
	CmcList []*Cmc `json:"cmc"`
}

// CmcResponse ...
type CmcResponse struct {
	CommonHTTPResponse
	Cmc *Cmc `json:"cmc"`
}

// ConversionResponse ...
type ConversionResponse struct {
	CommonHTTPResponse
	Conversion *Conversion `json:"conversion"`
}
