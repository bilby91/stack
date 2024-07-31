package activities

import (
	"context"
	"encoding/json"

	"github.com/formancehq/paymentsv3/internal/plugins/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) FetchNextAccounts(ctx context.Context, request models.FetchNextAccountsRequest) (models.FetchNextAccountsResponse, error) {
	resp, err := a.plugin.FetchNextAccounts(ctx, request)
	if err != nil {
		// TODO(polo): temporal errors
		return models.FetchNextAccountsResponse{}, err
	}

	return resp, nil
}

var FetchNextAccountsActivity = Activities{}.FetchNextAccounts

func FetchNextAccounts(ctx workflow.Context, fromPayload, state json.RawMessage, pageSize int) (models.FetchNextAccountsResponse, error) {
	ret := models.FetchNextAccountsResponse{}
	if err := executeActivity(ctx, FetchNextAccountsActivity, ret, models.FetchNextAccountsRequest{
		FromPayload: fromPayload,
		State:       state,
		PageSize:    pageSize,
	}); err != nil {
		return models.FetchNextAccountsResponse{}, err
	}
	return ret, nil
}
