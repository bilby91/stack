package workflow

import (
	"encoding/json"
	"time"

	"github.com/formancehq/paymentsv3/internal/connectors/engine/activities"
	"github.com/formancehq/paymentsv3/internal/models"
	"github.com/pkg/errors"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/workflow"
)

type InstallConnector struct {
	ConnectorID models.ConnectorID
	RawConfig   json.RawMessage
}

func (w Workflow) runInstallConnector(
	ctx workflow.Context,
	plugin models.Plugin,
	installConnector InstallConnector,
) error {
	// First step: store the connector inside the database
	connector := models.Connector{
		ID:        installConnector.ConnectorID,
		Name:      installConnector.ConnectorID.Reference,
		CreatedAt: time.Now().UTC(), // TODO(polo): workflow.Now()
		Provider:  installConnector.ConnectorID.Provider,
		Config:    installConnector.RawConfig,
	}
	err := activities.StorageStoreConnector(infiniteRetryContext(ctx), connector)
	if err != nil {
		return errors.Wrap(err, "failed to store connector")
	}

	// Second step: install the connector via the plugin and get the list of
	// capabilities and the workflow of polling data
	installResponse, err := activities.PluginInstallConnector(infiniteRetryContext(ctx), plugin, installConnector.RawConfig)
	if err != nil {
		return errors.Wrap(err, "failed to install connector")
	}

	// Third step: store the workflow of the connector
	err = activities.StorageStoreTasksTree(infiniteRetryContext(ctx), installConnector.ConnectorID, installResponse.Workflow)
	if err != nil {
		return errors.Wrap(err, "failed to store tasks tree")
	}

	// TODO(polo): store the capabilities

	var config models.Config
	if err := json.Unmarshal(installConnector.RawConfig, &config); err != nil {
		return errors.Wrap(err, "failed to unmarshal config")
	}

	// Fourth step: launch the workflow tree
	if err := workflow.ExecuteChildWorkflow(
		workflow.WithChildOptions(
			ctx,
			workflow.ChildWorkflowOptions{
				TaskQueue:         installConnector.ConnectorID.Reference,
				ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
			},
		),
		Run,
		plugin,
		config,
		installConnector.ConnectorID,
		nil,
		[]models.TaskTree(installResponse.Workflow),
	).Get(ctx, nil); err != nil {
		return errors.Wrap(err, "running next workflow")
	}

	return nil
}

var RunInstallConnector any

func init() {
	RunInstallConnector = Workflow{}.runInstallConnector
}
