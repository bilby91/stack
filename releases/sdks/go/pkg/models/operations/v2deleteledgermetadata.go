// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

package operations

import (
	"net/http"
)

type V2DeleteLedgerMetadataRequest struct {
	// Key to remove.
	Key string `pathParam:"style=simple,explode=false,name=key"`
	// Name of the ledger.
	Ledger string `pathParam:"style=simple,explode=false,name=ledger"`
}

func (o *V2DeleteLedgerMetadataRequest) GetKey() string {
	if o == nil {
		return ""
	}
	return o.Key
}

func (o *V2DeleteLedgerMetadataRequest) GetLedger() string {
	if o == nil {
		return ""
	}
	return o.Ledger
}

type V2DeleteLedgerMetadataResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
}

func (o *V2DeleteLedgerMetadataResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *V2DeleteLedgerMetadataResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *V2DeleteLedgerMetadataResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}
