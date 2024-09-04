package system

import (
	"context"
	ledger "github.com/formancehq/ledger/internal"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/formancehq/stack/libs/go-libs/metadata"
)

type Store interface {
	GetLedger(ctx context.Context, name string) (*ledger.Ledger, error)
	ListLedgers(ctx context.Context, query ListLedgersQuery) (*bunpaginate.Cursor[ledger.Ledger], error)
	UpdateLedgerMetadata(ctx context.Context, name string, m metadata.Metadata) error
	DeleteLedgerMetadata(ctx context.Context, param string, key string) error
}

type ListLedgersQuery bunpaginate.OffsetPaginatedQuery[PaginatedQueryOptions]

func (query ListLedgersQuery) WithPageSize(pageSize uint64) ListLedgersQuery {
	query.PageSize = pageSize
	return query
}

func NewListLedgersQuery(pageSize uint64) ListLedgersQuery {
	return ListLedgersQuery{
		PageSize: pageSize,
	}
}

type PaginatedQueryOptions struct {
	PageSize uint64 `json:"pageSize"`
}
