// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

package shared

type V2BulkElementCreateTransaction struct {
	Action string             `json:"action"`
	Data   *V2PostTransaction `json:"data,omitempty"`
	Ik     *string            `json:"ik,omitempty"`
}

func (o *V2BulkElementCreateTransaction) GetAction() string {
	if o == nil {
		return ""
	}
	return o.Action
}

func (o *V2BulkElementCreateTransaction) GetData() *V2PostTransaction {
	if o == nil {
		return nil
	}
	return o.Data
}

func (o *V2BulkElementCreateTransaction) GetIk() *string {
	if o == nil {
		return nil
	}
	return o.Ik
}
