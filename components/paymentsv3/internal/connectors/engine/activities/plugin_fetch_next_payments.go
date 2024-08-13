package activities

import (
	"context"
	"encoding/json"

	"github.com/formancehq/paymentsv3/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) PluginFetchNextPayments(ctx context.Context, plugin models.Plugin, request models.FetchNextPaymentsRequest) (models.FetchNextPaymentsResponse, error) {
	resp, err := plugin.FetchNextPayments(ctx, request)
	if err != nil {
		// TODO(polo): temporal errors
		return models.FetchNextPaymentsResponse{}, err
	}

	return resp, nil
}

var PluginFetchNextPaymentsActivity = Activities{}.PluginFetchNextPayments

func PluginFetchNextPayments(ctx workflow.Context, plugin models.Plugin, fromPayload, state json.RawMessage, pageSize int) (models.FetchNextPaymentsResponse, error) {
	ret := models.FetchNextPaymentsResponse{}
	if err := executeActivity(ctx, PluginFetchNextOthersActivity, ret, plugin, models.FetchNextPaymentsRequest{
		FromPayload: fromPayload,
		State:       state,
		PageSize:    pageSize,
	}); err != nil {
		return models.FetchNextPaymentsResponse{}, err
	}
	return ret, nil
}
