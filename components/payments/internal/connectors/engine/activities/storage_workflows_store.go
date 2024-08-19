package activities

import (
	"context"

	"github.com/formancehq/payments/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) StorageWorkflowsStore(ctx context.Context, workflow models.Workflow) error {
	return a.storage.WorkflowsUpsert(ctx, workflow)
}

var StorageWorkflowsStoreActivity = Activities{}.StorageWorkflowsStore

func StorageWorkflowsStore(ctx workflow.Context, workflow models.Workflow) error {
	return executeActivity(ctx, StorageWorkflowsStoreActivity, nil, workflow)
}
