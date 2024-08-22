package plugins

import (
	"context"
	"errors"

	"github.com/formancehq/payments/internal/connectors/grpc"
	"github.com/formancehq/payments/internal/connectors/grpc/proto"
	"github.com/formancehq/payments/internal/connectors/grpc/proto/services"
	"github.com/formancehq/payments/internal/models"
	"github.com/gogo/status"
	"google.golang.org/grpc/codes"
)

type impl struct {
	plugin models.Plugin
}

func NewGRPCImplem(plugin models.Plugin) *impl {
	return &impl{plugin: plugin}
}

func (i *impl) Install(ctx context.Context, req *services.InstallRequest) (*services.InstallResponse, error) {
	resp, err := i.plugin.Install(ctx, models.InstallRequest{
		Config: req.Config,
	})
	if err != nil {
		return nil, translateErrorToGRPC(err)
	}

	capabilities := make([]proto.Capability, 0, len(resp.Capabilities))
	for _, capability := range resp.Capabilities {
		capabilities = append(capabilities, proto.Capability(capability))
	}

	return &services.InstallResponse{
		Capabilities: capabilities,
		Workflow:     grpc.TranslateWorkflow(resp.Workflow),
	}, nil
}

func (i *impl) FetchNextAccounts(ctx context.Context, req *services.FetchNextAccountsRequest) (*services.FetchNextAccountsResponse, error) {
	resp, err := i.plugin.FetchNextAccounts(ctx, models.FetchNextAccountsRequest{
		FromPayload: req.FromPayload,
		State:       req.State,
		PageSize:    int(req.PageSize),
	})
	if err != nil {
		return nil, translateErrorToGRPC(err)
	}

	accounts := make([]*proto.Account, 0, len(resp.Accounts))
	for _, account := range resp.Accounts {
		accounts = append(accounts, grpc.TranslateAccount(account))
	}

	return &services.FetchNextAccountsResponse{
		Accounts: accounts,
		NewState: resp.NewState,
		HasMore:  resp.HasMore,
	}, nil
}

func (i *impl) FetchNextExternalAccounts(ctx context.Context, req *services.FetchNextExternalAccountsRequest) (*services.FetchNextExternalAccountsResponse, error) {
	resp, err := i.plugin.FetchNextExternalAccounts(ctx, models.FetchNextExternalAccountsRequest{
		FromPayload: req.FromPayload,
		State:       req.State,
		PageSize:    int(req.PageSize),
	})
	if err != nil {
		return nil, translateErrorToGRPC(err)
	}

	externalAccounts := make([]*proto.Account, 0, len(resp.ExternalAccounts))
	for _, account := range resp.ExternalAccounts {
		externalAccounts = append(externalAccounts, grpc.TranslateAccount(account))
	}

	return &services.FetchNextExternalAccountsResponse{
		Accounts: externalAccounts,
		NewState: resp.NewState,
		HasMore:  resp.HasMore,
	}, nil
}

func (i *impl) FetchNextPayments(ctx context.Context, req *services.FetchNextPaymentsRequest) (*services.FetchNextPaymentsResponse, error) {
	resp, err := i.plugin.FetchNextPayments(ctx, models.FetchNextPaymentsRequest{
		FromPayload: req.FromPayload,
		State:       req.State,
		PageSize:    int(req.PageSize),
	})
	if err != nil {
		return nil, translateErrorToGRPC(err)
	}

	payments := make([]*proto.Payment, 0, len(resp.Payments))
	for _, payment := range resp.Payments {
		payments = append(payments, grpc.TranslatePayment(payment))
	}

	return &services.FetchNextPaymentsResponse{
		Payments: payments,
		NewState: resp.NewState,
		HasMore:  resp.HasMore,
	}, nil
}

func (i *impl) FetchNextOthers(ctx context.Context, req *services.FetchNextOthersRequest) (*services.FetchNextOthersResponse, error) {
	resp, err := i.plugin.FetchNextOthers(ctx, models.FetchNextOthersRequest{
		FromPayload: req.FromPayload,
		State:       req.State,
		PageSize:    int(req.PageSize),
		Name:        req.Name,
	})
	if err != nil {
		return nil, translateErrorToGRPC(err)
	}

	others := make([]*proto.Other, 0, len(resp.Others))
	for _, other := range resp.Others {
		others = append(others, &proto.Other{
			Id:    other.ID,
			Other: other.Other,
		})
	}

	return &services.FetchNextOthersResponse{
		Others:   others,
		NewState: resp.NewState,
		HasMore:  resp.HasMore,
	}, nil
}

func (i *impl) CreateBankAccount(ctx context.Context, req *services.CreateBankAccountRequest) (*services.CreateBankAccountResponse, error) {
	resp, err := i.plugin.CreateBankAccount(ctx, models.CreateBankAccountRequest{
		BankAccount: grpc.TranslateProtoBankAccount(req.BankAccount),
	})
	if err != nil {
		return nil, translateErrorToGRPC(err)
	}

	return &services.CreateBankAccountResponse{
		RelatedAccount: grpc.TranslateAccount(resp.RelatedAccount),
	}, nil
}

var _ grpc.PSP = &impl{}

func translateErrorToGRPC(err error) error {
	switch {
	case errors.Is(err, models.ErrInvalidConfig):
		return status.Errorf(codes.InvalidArgument, err.Error())
	default:
		return status.Errorf(codes.Internal, err.Error())
	}
}
