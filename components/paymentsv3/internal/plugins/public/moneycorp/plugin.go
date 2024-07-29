package moneycorp

import (
	"context"

	"github.com/formancehq/stack/components/paymentsv3/internal/plugins/models"
	"github.com/formancehq/stack/components/paymentsv3/internal/plugins/public/moneycorp/client"
)

type Plugin struct {
	client *client.Client
}

func (p *Plugin) Install(ctx context.Context, req models.InstallRequest) (models.InstallResponse, error) {
	config, err := unmarshalAndValidateConfig(req.Config)
	if err != nil {
		return models.InstallResponse{}, err
	}

	client, err := client.NewClient(config.ClientID, config.APIKey, config.Endpoint, config.PageSize)
	if err != nil {
		return models.InstallResponse{}, err
	}
	p.client = client

	return models.InstallResponse{
		Capabilities: capabilities,
		Workflow:     workflow(),
	}, nil
}

func (p Plugin) FetchAccounts(ctx context.Context, req models.FetchAccountsRequest) (models.FetchAccountsResponse, error) {
	return p.fetchAccounts(ctx, req)
}

func (p Plugin) FetchPayments(ctx context.Context, req models.FetchPaymentsRequest) (models.FetchPaymentsResponse, error) {
	return p.fetchPayments(ctx, req)
}

func (p Plugin) FetchOthers(ctx context.Context, req models.FetchOthersRequest) (models.FetchOthersResponse, error) {
	return models.FetchOthersResponse{}, nil
}

var _ models.Plugin = &Plugin{}
