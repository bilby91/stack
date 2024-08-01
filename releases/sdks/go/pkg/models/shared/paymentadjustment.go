// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

package shared

import (
	"github.com/formancehq/formance-sdk-go/v2/pkg/utils"
	"math/big"
	"time"
)

type PaymentAdjustmentRaw struct {
}

type PaymentAdjustment struct {
	Amount    *big.Int             `json:"amount"`
	CreatedAt time.Time            `json:"createdAt"`
	Raw       PaymentAdjustmentRaw `json:"raw"`
	Reference string               `json:"reference"`
	Status    PaymentStatus        `json:"status"`
}

func (p PaymentAdjustment) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(p, "", false)
}

func (p *PaymentAdjustment) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &p, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *PaymentAdjustment) GetAmount() *big.Int {
	if o == nil {
		return big.NewInt(0)
	}
	return o.Amount
}

func (o *PaymentAdjustment) GetCreatedAt() time.Time {
	if o == nil {
		return time.Time{}
	}
	return o.CreatedAt
}

func (o *PaymentAdjustment) GetRaw() PaymentAdjustmentRaw {
	if o == nil {
		return PaymentAdjustmentRaw{}
	}
	return o.Raw
}

func (o *PaymentAdjustment) GetReference() string {
	if o == nil {
		return ""
	}
	return o.Reference
}

func (o *PaymentAdjustment) GetStatus() PaymentStatus {
	if o == nil {
		return PaymentStatus("")
	}
	return o.Status
}
