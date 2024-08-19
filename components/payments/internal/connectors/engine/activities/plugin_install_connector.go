package activities

import (
	"context"
	"encoding/json"

	"github.com/formancehq/payments/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) PluginInstallConnector(ctx context.Context, plugin models.Plugin, request models.InstallRequest) (models.InstallResponse, error) {
	resp, err := plugin.Install(ctx, request)
	if err != nil {
		// TODO(polo): temporal errors
		return models.InstallResponse{}, err
	}
	return resp, err
}

var PluginInstallConnectorActivity = Activities{}.PluginInstallConnector

func PluginInstallConnector(ctx workflow.Context, plugin models.Plugin, config json.RawMessage) (models.InstallResponse, error) {
	ret := models.InstallResponse{}
	if err := executeActivity(ctx, PluginInstallConnectorActivity, ret, plugin, models.InstallRequest{
		Config: config,
	}); err != nil {
		return models.InstallResponse{}, err
	}
	return ret, nil
}
