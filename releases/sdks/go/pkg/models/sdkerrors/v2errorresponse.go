// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

package sdkerrors

import (
	"encoding/json"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
)

// V2ErrorResponse - Error
type V2ErrorResponse struct {
	Details      *string             `json:"details,omitempty"`
	ErrorCode    shared.V2ErrorsEnum `json:"errorCode"`
	ErrorMessage string              `json:"errorMessage"`
}

var _ error = &V2ErrorResponse{}

func (e *V2ErrorResponse) Error() string {
	data, _ := json.Marshal(e)
	return string(data)
}
