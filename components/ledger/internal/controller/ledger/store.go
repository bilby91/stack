package ledger

import (
	"context"
	"encoding/json"
	ledger "github.com/formancehq/ledger/internal"
	"github.com/formancehq/ledger/internal/controller/ledger/writer"
	"github.com/formancehq/ledger/internal/machine/vm"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"github.com/formancehq/stack/libs/go-libs/query"
	"github.com/formancehq/stack/libs/go-libs/time"
)

type Store interface {
	writer.Store
	vm.Store
	// todo: move queries in controller package
	GetTransactions(ctx context.Context, q GetTransactionsQuery) (*bunpaginate.Cursor[ledger.ExpandedTransaction], error)
	CountTransactions(ctx context.Context, q GetTransactionsQuery) (int, error)
	GetTransactionWithVolumes(ctx context.Context, query GetTransactionQuery) (*ledger.ExpandedTransaction, error)
	CountAccounts(ctx context.Context, a GetAccountsQuery) (int, error)
	GetAccountsWithVolumes(ctx context.Context, a GetAccountsQuery) (*bunpaginate.Cursor[ledger.ExpandedAccount], error)
	GetAccountWithVolumes(ctx context.Context, q GetAccountQuery) (*ledger.ExpandedAccount, error)
	GetAggregatedBalances(ctx context.Context, q GetAggregatedBalanceQuery) (ledger.BalancesByAssets, error)
	GetLogs(ctx context.Context, q GetLogsQuery) (*bunpaginate.Cursor[ledger.ChainedLog], error)
	GetVolumesWithBalances(ctx context.Context, q GetVolumesWithBalancesQuery) (*bunpaginate.Cursor[ledger.VolumesWithBalanceByAssetByAccount], error)
	IsUpToDate(ctx context.Context) (bool, error)
}

type StorageDriver interface {
	OpenLedger(context.Context, string) (Store, error)
	CreateLedger(context.Context, string, ledger.Configuration) error
}

type GetTransactionsQuery bunpaginate.ColumnPaginatedQuery[PaginatedQueryOptions[PITFilterWithVolumes]]

func (q GetTransactionsQuery) WithExpandVolumes() GetTransactionsQuery {
	q.Options.Options.ExpandVolumes = true

	return q
}

func (q GetTransactionsQuery) WithExpandEffectiveVolumes() GetTransactionsQuery {
	q.Options.Options.ExpandEffectiveVolumes = true

	return q
}

func (q GetTransactionsQuery) WithColumn(column string) GetTransactionsQuery {
	ret := pointer.For((bunpaginate.ColumnPaginatedQuery[PaginatedQueryOptions[PITFilterWithVolumes]])(q))
	ret = ret.WithColumn(column)

	return GetTransactionsQuery(*ret)
}

func NewGetTransactionsQuery(options PaginatedQueryOptions[PITFilterWithVolumes]) GetTransactionsQuery {
	return GetTransactionsQuery{
		PageSize: options.PageSize,
		Column:   "id",
		Order:    bunpaginate.OrderDesc,
		Options:  options,
	}
}

type GetTransactionQuery struct {
	PITFilterWithVolumes
	ID int
}

func (q GetTransactionQuery) WithExpandVolumes() GetTransactionQuery {
	q.ExpandVolumes = true

	return q
}

func (q GetTransactionQuery) WithExpandEffectiveVolumes() GetTransactionQuery {
	q.ExpandEffectiveVolumes = true

	return q
}

func NewGetTransactionQuery(id int) GetTransactionQuery {
	return GetTransactionQuery{
		PITFilterWithVolumes: PITFilterWithVolumes{},
		ID:                   id,
	}
}

type GetAccountsQuery bunpaginate.OffsetPaginatedQuery[PaginatedQueryOptions[PITFilterWithVolumes]]

func (q GetAccountsQuery) WithExpandVolumes() GetAccountsQuery {
	q.Options.Options.ExpandVolumes = true

	return q
}

func (q GetAccountsQuery) WithExpandEffectiveVolumes() GetAccountsQuery {
	q.Options.Options.ExpandEffectiveVolumes = true

	return q
}

