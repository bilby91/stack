package grpc

import (
	"context"

	"github.com/formancehq/paymentsv3/internal/grpc/proto/services"
)

type GRPCClient struct {
	client services.PluginClient
}

func (c *GRPCClient) Install(ctx context.Context, req *services.InstallRequest) (*services.InstallResponse, error) {
	return c.client.Install(ctx, req)
}

func (c *GRPCClient) FetchNextAccounts(ctx context.Context, req *services.FetchNextAccountsRequest) (*services.FetchNextAccountsResponse, error) {
	return c.client.FetchNextAccounts(ctx, req)
}

func (c *GRPCClient) FetchNextPayments(ctx context.Context, req *services.FetchNextPaymentsRequest) (*services.FetchNextPaymentsResponse, error) {
	return c.client.FetchNextPayments(ctx, req)
}

func (c *GRPCClient) FetchNextOthers(ctx context.Context, req *services.FetchNextOthersRequest) (*services.FetchNextOthersResponse, error) {
	return c.client.FetchNextOthers(ctx, req)
}
