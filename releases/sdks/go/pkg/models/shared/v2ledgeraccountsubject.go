// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

package shared

type V2LedgerAccountSubject struct {
	Identifier string `json:"identifier"`
	Type       string `json:"type"`
}

func (o *V2LedgerAccountSubject) GetIdentifier() string {
	if o == nil {
		return ""
	}
	return o.Identifier
}

func (o *V2LedgerAccountSubject) GetType() string {
	if o == nil {
		return ""
	}
	return o.Type
}
