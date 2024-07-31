package activities

import (
	"context"
	"encoding/json"

	"github.com/formancehq/paymentsv3/internal/plugins/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) FetchNextOthers(ctx context.Context, request models.FetchNextOthersRequest) (models.FetchNextOthersResponse, error) {
	resp, err := a.plugin.FetchNextOthers(ctx, request)
	if err != nil {
		// TODO(polo): temporal errors
		return models.FetchNextOthersResponse{}, err
	}

	return resp, nil
}

var FetchNextOthersActivity = Activities{}.FetchNextOthers

func FetchNextOthers(ctx workflow.Context, name string, fromPayload, state json.RawMessage, pageSize int) (models.FetchNextOthersResponse, error) {
	ret := models.FetchNextOthersResponse{}
	if err := executeActivity(ctx, FetchNextOthersActivity, ret, models.FetchNextOthersRequest{
		Name:        name,
		FromPayload: fromPayload,
		State:       state,
		PageSize:    pageSize,
	}); err != nil {
		return models.FetchNextOthersResponse{}, err
	}
	return ret, nil
}
