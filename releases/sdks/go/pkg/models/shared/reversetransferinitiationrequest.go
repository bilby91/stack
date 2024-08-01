// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

package shared

import (
	"github.com/formancehq/formance-sdk-go/v2/pkg/utils"
	"math/big"
)

type ReverseTransferInitiationRequest struct {
	Amount      *big.Int          `json:"amount"`
	Asset       string            `json:"asset"`
	Description string            `json:"description"`
	Metadata    map[string]string `json:"metadata"`
	Reference   string            `json:"reference"`
}

func (r ReverseTransferInitiationRequest) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(r, "", false)
}

func (r *ReverseTransferInitiationRequest) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &r, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *ReverseTransferInitiationRequest) GetAmount() *big.Int {
	if o == nil {
		return big.NewInt(0)
	}
	return o.Amount
}

func (o *ReverseTransferInitiationRequest) GetAsset() string {
	if o == nil {
		return ""
	}
	return o.Asset
}

func (o *ReverseTransferInitiationRequest) GetDescription() string {
	if o == nil {
		return ""
	}
	return o.Description
}

func (o *ReverseTransferInitiationRequest) GetMetadata() map[string]string {
	if o == nil {
		return nil
	}
	return o.Metadata
}

func (o *ReverseTransferInitiationRequest) GetReference() string {
	if o == nil {
		return ""
	}
	return o.Reference
}
