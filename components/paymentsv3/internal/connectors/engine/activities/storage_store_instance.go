package activities

import (
	"context"

	"github.com/formancehq/paymentsv3/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) StorageStoreInstance(ctx context.Context, instance models.Instance) error {
	return a.storage.InsertNewInstance(ctx, instance)
}

var StorageStoreInstanceActivity = Activities{}.StorageStoreInstance

func StorageStoreInstance(ctx workflow.Context, instance models.Instance) error {
	return executeActivity(ctx, StorageStoreInstanceActivity, nil, instance)
}
