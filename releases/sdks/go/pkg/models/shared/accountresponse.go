// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

package shared

type AccountResponse struct {
	Data AccountWithVolumesAndBalances `json:"data"`
}

func (o *AccountResponse) GetData() AccountWithVolumesAndBalances {
	if o == nil {
		return AccountWithVolumesAndBalances{}
	}
	return o.Data
}
