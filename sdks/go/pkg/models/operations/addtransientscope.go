// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package operations

import (
	"net/http"
)

type AddTransientScopeRequest struct {
	// Scope ID
	ScopeID string `pathParam:"style=simple,explode=false,name=scopeId"`
	// Transient scope ID
	TransientScopeID string `pathParam:"style=simple,explode=false,name=transientScopeId"`
}

func (o *AddTransientScopeRequest) GetScopeID() string {
	if o == nil {
		return ""
	}
	return o.ScopeID
}

func (o *AddTransientScopeRequest) GetTransientScopeID() string {
	if o == nil {
		return ""
	}
	return o.TransientScopeID
}

type AddTransientScopeResponse struct {
	ContentType string
	StatusCode  int
	RawResponse *http.Response
}

func (o *AddTransientScopeResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *AddTransientScopeResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *AddTransientScopeResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}
