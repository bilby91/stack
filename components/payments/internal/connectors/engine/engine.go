package engine

import (
	"context"
	"encoding/json"

	"github.com/formancehq/payments/internal/connectors/engine/plugins"
	"github.com/formancehq/payments/internal/connectors/engine/workflow"
	"github.com/formancehq/payments/internal/models"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
)

type Engine interface {
	InstallConnector(ctx context.Context, provider string, rawConfig json.RawMessage) (models.ConnectorID, error)
	UninstallConnector(ctx context.Context, connectorID models.ConnectorID) error
	CreateBankAccount(ctx context.Context, bankAccountID uuid.UUID, connectorID models.ConnectorID) (*models.BankAccount, error)
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
	config := models.DefaultConfig()
	if err := json.Unmarshal(rawConfig, &config); err != nil {
		return models.ConnectorID{}, err
	}

	if err := config.Validate(); err != nil {
		return models.ConnectorID{}, errors.Wrap(ErrValidation, err.Error())
	}

	connectorID := models.ConnectorID{
		Reference: config.Name,
		Provider:  provider,
	}

	plugin, err := e.plugins.RegisterPlugin(connectorID)
	if err != nil {
		return models.ConnectorID{}, handlePluginError(err)
	}

	err = e.workers.AddWorker(connectorID)
	if err != nil {
		return models.ConnectorID{}, err
	}

	// Launch the workflow
	run, err := e.temporalClient.ExecuteWorkflow(
		ctx,
		client.StartWorkflowOptions{
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
		return handlePluginError(err)
	}

	return nil
}

func (e *engine) CreateBankAccount(ctx context.Context, bankAccountID uuid.UUID, connectorID models.ConnectorID) (*models.BankAccount, error) {
	plugin, err := e.plugins.Get(connectorID)
	if err != nil {
		return nil, handlePluginError(err)
	}

	run, err := e.temporalClient.ExecuteWorkflow(
		ctx,
		client.StartWorkflowOptions{
			TaskQueue:                                connectorID.Reference,
			WorkflowIDReusePolicy:                    enums.WORKFLOW_ID_REUSE_POLICY_REJECT_DUPLICATE,
			WorkflowExecutionErrorWhenAlreadyStarted: false,
		},
		workflow.RunCreateBankAccount,
		plugin,
		workflow.CreateBankAccount{
			ConnectorID:   connectorID,
			BankAccountID: bankAccountID,
		},
	)
	if err != nil {
		return nil, err
	}

	var bankAccount models.BankAccount
	// Wait for bank account creation to complete
	if err := run.Get(ctx, &bankAccount); err != nil {
		return nil, err
	}

	return &bankAccount, nil
}

var _ Engine = &engine{}
