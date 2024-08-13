package activities

import (
	"context"

	"github.com/formancehq/paymentsv3/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) StorageStoreWorkflow(ctx context.Context, workflow models.Workflow) error {
	return a.storage.UpsertWorkflow(ctx, workflow)
}

var StorageStoreWorkflowActivity = Activities{}.StorageStoreWorkflow

func StorageStoreWorkflow(ctx workflow.Context, workflow models.Workflow) error {
	return executeActivity(ctx, StorageStoreWorkflowActivity, nil, workflow)
}
