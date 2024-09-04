package ledger

import (
	"context"
	"database/sql"
	"github.com/formancehq/ledger/internal/controller/ledger"
	"github.com/formancehq/ledger/internal/controller/ledger/writer"
	"github.com/uptrace/bun"
)

type TX struct {
	*Store
	sqlTX bun.Tx
}

func (t *TX) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return t.sqlTX.QueryContext(ctx, query, args...)
}

func (t *TX) Commit(_ context.Context) error {
	return t.sqlTX.Commit()
}

func (t *TX) Rollback(_ context.Context) error {
	return t.sqlTX.Rollback()
}

type DefaultStoreAdapter struct {
	*Store
}

func (d *DefaultStoreAdapter) BeginTX(ctx context.Context) (writer.TX, error) {
	tx, err := d.GetDB().BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}

	return &TX{
		Store: d.Store.WithDB(tx),
		sqlTX: tx,
	}, nil
}

func NewDefaultStoreAdapter(store *Store) *DefaultStoreAdapter {
	return &DefaultStoreAdapter{
		Store: store,
	}
}

var _ ledger.Store = (*DefaultStoreAdapter)(nil)
