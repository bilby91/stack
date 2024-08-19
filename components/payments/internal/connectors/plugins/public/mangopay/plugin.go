package mangopay

import (
	"context"
	"errors"

	"github.com/formancehq/payments/internal/connectors/plugins/public/mangopay/client"
	"github.com/formancehq/payments/internal/models"
)

const (
	usersTaskName = "fetch_users"
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

func (p Plugin) FetchNextOthers(ctx context.Context, req models.FetchNextOthersRequest) (models.FetchNextOthersResponse, error) {
	switch req.Name {
	case usersTaskName:
		return p.fetchNextUsers(ctx, req)
	default:
		return models.FetchNextOthersResponse{}, errors.New("unknown task name")
	}
}
