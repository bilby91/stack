package activities

import (
	"context"

	"github.com/formancehq/paymentsv3/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) StorageUpdateInstance(ctx context.Context, instance models.Instance) error {
	return a.storage.UpdateInstance(ctx, instance)
}

var StorageUpdateInstanceActivity = Activities{}.StorageUpdateInstance

func StorageUpdateInstance(ctx workflow.Context, instance models.Instance) error {
	return executeActivity(ctx, StorageUpdateInstanceActivity, nil, instance)
}
