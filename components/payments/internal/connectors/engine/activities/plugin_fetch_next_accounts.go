package activities

import (
	"context"
	"encoding/json"

	"github.com/formancehq/payments/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) PluginFetchNextAccounts(ctx context.Context, plugin models.Plugin, request models.FetchNextAccountsRequest) (models.FetchNextAccountsResponse, error) {
	resp, err := plugin.FetchNextAccounts(ctx, request)
	if err != nil {
		// TODO(polo): temporal errors
		return models.FetchNextAccountsResponse{}, err
	}
	return resp, nil
}

var PluginFetchNextAccountsActivity = Activities{}.PluginFetchNextAccounts

func PluginFetchNextAccounts(ctx workflow.Context, plugin models.Plugin, fromPayload, state json.RawMessage, pageSize int) (models.FetchNextAccountsResponse, error) {
	ret := models.FetchNextAccountsResponse{}
	if err := executeActivity(ctx, PluginFetchNextAccountsActivity, ret, plugin, models.FetchNextAccountsRequest{
		FromPayload: fromPayload,
		State:       state,
		PageSize:    pageSize,
	}); err != nil {
		return models.FetchNextAccountsResponse{}, err
	}
	return ret, nil
}
