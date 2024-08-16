package activities

import (
	"context"

	"github.com/formancehq/paymentsv3/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) StorageStatesGet(ctx context.Context, id models.StateID) (models.State, error) {
	return a.storage.StatesGet(ctx, id)
}

var StorageStatesGetActivity = Activities{}.StorageStatesGet

func StorageStatesGet(ctx workflow.Context, id models.StateID) (models.State, error) {
	ret := models.State{}
	if err := executeActivity(ctx, StorageStatesGetActivity, ret, id); err != nil {
		return models.State{}, err
	}
	return ret, nil
}
