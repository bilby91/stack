package backend

import (
	"context"

	"github.com/formancehq/payments/cmd/api/internal/storage"
	"github.com/formancehq/payments/internal/models"
	"github.com/google/uuid"
)

//go:generate mockgen -source backend.go -destination backend_generated.go -package backend . Service
type Service interface {
	Ping() error
	ListAccounts(ctx context.Context, pagination storage.PaginatorQuery) ([]*models.Account, storage.PaginationDetails, error)
	GetAccount(ctx context.Context, id string) (*models.Account, error)
	ListBalances(ctx context.Context, query storage.BalanceQuery) ([]*models.Balance, storage.PaginationDetails, error)
	ListBankAccounts(ctx context.Context, pagination storage.PaginatorQuery) ([]*models.BankAccount, storage.PaginationDetails, error)
	GetBankAccount(ctx context.Context, id uuid.UUID, expand bool) (*models.BankAccount, error)
	ListTransferInitiations(ctx context.Context, pagination storage.PaginatorQuery) ([]*models.TransferInitiation, storage.PaginationDetails, error)
	ReadTransferInitiation(ctx context.Context, id models.TransferInitiationID) (*models.TransferInitiation, error)
	ListPayments(ctx context.Context, pagination storage.PaginatorQuery) ([]*models.Payment, storage.PaginationDetails, error)
	GetPayment(ctx context.Context, id string) (*models.Payment, error)
	UpdatePaymentMetadata(ctx context.Context, paymentID models.PaymentID, metadata map[string]string) error
}

type Backend interface {
	GetService() Service
}

type DefaultBackend struct {
	service Service
}

func (d DefaultBackend) GetService() Service {
	return d.service
}

func NewDefaultBackend(service Service) Backend {
	return &DefaultBackend{
		service: service,
	}
}
