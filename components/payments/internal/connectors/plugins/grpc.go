package plugins

import (
	"context"
	"errors"
	"os"

	"github.com/formancehq/payments/internal/connectors/grpc"
	"github.com/formancehq/payments/internal/connectors/grpc/proto"
	"github.com/formancehq/payments/internal/connectors/grpc/proto/services"
	"github.com/formancehq/payments/internal/models"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type impl struct {
	logger hclog.Logger

	plugin models.Plugin
}

func NewGRPCImplem(plugin models.Plugin) *impl {
	logger := hclog.New(&hclog.LoggerOptions{
		Level:  hclog.Debug,
		Output: os.Stderr,
	})

	return &impl{
		logger: logger,
		plugin: plugin,
	}
}

func (i *impl) Install(ctx context.Context, req *services.InstallRequest) (*services.InstallResponse, error) {
	i.logger.Info("installing...")

	resp, err := i.plugin.Install(ctx, models.InstallRequest{
		Config: req.Config,
	})
	if err != nil {
		i.logger.Error("install failed: ", err)
		return nil, translateErrorToGRPC(err)
	}

	capabilities := make([]proto.Capability, 0, len(resp.Capabilities))
	for _, capability := range resp.Capabilities {
		capabilities = append(capabilities, proto.Capability(capability))
	}

	i.logger.Info("installed!")

	return &services.InstallResponse{
		Capabilities: capabilities,
		Workflow:     grpc.TranslateWorkflow(resp.Workflow),
	}, nil
}

func (i *impl) FetchNextAccounts(ctx context.Context, req *services.FetchNextAccountsRequest) (*services.FetchNextAccountsResponse, error) {
	i.logger.Info("fetching next accounts...")

	resp, err := i.plugin.FetchNextAccounts(ctx, models.FetchNextAccountsRequest{
		FromPayload: req.FromPayload,
		State:       req.State,
		PageSize:    int(req.PageSize),
	})
	if err != nil {
		i.logger.Error("fetching next accounts failed: ", err)
		return nil, translateErrorToGRPC(err)
	}

	accounts := make([]*proto.Account, 0, len(resp.Accounts))
	for _, account := range resp.Accounts {
		accounts = append(accounts, grpc.TranslateAccount(account))
	}

	i.logger.Info("fetched next accounts succeeded!")

	return &services.FetchNextAccountsResponse{
		Accounts: accounts,
		NewState: resp.NewState,
		HasMore:  resp.HasMore,
	}, nil
}

func (i *impl) FetchNextExternalAccounts(ctx context.Context, req *services.FetchNextExternalAccountsRequest) (*services.FetchNextExternalAccountsResponse, error) {
	i.logger.Info("fetching next external accounts...")

	resp, err := i.plugin.FetchNextExternalAccounts(ctx, models.FetchNextExternalAccountsRequest{
		FromPayload: req.FromPayload,
		State:       req.State,
		PageSize:    int(req.PageSize),
	})
	if err != nil {
		i.logger.Error("fetching next external accounts failed: ", err)
		return nil, translateErrorToGRPC(err)
	}

	externalAccounts := make([]*proto.Account, 0, len(resp.ExternalAccounts))
	for _, account := range resp.ExternalAccounts {
		externalAccounts = append(externalAccounts, grpc.TranslateAccount(account))
	}

	i.logger.Info("fetched next external accounts succeeded!")

	return &services.FetchNextExternalAccountsResponse{
		Accounts: externalAccounts,
		NewState: resp.NewState,
		HasMore:  resp.HasMore,
	}, nil
}

func (i *impl) FetchNextPayments(ctx context.Context, req *services.FetchNextPaymentsRequest) (*services.FetchNextPaymentsResponse, error) {
	i.logger.Info("fetching next payments...")

	resp, err := i.plugin.FetchNextPayments(ctx, models.FetchNextPaymentsRequest{
		FromPayload: req.FromPayload,
		State:       req.State,
		PageSize:    int(req.PageSize),
	})
	if err != nil {
		i.logger.Error("fetching next payments failed: ", err)
		return nil, translateErrorToGRPC(err)
	}

	payments := make([]*proto.Payment, 0, len(resp.Payments))
	for _, payment := range resp.Payments {
		payments = append(payments, grpc.TranslatePayment(payment))
	}

	i.logger.Info("fetched next payments succeeded!")

	return &services.FetchNextPaymentsResponse{
		Payments: payments,
		NewState: resp.NewState,
		HasMore:  resp.HasMore,
	}, nil
}

func (i *impl) FetchNextOthers(ctx context.Context, req *services.FetchNextOthersRequest) (*services.FetchNextOthersResponse, error) {
	i.logger.Info("fetching next others...")

	resp, err := i.plugin.FetchNextOthers(ctx, models.FetchNextOthersRequest{
		FromPayload: req.FromPayload,
		State:       req.State,
		PageSize:    int(req.PageSize),
		Name:        req.Name,
	})
	if err != nil {
		i.logger.Error("fetching next others failed: ", err)
		return nil, translateErrorToGRPC(err)
	}

	others := make([]*proto.Other, 0, len(resp.Others))
	for _, other := range resp.Others {
		others = append(others, &proto.Other{
			Id:    other.ID,
			Other: other.Other,
		})
	}

	i.logger.Info("fetched next others succeeded!")

	return &services.FetchNextOthersResponse{
		Others:   others,
		NewState: resp.NewState,
		HasMore:  resp.HasMore,
	}, nil
}

func (i *impl) CreateBankAccount(ctx context.Context, req *services.CreateBankAccountRequest) (*services.CreateBankAccountResponse, error) {
	i.logger.Info("creating bank account...")

	resp, err := i.plugin.CreateBankAccount(ctx, models.CreateBankAccountRequest{
		BankAccount: grpc.TranslateProtoBankAccount(req.BankAccount),
	})
	if err != nil {
		i.logger.Error("creating bank account failed: ", err)
		return nil, translateErrorToGRPC(err)
	}

	i.logger.Info("created bank account succeeded!")

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
