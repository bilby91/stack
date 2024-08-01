// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

package shared

import (
	"github.com/formancehq/formance-sdk-go/v2/pkg/utils"
	"math/big"
)

type V2AggregateBalancesResponse struct {
	Data map[string]*big.Int `json:"data"`
}

func (v V2AggregateBalancesResponse) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(v, "", false)
}

func (v *V2AggregateBalancesResponse) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &v, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *V2AggregateBalancesResponse) GetData() map[string]*big.Int {
	if o == nil {
		return map[string]*big.Int{}
	}
	return o.Data
}
