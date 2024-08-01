package activities

import (
	"context"
	"encoding/json"

	"github.com/formancehq/paymentsv3/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) FetchNextPayments(ctx context.Context, request models.FetchNextPaymentsRequest) (models.FetchNextPaymentsResponse, error) {
	resp, err := a.plugin.FetchNextPayments(ctx, request)
	if err != nil {
		// TODO(polo): temporal errors
		return models.FetchNextPaymentsResponse{}, err
	}

	return resp, nil
}

var FetchNextPaymentsActivity = Activities{}.FetchNextPayments

func FetchNextPayments(ctx workflow.Context, fromPayload, state json.RawMessage, pageSize int) (models.FetchNextPaymentsResponse, error) {
	ret := models.FetchNextPaymentsResponse{}
	if err := executeActivity(ctx, FetchNextOthersActivity, ret, models.FetchNextPaymentsRequest{
		FromPayload: fromPayload,
		State:       state,
		PageSize:    pageSize,
	}); err != nil {
		return models.FetchNextPaymentsResponse{}, err
	}
	return ret, nil
}
