// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

package shared

type UpdateClientResponse struct {
	Data *Client `json:"data,omitempty"`
}

func (o *UpdateClientResponse) GetData() *Client {
	if o == nil {
		return nil
	}
	return o.Data
}