func NewGetAccountsQuery(opts PaginatedQueryOptions[PITFilterWithVolumes]) GetAccountsQuery {
	return GetAccountsQuery{
		PageSize: opts.PageSize,
		Order:    bunpaginate.OrderAsc,
		Options:  opts,
	}
}

type GetAccountQuery struct {
	PITFilterWithVolumes
	Addr string
}

func (q GetAccountQuery) WithPIT(pit time.Time) GetAccountQuery {
	q.PIT = &pit

	return q
}

func (q GetAccountQuery) WithExpandVolumes() GetAccountQuery {
	q.ExpandVolumes = true

	return q
}

func (q GetAccountQuery) WithExpandEffectiveVolumes() GetAccountQuery {
	q.ExpandEffectiveVolumes = true

	return q
}

func NewGetAccountQuery(addr string) GetAccountQuery {
	return GetAccountQuery{
		Addr: addr,
	}
}

type GetAggregatedBalanceQuery struct {
	PITFilter
	QueryBuilder     query.Builder
	UseInsertionDate bool
}

func NewGetAggregatedBalancesQuery(filter PITFilter, qb query.Builder, useInsertionDate bool) GetAggregatedBalanceQuery {
	return GetAggregatedBalanceQuery{
		PITFilter:        filter,
		QueryBuilder:     qb,
		UseInsertionDate: useInsertionDate,
	}
}

type GetVolumesWithBalancesQuery bunpaginate.OffsetPaginatedQuery[PaginatedQueryOptions[FiltersForVolumes]]

func NewGetVolumesWithBalancesQuery(opts PaginatedQueryOptions[FiltersForVolumes]) GetVolumesWithBalancesQuery {
	return GetVolumesWithBalancesQuery{
		PageSize: opts.PageSize,
		Order:    bunpaginate.OrderAsc,
		Options:  opts,
	}
}

type GetLogsQuery bunpaginate.ColumnPaginatedQuery[PaginatedQueryOptions[any]]

func (q GetLogsQuery) WithOrder(order bunpaginate.Order) GetLogsQuery {
	q.Order = order
	return q
}

func NewGetLogsQuery(options PaginatedQueryOptions[any]) GetLogsQuery {
	return GetLogsQuery{
		PageSize: options.PageSize,
		Column:   "id",
		Order:    bunpaginate.OrderDesc,
		Options:  options,
	}
}

type PaginatedQueryOptions[T any] struct {
	QueryBuilder query.Builder `json:"qb"`
	PageSize     uint64        `json:"pageSize"`
	Options      T             `json:"options"`
}

func (v *PaginatedQueryOptions[T]) UnmarshalJSON(data []byte) error {
	type aux struct {
		QueryBuilder json.RawMessage `json:"qb"`
		PageSize     uint64          `json:"pageSize"`
		Options      T               `json:"options"`
	}
	x := &aux{}
	if err := json.Unmarshal(data, x); err != nil {
		return err
	}

	*v = PaginatedQueryOptions[T]{
		PageSize: x.PageSize,
		Options:  x.Options,
	}

	var err error
	if x.QueryBuilder != nil {
		v.QueryBuilder, err = query.ParseJSON(string(x.QueryBuilder))
		if err != nil {
			return err
		}
	}

	return nil
}

func (opts PaginatedQueryOptions[T]) WithQueryBuilder(qb query.Builder) PaginatedQueryOptions[T] {
	opts.QueryBuilder = qb

	return opts
}

func (opts PaginatedQueryOptions[T]) WithPageSize(pageSize uint64) PaginatedQueryOptions[T] {
	opts.PageSize = pageSize

	return opts
}

func NewPaginatedQueryOptions[T any](options T) PaginatedQueryOptions[T] {
	return PaginatedQueryOptions[T]{
		Options:  options,
		PageSize: bunpaginate.QueryDefaultPageSize,
	}
}

type PITFilter struct {
	PIT *time.Time `json:"pit"`
	OOT *time.Time `json:"oot"`
}

type PITFilterWithVolumes struct {
	PITFilter
	ExpandVolumes          bool `json:"volumes"`
	ExpandEffectiveVolumes bool `json:"effectiveVolumes"`
}

type FiltersForVolumes struct {
	PITFilter
	UseInsertionDate bool
	GroupLvl         uint
}
