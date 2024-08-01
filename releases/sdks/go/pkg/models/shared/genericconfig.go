// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

package shared

import (
	"github.com/formancehq/formance-sdk-go/v2/pkg/utils"
)

type GenericConfig struct {
	APIKey   string `json:"apiKey"`
	Endpoint string `json:"endpoint"`
	Name     string `json:"name"`
	// The frequency at which the connector will try to fetch new BalanceTransaction objects from the API.
	//
	PollingPeriod *string `default:"120s" json:"pollingPeriod"`
}

func (g GenericConfig) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(g, "", false)
}

func (g *GenericConfig) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &g, "", false, true); err != nil {
		return err
	}
	return nil
}

func (o *GenericConfig) GetAPIKey() string {
	if o == nil {
		return ""
	}
	return o.APIKey
}

func (o *GenericConfig) GetEndpoint() string {
	if o == nil {
		return ""
	}
	return o.Endpoint
}

func (o *GenericConfig) GetName() string {
	if o == nil {
		return ""
	}
	return o.Name
}

func (o *GenericConfig) GetPollingPeriod() *string {
	if o == nil {
		return nil
	}
	return o.PollingPeriod
}
