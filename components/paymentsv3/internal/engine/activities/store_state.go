package activities

import (
	"context"

	"github.com/formancehq/paymentsv3/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) StoreState(ctx context.Context, state models.State) error {
	return a.storage.UpsertState(ctx, state)
}

var StoreStateActivity = Activities{}.StoreState

func StoreState(ctx workflow.Context, state models.State) error {
	return executeActivity(ctx, StoreStateActivity, nil, state)
}
