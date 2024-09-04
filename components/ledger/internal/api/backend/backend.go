package backend

import (
	"context"
	ledgercontroller "github.com/formancehq/ledger/internal/controller/ledger"
	"github.com/formancehq/ledger/internal/controller/ledger/writer"
	systemcontroller "github.com/formancehq/ledger/internal/controller/system"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"

	ledger "github.com/formancehq/ledger/internal"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/formancehq/stack/libs/go-libs/migrations"
)

//go:generate mockgen -source backend.go -destination backend_generated.go -package backend . Ledger

type Ledger interface {
	GetAccountWithVolumes(ctx context.Context, query ledgercontroller.GetAccountQuery) (*ledger.ExpandedAccount, error)
	GetAccountsWithVolumes(ctx context.Context, query ledgercontroller.GetAccountsQuery) (*bunpaginate.Cursor[ledger.ExpandedAccount], error)
	CountAccounts(ctx context.Context, query ledgercontroller.GetAccountsQuery) (int, error)
	GetAggregatedBalances(ctx context.Context, q ledgercontroller.GetAggregatedBalanceQuery) (ledger.BalancesByAssets, error)
	GetMigrationsInfo(ctx context.Context) ([]migrations.Info, error)
	Stats(ctx context.Context) (ledgercontroller.Stats, error)
	GetLogs(ctx context.Context, query ledgercontroller.GetLogsQuery) (*bunpaginate.Cursor[ledger.ChainedLog], error)
	CountTransactions(ctx context.Context, query ledgercontroller.GetTransactionsQuery) (int, error)
	GetTransactions(ctx context.Context, query ledgercontroller.GetTransactionsQuery) (*bunpaginate.Cursor[ledger.ExpandedTransaction], error)
	GetTransactionWithVolumes(ctx context.Context, query ledgercontroller.GetTransactionQuery) (*ledger.ExpandedTransaction, error)

	CreateTransaction(ctx context.Context, parameters writer.Parameters, data ledger.RunScript) (*ledger.Transaction, error)
	RevertTransaction(ctx context.Context, parameters writer.Parameters, id int, force, atEffectiveDate bool) (*ledger.Transaction, error)
	SaveMeta(ctx context.Context, parameters writer.Parameters, targetType string, targetID any, m metadata.Metadata) error
	DeleteMetadata(ctx context.Context, parameters writer.Parameters, targetType string, targetID any, key string) error
	Import(ctx context.Context, stream chan *ledger.ChainedLog) error
	Export(ctx context.Context, w ledgercontroller.ExportWriter) error

	IsDatabaseUpToDate(ctx context.Context) (bool, error)

	GetVolumesWithBalances(ctx context.Context, q ledgercontroller.GetVolumesWithBalancesQuery) (*bunpaginate.Cursor[ledger.VolumesWithBalanceByAssetByAccount], error)
}

type Backend interface {
	GetLedgerController(ctx context.Context, name string) (Ledger, error)
	GetLedger(ctx context.Context, name string) (*ledger.Ledger, error)
	ListLedgers(ctx context.Context, query systemcontroller.ListLedgersQuery) (*bunpaginate.Cursor[ledger.Ledger], error)
	CreateLedger(ctx context.Context, name string, configuration ledger.Configuration) error
	UpdateLedgerMetadata(ctx context.Context, name string, m map[string]string) error
	GetVersion() string
	DeleteLedgerMetadata(ctx context.Context, param string, key string) error
}

type DefaultBackend struct {
	*systemcontroller.Controller
	version string
}

func (d *DefaultBackend) GetVersion() string {
	return d.version
}

func (d *DefaultBackend) GetLedgerController(ctx context.Context, name string) (Ledger, error) {
	return d.Controller.GetLedgerController(ctx, name)
}

var _ Backend = (*DefaultBackend)(nil)

func NewDefaultBackend(systemController *systemcontroller.Controller, version string) *DefaultBackend {
	return &DefaultBackend{
		Controller: systemController,
		version:    version,
	}
}
