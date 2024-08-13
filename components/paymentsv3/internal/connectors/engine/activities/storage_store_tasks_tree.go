package activities

import (
	"context"

	"github.com/formancehq/paymentsv3/internal/models"
	"go.temporal.io/sdk/workflow"
)

type StoreTasksTreeRequest struct {
	ConnectorID models.ConnectorID
	Workflow    models.Tasks
}

func (a Activities) StorageStoreTasksTree(ctx context.Context, request StoreTasksTreeRequest) error {
	return a.storage.UpsertTasks(ctx, request.ConnectorID, request.Workflow)
}

var StorageStoreTasksTreeActivity = Activities{}.StorageStoreTasksTree

func StorageStoreTasksTree(ctx workflow.Context, connectorID models.ConnectorID, workflow models.Tasks) error {
	return executeActivity(ctx, StorageStoreTasksTreeActivity, nil, StoreTasksTreeRequest{
		ConnectorID: connectorID,
		Workflow:    workflow,
	})
}
