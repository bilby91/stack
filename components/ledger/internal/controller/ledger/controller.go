package ledger

import (
	"context"
	"github.com/formancehq/ledger/internal/controller/ledger/writer"
	"github.com/formancehq/ledger/internal/opentelemetry/tracer"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	"sync"

	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"

	"github.com/ThreeDotsLabs/watermill/message"
	ledger "github.com/formancehq/ledger/internal"
)

type Controller struct {
	engine *writer.Writer
	store  Store
	mu     sync.Mutex
}

func New(
	name string,
	store Store,
	publisher message.Publisher,
	machineFactory writer.MachineFactory,
) *Controller {
	// TODO: restore
	//var monitor bus.Monitor = bus.NewNoOpMonitor()
	//if publisher != nil {
	//	monitor = bus.NewLedgerMonitor(publisher, store.Name())
	//}
	ret := &Controller{
		engine: writer.New(store, machineFactory),
		store:  store,
	}
	return ret
}

func (l *Controller) GetTransactions(ctx context.Context, q GetTransactionsQuery) (*bunpaginate.Cursor[ledger.ExpandedTransaction], error) {
	return tracer.Trace(ctx, "GetTransactions", func(ctx context.Context) (*bunpaginate.Cursor[ledger.ExpandedTransaction], error) {
		txs, err := l.store.GetTransactions(ctx, q)
		return txs, newStorageError(err, "getting transactions")
	})
}

func (l *Controller) CountTransactions(ctx context.Context, q GetTransactionsQuery) (int, error) {
	return tracer.Trace(ctx, "CountTransactions", func(ctx context.Context) (int, error) {
		count, err := l.store.CountTransactions(ctx, q)
		return count, newStorageError(err, "counting transactions")
	})
}

func (l *Controller) GetTransactionWithVolumes(ctx context.Context, query GetTransactionQuery) (*ledger.ExpandedTransaction, error) {
	return tracer.Trace(ctx, "GetTransactionWithVolumes", func(ctx context.Context) (*ledger.ExpandedTransaction, error) {
		tx, err := l.store.GetTransactionWithVolumes(ctx, query)
		return tx, newStorageError(err, "getting transaction")
	})
}

func (l *Controller) CountAccounts(ctx context.Context, a GetAccountsQuery) (int, error) {
	return tracer.Trace(ctx, "CountAccounts", func(ctx context.Context) (int, error) {
		count, err := l.store.CountAccounts(ctx, a)
		return count, newStorageError(err, "counting accounts")
	})
}

func (l *Controller) GetAccountsWithVolumes(ctx context.Context, a GetAccountsQuery) (*bunpaginate.Cursor[ledger.ExpandedAccount], error) {
	return tracer.Trace(ctx, "GetAccountsWithVolumes", func(ctx context.Context) (*bunpaginate.Cursor[ledger.ExpandedAccount], error) {
		accounts, err := l.store.GetAccountsWithVolumes(ctx, a)
		return accounts, newStorageError(err, "getting accounts")
	})
}

func (l *Controller) GetAccountWithVolumes(ctx context.Context, q GetAccountQuery) (*ledger.ExpandedAccount, error) {
	return tracer.Trace(ctx, "GetAccountWithVolumes", func(ctx context.Context) (*ledger.ExpandedAccount, error) {
		accounts, err := l.store.GetAccountWithVolumes(ctx, q)
		return accounts, newStorageError(err, "getting account")
	})
}

func (l *Controller) GetAggregatedBalances(ctx context.Context, q GetAggregatedBalanceQuery) (ledger.BalancesByAssets, error) {
	return tracer.Trace(ctx, "GetAggregatedBalances", func(ctx context.Context) (ledger.BalancesByAssets, error) {
		balances, err := l.store.GetAggregatedBalances(ctx, q)
		return balances, newStorageError(err, "getting balances aggregated")
	})
}

func (l *Controller) GetLogs(ctx context.Context, q GetLogsQuery) (*bunpaginate.Cursor[ledger.ChainedLog], error) {
	return tracer.Trace(ctx, "GetLogs", func(ctx context.Context) (*bunpaginate.Cursor[ledger.ChainedLog], error) {
		logs, err := l.store.GetLogs(ctx, q)
		return logs, newStorageError(err, "getting logs")
	})
}

func (l *Controller) markInUseIfNeeded(ctx context.Context) {
	//// todo: keep in memory to avoid repeating the same request again and again
	//if err := l.systemStore.UpdateLedgerState(ctx, l.store.Name(), system.StateInUse); err != nil {
	//	logging.FromContext(ctx).Error("Unable to declare ledger as in use")
	//	return
	//}
}

func (l *Controller) IsDatabaseUpToDate(ctx context.Context) (bool, error) {
	return tracer.Trace(ctx, "IsDatabaseUpToDate", func(ctx context.Context) (bool, error) {
		return l.store.IsUpToDate(ctx)
	})
}

func (l *Controller) CreateTransaction(ctx context.Context, parameters writer.Parameters, runScript ledger.RunScript) (*ledger.Transaction, error) {
	return tracer.Trace(ctx, "CreateTransaction", func(ctx context.Context) (*ledger.Transaction, error) {
		return l.engine.CreateTransaction(ctx, parameters, runScript)
	})
}

func (l *Controller) RevertTransaction(ctx context.Context, parameters writer.Parameters, id int, force, atEffectiveDate bool) (*ledger.Transaction, error) {
	return tracer.Trace(ctx, "RevertTransaction", func(ctx context.Context) (*ledger.Transaction, error) {
		return l.engine.RevertTransaction(ctx, parameters, id, force, atEffectiveDate)
	})
}

func (l *Controller) SaveMeta(ctx context.Context, parameters writer.Parameters, targetType string, targetID any, m metadata.Metadata) error {
	return tracer.SkipResult(tracer.Trace(ctx, "SaveMeta", tracer.NoResult(func(ctx context.Context) error {
		return l.engine.SaveMeta(ctx, parameters, targetType, targetID, m)
	})))
}

func (l *Controller) DeleteMetadata(ctx context.Context, parameters writer.Parameters, targetType string, targetID any, key string) error {
	return tracer.SkipResult(tracer.Trace(ctx, "DeleteMetadata", tracer.NoResult(func(ctx context.Context) error {
		return l.engine.DeleteMetadata(ctx, parameters, targetType, targetID, key)
	})))
}

func (l *Controller) GetVolumesWithBalances(ctx context.Context, q GetVolumesWithBalancesQuery) (*bunpaginate.Cursor[ledger.VolumesWithBalanceByAssetByAccount], error) {
	return tracer.Trace(ctx, "GetVolumesWithBalances", func(ctx context.Context) (*bunpaginate.Cursor[ledger.VolumesWithBalanceByAssetByAccount], error) {
		volumes, err := l.store.GetVolumesWithBalances(ctx, q)
		return volumes, newStorageError(err, "getting Volumes with balances")
	})
}
