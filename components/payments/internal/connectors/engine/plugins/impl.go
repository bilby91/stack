package plugins

import (
	"context"

	"github.com/formancehq/payments/internal/connectors/grpc"
	"github.com/formancehq/payments/internal/connectors/grpc/proto/services"
	"github.com/formancehq/payments/internal/models"
)

type impl struct {
	pluginClient grpc.PSP
}

func (i *impl) Install(ctx context.Context, req models.InstallRequest) (models.InstallResponse, error) {
	resp, err := i.pluginClient.Install(ctx, &services.InstallRequest{
		Config: req.Config,
	})
	if err != nil {
		return models.InstallResponse{}, err
	}

	capabilities := make([]models.Capability, 0, len(resp.Capabilities))
	for _, capability := range resp.Capabilities {
		capabilities = append(capabilities, models.Capability(capability))
	}

	return models.InstallResponse{
		Capabilities: capabilities,
		Workflow:     grpc.TranslateProtoWorkflow(resp.Workflow),
	}, nil
}

func (i *impl) FetchNextAccounts(ctx context.Context, req models.FetchNextAccountsRequest) (models.FetchNextAccountsResponse, error) {
	resp, err := i.pluginClient.FetchNextAccounts(ctx, &services.FetchNextAccountsRequest{
		FromPayload: req.FromPayload,
		State:       req.State,
		PageSize:    int64(req.PageSize),
	})
	if err != nil {
		return models.FetchNextAccountsResponse{}, err
	}

	accounts := make([]models.PSPAccount, 0, len(resp.Accounts))
	for _, account := range resp.Accounts {
		accounts = append(accounts, grpc.TranslateProtoAccount(account))
	}

	return models.FetchNextAccountsResponse{
		Accounts: accounts,
		NewState: resp.NewState,
		HasMore:  resp.HasMore,
	}, nil
}

func (i *impl) FetchNextExternalAccounts(ctx context.Context, req models.FetchNextExternalAccountsRequest) (models.FetchNextExternalAccountsResponse, error) {
	resp, err := i.pluginClient.FetchNextExternalAccounts(ctx, &services.FetchNextExternalAccountsRequest{
		FromPayload: req.FromPayload,
		State:       req.State,
		PageSize:    int64(req.PageSize),
	})
	if err != nil {
		return models.FetchNextExternalAccountsResponse{}, err
	}

	externalAccounts := make([]models.PSPAccount, 0, len(resp.Accounts))
	for _, account := range resp.Accounts {
		externalAccounts = append(externalAccounts, grpc.TranslateProtoAccount(account))
	}

	return models.FetchNextExternalAccountsResponse{
		ExternalAccounts: externalAccounts,
		NewState:         resp.NewState,
		HasMore:          resp.HasMore,
	}, nil
}

func (i *impl) FetchNextPayments(ctx context.Context, req models.FetchNextPaymentsRequest) (models.FetchNextPaymentsResponse, error) {
	resp, err := i.pluginClient.FetchNextPayments(ctx, &services.FetchNextPaymentsRequest{
		FromPayload: req.FromPayload,
		State:       req.State,
		PageSize:    int64(req.PageSize),
	})
	if err != nil {
		return models.FetchNextPaymentsResponse{}, err
	}

	payments := make([]models.PSPPayment, 0, len(resp.Payments))
	for _, payment := range resp.Payments {
		p, err := grpc.TranslateProtoPayment(payment)
		if err != nil {
			return models.FetchNextPaymentsResponse{}, err
		}
		payments = append(payments, p)
	}

	return models.FetchNextPaymentsResponse{
		Payments: payments,
		NewState: resp.NewState,
		HasMore:  resp.HasMore,
	}, nil
}

func (i *impl) FetchNextOthers(ctx context.Context, req models.FetchNextOthersRequest) (models.FetchNextOthersResponse, error) {
	resp, err := i.pluginClient.FetchNextOthers(ctx, &services.FetchNextOthersRequest{
		Name:        req.Name,
		FromPayload: req.FromPayload,
		State:       req.State,
		PageSize:    int64(req.PageSize),
	})
	if err != nil {
		return models.FetchNextOthersResponse{}, err
	}

	others := make([]models.PSPOther, 0, len(resp.Others))
	for _, other := range resp.Others {
		others = append(others, models.PSPOther{
			ID:    other.Id,
			Other: other.Other,
		})
	}

	return models.FetchNextOthersResponse{
		Others:   others,
		NewState: resp.NewState,
		HasMore:  resp.HasMore,
	}, nil
}

func (i *impl) CreateBankAccount(ctx context.Context, req models.CreateBankAccountRequest) (models.CreateBankAccountResponse, error) {
	resp, err := i.pluginClient.CreateBankAccount(ctx, &services.CreateBankAccountRequest{
		BankAccount: grpc.TranslateBankAccount(req.BankAccount),
	})
	if err != nil {
		return models.CreateBankAccountResponse{}, err
	}

	return models.CreateBankAccountResponse{
		RelatedAccount: grpc.TranslateProtoAccount(resp.RelatedAccount),
	}, nil
}

var _ models.Plugin = &impl{}
