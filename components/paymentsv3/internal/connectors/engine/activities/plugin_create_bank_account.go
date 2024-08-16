package activities

import (
	"context"

	"github.com/formancehq/paymentsv3/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) PluginCreateBankAccount(ctx context.Context, plugin models.Plugin, request models.CreateBankAccountRequest) (models.CreateBankAccountResponse, error) {
	resp, err := plugin.CreateBankAccount(ctx, request)
	if err != nil {
		// TODO(polo): temporal errors
		return models.CreateBankAccountResponse{}, err
	}
	return resp, nil
}

var PluginCreateBankAccountActivity = Activities{}.PluginCreateBankAccount

func PluginCreateBankAccount(ctx workflow.Context, plugin models.Plugin, request models.CreateBankAccountRequest) (models.CreateBankAccountResponse, error) {
	ret := models.CreateBankAccountResponse{}
	if err := executeActivity(ctx, PluginCreateBankAccountActivity, ret, plugin, request); err != nil {
		return models.CreateBankAccountResponse{}, err
	}
	return ret, nil
}
