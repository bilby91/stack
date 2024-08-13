package engine

import (
	"context"
	"encoding/json"

	"github.com/formancehq/paymentsv3/internal/connectors/engine/plugins"
	"github.com/formancehq/paymentsv3/internal/connectors/engine/workflow"
	"github.com/formancehq/paymentsv3/internal/models"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
)

type Engine interface {
	InstallConnector(ctx context.Context, provider string, rawConfig json.RawMessage) (models.ConnectorID, error)
	UninstallConnector(ctx context.Context, connectorID models.ConnectorID) error
}

type engine struct {
	temporalClient client.Client

	workers *Workers
	plugins plugins.Plugins
}

func New(workers *Workers, plugins plugins.Plugins) Engine {
	return &engine{
		workers: workers,
		plugins: plugins,
	}
}

func (e *engine) InstallConnector(ctx context.Context, provider string, rawConfig json.RawMessage) (models.ConnectorID, error) {
	var config models.Config
	if err := json.Unmarshal(rawConfig, &config); err != nil {
		return models.ConnectorID{}, err
	}

	connectorID := models.ConnectorID{
		Reference: config.Name,
		Provider:  provider,
	}

	plugin, err := e.plugins.RegisterPlugin(connectorID)
	if err != nil {
		return models.ConnectorID{}, err
	}

	err = e.workers.AddWorker(connectorID)
	if err != nil {
		return models.ConnectorID{}, err
	}

	// Launch the workflow without waiting for the result
	run, err := e.temporalClient.ExecuteWorkflow(
		ctx,
		client.StartWorkflowOptions{
			ID:                                       connectorID.String(),
			TaskQueue:                                connectorID.Reference,
			WorkflowIDReusePolicy:                    enums.WORKFLOW_ID_REUSE_POLICY_REJECT_DUPLICATE,
			WorkflowExecutionErrorWhenAlreadyStarted: false,
		},
		workflow.RunInstallConnector,
		plugin,
		workflow.InstallConnector{
			ConnectorID: connectorID,
			RawConfig:   rawConfig,
		},
	)
	if err != nil {
		return models.ConnectorID{}, err
	}

	// Wait for installation to complete
	if err := run.Get(ctx, nil); err != nil {
		return models.ConnectorID{}, err
	}

	return connectorID, nil
}

func (e *engine) UninstallConnector(ctx context.Context, connectorID models.ConnectorID) error {
	run, err := e.temporalClient.ExecuteWorkflow(
		ctx,
		client.StartWorkflowOptions{
			ID:                                       connectorID.String(),
			TaskQueue:                                connectorID.Reference,
			WorkflowIDReusePolicy:                    enums.WORKFLOW_ID_REUSE_POLICY_REJECT_DUPLICATE,
			WorkflowExecutionErrorWhenAlreadyStarted: false,
		},
		workflow.RunUninstallConnector,
		workflow.UninstallConnector{
			ConnectorID: connectorID,
		},
	)
	if err != nil {
		return err
	}

	// Wait for uninstallation to complete
	if err := run.Get(ctx, nil); err != nil {
		return err
	}

	if err := e.workers.RemoveWorker(connectorID); err != nil {
		return err
	}

	if err := e.plugins.UnregisterPlugin(connectorID); err != nil {
		return err
	}

	return nil
}

var _ Engine = &engine{}
