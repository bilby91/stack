// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

package components

type V2StatsResponse struct {
	Data V2Stats `json:"data"`
}

func (o *V2StatsResponse) GetData() V2Stats {
	if o == nil {
		return V2Stats{}
	}
	return o.Data
}
