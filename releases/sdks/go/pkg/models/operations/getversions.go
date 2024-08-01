// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

package operations

import (
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"net/http"
)

type GetVersionsResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// OK
	GetVersionsResponse *shared.GetVersionsResponse
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
}

func (o *GetVersionsResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *GetVersionsResponse) GetGetVersionsResponse() *shared.GetVersionsResponse {
	if o == nil {
		return nil
	}
	return o.GetVersionsResponse
}

func (o *GetVersionsResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *GetVersionsResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}
