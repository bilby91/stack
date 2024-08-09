package activities

import (
	"context"

	"github.com/formancehq/paymentsv3/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) FetchState(ctx context.Context, id models.StateID) (models.State, error) {
	return a.storage.GetState(ctx, id)
}

var FetchStateActivity = Activities{}.FetchState

func FetchState(ctx workflow.Context, id models.StateID) (models.State, error) {
	ret := models.State{}
	if err := executeActivity(ctx, FetchStateActivity, ret, id); err != nil {
		return models.State{}, err
	}
	return ret, nil
}
