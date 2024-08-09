package activities

import (
	"context"

	"github.com/formancehq/paymentsv3/internal/models"
	"go.temporal.io/sdk/workflow"
)

type StoreTasksTreeRequest struct {
	ConnectorID models.ConnectorID
	Workflow    models.Workflow
}

func (a Activities) StoreTasksTree(ctx context.Context, request StoreTasksTreeRequest) error {
	return a.storage.UpsertWorkflow(ctx, request.ConnectorID, request.Workflow)
}

var StoreTasksTreeActivity = Activities{}.StoreTasksTree

func StoreTasksTree(ctx workflow.Context, connectorID models.ConnectorID, workflow models.Workflow) error {
	return executeActivity(ctx, StoreTasksTreeActivity, nil, StoreTasksTreeRequest{
		ConnectorID: connectorID,
		Workflow:    workflow,
	})
}
