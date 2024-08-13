package activities

import (
	"context"
	"encoding/json"

	"github.com/formancehq/paymentsv3/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) PluginFetchNextExternalAccounts(ctx context.Context, plugin models.Plugin, request models.FetchNextExternalAccountsRequest) (models.FetchNextExternalAccountsResponse, error) {
	resp, err := plugin.FetchNextExternalAccounts(ctx, request)
	if err != nil {
		// TODO(polo): temporal errors
		return models.FetchNextExternalAccountsResponse{}, err
	}

	return resp, nil
}

var PluginFetchNextExternalAccountsActivity = Activities{}.PluginFetchNextExternalAccounts

func PluginFetchNextExternalAccounts(ctx workflow.Context, plugin models.Plugin, fromPayload, state json.RawMessage, pageSize int) (models.FetchNextExternalAccountsResponse, error) {
	ret := models.FetchNextExternalAccountsResponse{}
	if err := executeActivity(ctx, PluginFetchNextExternalAccountsActivity, ret, plugin, models.FetchNextExternalAccountsRequest{
		FromPayload: fromPayload,
		State:       state,
		PageSize:    pageSize,
	}); err != nil {
		return models.FetchNextExternalAccountsResponse{}, err
	}
	return ret, nil
}
