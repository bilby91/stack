package activities

import (
	"context"

	"github.com/formancehq/paymentsv3/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) StorageStoreState(ctx context.Context, state models.State) error {
	return a.storage.UpsertState(ctx, state)
}

var StorageStoreStateActivity = Activities{}.StorageStoreState

func StorageStoreState(ctx workflow.Context, state models.State) error {
	return executeActivity(ctx, StorageStoreStateActivity, nil, state)
}
