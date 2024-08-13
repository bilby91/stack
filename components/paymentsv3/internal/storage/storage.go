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

	// Tasks
	UpsertTasks(ctx context.Context, connectorID models.ConnectorID, tasks models.Tasks) error
	GetTasks(ctx context.Context, connectorID models.ConnectorID) (*models.Tasks, error)
	DeleteTasksFromConnectorID(ctx context.Context, connectorID models.ConnectorID) error

	// Workflows
	UpsertWorkflow(ctx context.Context, workflow models.Workflow) error
	GetWorflow(ctx context.Context, id string) (*models.Workflow, error)
	ListWorkflows(ctx context.Context, q ListWorkflowsQuery) (*bunpaginate.Cursor[models.Workflow], error)
	DeleteWorkflowsFromConnectorID(ctx context.Context, connectorID models.ConnectorID) error

	// Accounts
	UpsertAccounts(ctx context.Context, accounts []models.Account) error
	GetAccount(ctx context.Context, id models.AccountID) (*models.Account, error)
	ListAccounts(ctx context.Context, q ListAccountsQuery) (*bunpaginate.Cursor[models.Account], error)
	DeleteAccountsFromConnectorID(ctx context.Context, connectorID models.ConnectorID) error

	// Payments
	UpsertPayment(ctx context.Context, payments []models.Payment) error
	GetPayment(ctx context.Context, id models.PaymentID) (*models.Payment, error)
	ListPayments(ctx context.Context, q ListPaymentsQuery) (*bunpaginate.Cursor[models.Payment], error)
	DeletePaymentsFromConnectorID(ctx context.Context, connectorID models.ConnectorID) error

	// State
	UpsertState(ctx context.Context, state models.State) error
	GetState(ctx context.Context, id models.StateID) (models.State, error)
	DeleteStatesFromConnectorID(ctx context.Context, connectorID models.ConnectorID) error

	// Schedules
	UpsertSchedule(ctx context.Context, schedule models.Schedule) error
	ListSchedules(ctx context.Context, q ListSchedulesQuery) (*bunpaginate.Cursor[models.Schedule], error)
	DeleteSchedulesFromConnectorID(ctx context.Context, connectorID models.ConnectorID) error

	// Workflow Instances
	InsertNewInstance(ctx context.Context, instance models.Instance) error
	UpdateInstance(ctx context.Context, instance models.Instance) error
	ListInstances(ctx context.Context, q ListInstancesQuery) (*bunpaginate.Cursor[models.Instance], error)
	DeleteInstancesFromConnectorID(ctx context.Context, connectorID models.ConnectorID) error
}

const encryptionOptions = "compress-algo=1, cipher-algo=aes256"

type store struct {
	db                  *bun.DB
	configEncryptionKey string
}

func newStorage(db *bun.DB, configEncryptionKey string) Storage {
	return &store{db: db, configEncryptionKey: configEncryptionKey}
}
