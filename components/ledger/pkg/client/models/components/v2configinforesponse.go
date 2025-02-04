// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

package components

type V2ConfigInfoResponse struct {
	Server  string `json:"server"`
	Version string `json:"version"`
}

func (o *V2ConfigInfoResponse) GetServer() string {
	if o == nil {
		return ""
	}
	return o.Server
}

func (o *V2ConfigInfoResponse) GetVersion() string {
	if o == nil {
		return ""
	}
	return o.Version
}
