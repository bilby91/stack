// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

package shared

import (
	"github.com/formancehq/formance-sdk-go/v2/pkg/utils"
	"time"
)

type DebitWalletRequest struct {
	Amount      Monetary `json:"amount"`
	Balances    []string `json:"balances,omitempty"`
	Description *string  `json:"description,omitempty"`
	Destination *Subject `json:"destination,omitempty"`
	// Metadata associated with the wallet.
	Metadata map[string]string `json:"metadata"`
	// Set to true to create a pending hold. If false, the wallet will be debited immediately.
	Pending *bool `json:"pending,omitempty"`
	// cannot be used in conjunction with `pending` property
	Timestamp *time.Time `json:"timestamp,omitempty"`
}

func (d DebitWalletRequest) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(d, "", false)
}

func (d *DebitWalletRequest) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &d, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *DebitWalletRequest) GetAmount() Monetary {
	if o == nil {
		return Monetary{}
	}
	return o.Amount
}

func (o *DebitWalletRequest) GetBalances() []string {
	if o == nil {
		return nil
	}
	return o.Balances
}

func (o *DebitWalletRequest) GetDescription() *string {
	if o == nil {
		return nil
	}
	return o.Description
}

func (o *DebitWalletRequest) GetDestination() *Subject {
	if o == nil {
		return nil
	}
	return o.Destination
}

func (o *DebitWalletRequest) GetDestinationAccount() *LedgerAccountSubject {
	if v := o.GetDestination(); v != nil {
		return v.LedgerAccountSubject
	}
	return nil
}

func (o *DebitWalletRequest) GetDestinationWallet() *WalletSubject {
	if v := o.GetDestination(); v != nil {
		return v.WalletSubject
	}
	return nil
}

func (o *DebitWalletRequest) GetMetadata() map[string]string {
	if o == nil {
		return map[string]string{}
	}
	return o.Metadata
}

func (o *DebitWalletRequest) GetPending() *bool {
	if o == nil {
		return nil
	}
	return o.Pending
}

func (o *DebitWalletRequest) GetTimestamp() *time.Time {
	if o == nil {
		return nil
	}
	return o.Timestamp
}
