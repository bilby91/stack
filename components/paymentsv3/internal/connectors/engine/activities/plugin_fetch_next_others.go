package activities

import (
	"context"
	"encoding/json"

	"github.com/formancehq/paymentsv3/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) PluginFetchNextOthers(ctx context.Context, plugin models.Plugin, request models.FetchNextOthersRequest) (models.FetchNextOthersResponse, error) {
	resp, err := plugin.FetchNextOthers(ctx, request)
	if err != nil {
		// TODO(polo): temporal errors
		return models.FetchNextOthersResponse{}, err
	}

	return resp, nil
}

var PluginFetchNextOthersActivity = Activities{}.PluginFetchNextOthers

func PluginFetchNextOthers(ctx workflow.Context, plugin models.Plugin, name string, fromPayload, state json.RawMessage, pageSize int) (models.FetchNextOthersResponse, error) {
	ret := models.FetchNextOthersResponse{}
	if err := executeActivity(ctx, PluginFetchNextOthersActivity, ret, plugin, models.FetchNextOthersRequest{
		Name:        name,
		FromPayload: fromPayload,
		State:       state,
		PageSize:    pageSize,
	}); err != nil {
		return models.FetchNextOthersResponse{}, err
	}
	return ret, nil
}
