package activities

import (
	"context"
	"encoding/json"

	"github.com/formancehq/paymentsv3/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) FetchNextExternalAccounts(ctx context.Context, request models.FetchNextExternalAccountsRequest) (models.FetchNextExternalAccountsResponse, error) {
	resp, err := a.plugin.FetchNextExternalAccounts(ctx, request)
	if err != nil {
		// TODO(polo): temporal errors
		return models.FetchNextExternalAccountsResponse{}, err
	}

	return resp, nil
}

var FetchNextExternalAccountsActivity = Activities{}.FetchNextExternalAccounts

func FetchNextExternalAccounts(ctx workflow.Context, fromPayload, state json.RawMessage, pageSize int) (models.FetchNextExternalAccountsResponse, error) {
	ret := models.FetchNextExternalAccountsResponse{}
	if err := executeActivity(ctx, FetchNextExternalAccountsActivity, ret, models.FetchNextExternalAccountsRequest{
		FromPayload: fromPayload,
		State:       state,
		PageSize:    pageSize,
	}); err != nil {
		return models.FetchNextExternalAccountsResponse{}, err
	}
	return ret, nil
}
