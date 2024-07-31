package moneycorp

import (
	"context"

	"github.com/formancehq/paymentsv3/internal/plugins/models"
	"github.com/formancehq/paymentsv3/internal/plugins/public/moneycorp/client"
)

type Plugin struct {
	client *client.Client
}

func (p *Plugin) Install(ctx context.Context, req models.InstallRequest) (models.InstallResponse, error) {
	config, err := unmarshalAndValidateConfig(req.Config)
	if err != nil {
		return models.InstallResponse{}, err
	}

	client, err := client.NewClient(config.ClientID, config.APIKey, config.Endpoint)
	if err != nil {
		return models.InstallResponse{}, err
	}
	p.client = client

	return models.InstallResponse{
		Capabilities: capabilities,
		Workflow:     workflow(),
	}, nil
}

func (p Plugin) FetchNextAccounts(ctx context.Context, req models.FetchNextAccountsRequest) (models.FetchNextAccountsResponse, error) {
	return p.fetchNextAccounts(ctx, req)
}

func (p Plugin) FetchNextPayments(ctx context.Context, req models.FetchNextPaymentsRequest) (models.FetchNextPaymentsResponse, error) {
	return p.fetchNextPayments(ctx, req)
}

func (p Plugin) FetchNextOthers(ctx context.Context, req models.FetchNextOthersRequest) (models.FetchNextOthersResponse, error) {
	return models.FetchNextOthersResponse{}, nil
}

var _ models.Plugin = &Plugin{}
