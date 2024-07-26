// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

package operations

import (
	"net/http"
)

type DeleteTransferInitiationRequest struct {
	// The transfer ID.
	TransferID string `pathParam:"style=simple,explode=false,name=transferId"`
}

func (o *DeleteTransferInitiationRequest) GetTransferID() string {
	if o == nil {
		return ""
	}
	return o.TransferID
}

type DeleteTransferInitiationResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
}

func (o *DeleteTransferInitiationResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *DeleteTransferInitiationResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *DeleteTransferInitiationResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}
