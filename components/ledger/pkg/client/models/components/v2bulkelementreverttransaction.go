// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

package components

import (
	"github.com/formancehq/stack/ledger/client/internal/utils"
	"math/big"
)

type V2BulkElementRevertTransactionData struct {
	ID              *big.Int `json:"id"`
	Force           *bool    `json:"force,omitempty"`
	AtEffectiveDate *bool    `json:"atEffectiveDate,omitempty"`
}

func (v V2BulkElementRevertTransactionData) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(v, "", false)
}

func (v *V2BulkElementRevertTransactionData) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &v, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *V2BulkElementRevertTransactionData) GetID() *big.Int {
	if o == nil {
		return big.NewInt(0)
	}
	return o.ID
}

func (o *V2BulkElementRevertTransactionData) GetForce() *bool {
	if o == nil {
		return nil
	}
	return o.Force
}

func (o *V2BulkElementRevertTransactionData) GetAtEffectiveDate() *bool {
	if o == nil {
		return nil
	}
	return o.AtEffectiveDate
}

type V2BulkElementRevertTransaction struct {
	Action string                              `json:"action"`
	Ik     *string                             `json:"ik,omitempty"`
	Data   *V2BulkElementRevertTransactionData `json:"data,omitempty"`
}

func (o *V2BulkElementRevertTransaction) GetAction() string {
	if o == nil {
		return ""
	}
	return o.Action
}

func (o *V2BulkElementRevertTransaction) GetIk() *string {
	if o == nil {
		return nil
	}
	return o.Ik
}

func (o *V2BulkElementRevertTransaction) GetData() *V2BulkElementRevertTransactionData {
	if o == nil {
		return nil
	}
	return o.Data
}
