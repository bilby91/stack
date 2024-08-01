// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

package operations

import (
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"net/http"
)

type GetWorkflowRequest struct {
	// The flow id
	FlowID string `pathParam:"style=simple,explode=false,name=flowId"`
}

func (o *GetWorkflowRequest) GetFlowID() string {
	if o == nil {
		return ""
	}
	return o.FlowID
}

type GetWorkflowResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// The workflow
	GetWorkflowResponse *shared.GetWorkflowResponse
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
}

func (o *GetWorkflowResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *GetWorkflowResponse) GetGetWorkflowResponse() *shared.GetWorkflowResponse {
	if o == nil {
		return nil
	}
	return o.GetWorkflowResponse
}

func (o *GetWorkflowResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *GetWorkflowResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}
