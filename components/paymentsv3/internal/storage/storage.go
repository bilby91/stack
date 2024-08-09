package storage

import (
	"context"

	"github.com/formancehq/paymentsv3/internal/models"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/uptrace/bun"
)

type Storage interface {
	// Connectors
	InstallConnector(ctx context.Context, c models.Connector) error
	UninstallConnector(ctx context.Context, id models.ConnectorID) error
	GetConnector(ctx context.Context, id models.ConnectorID) (*models.Connector, error)
	ListConnectors(ctx context.Context, q ListConnectorssQuery) (*bunpaginate.Cursor[models.Connector], error)

	// Workflow
	UpsertWorkflow(ctx context.Context, connectorID models.ConnectorID, tasks models.Workflow) error
	GetWorkflow(ctx context.Context, connectorID models.ConnectorID) (*models.Workflow, error)

	// Accounts
	UpsertAccounts(ctx context.Context, accounts []models.Account) error
	GetAccount(ctx context.Context, id models.AccountID) (*models.Account, error)
	ListAccounts(ctx context.Context, q ListAccountsQuery) (*bunpaginate.Cursor[models.Account], error)

	// Payments
	UpsertPayment(ctx context.Context, payments []models.Payment) error
	GetPayment(ctx context.Context, id models.PaymentID) (*models.Payment, error)
	ListPayments(ctx context.Context, q ListPaymentsQuery) (*bunpaginate.Cursor[models.Payment], error)

	// State
	UpsertState(ctx context.Context, state models.State) error
	GetState(ctx context.Context, id models.StateID) (models.State, error)
}

const encryptionOptions = "compress-algo=1, cipher-algo=aes256"

type store struct {
	db                  *bun.DB
	configEncryptionKey string
}

func newStorage(db *bun.DB, configEncryptionKey string) Storage {
	return &store{db: db, configEncryptionKey: configEncryptionKey}
}
