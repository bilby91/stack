// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

type ActivityGetWallet struct {
	ID string `json:"id"`
}

func (o *ActivityGetWallet) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}
