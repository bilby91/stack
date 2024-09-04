package ledger

import (
	"context"

	"github.com/formancehq/stack/libs/go-libs/migrations"
)

func (l *Controller) GetMigrationsInfo(ctx context.Context) ([]migrations.Info, error) {
	panic("not implemented")
	//return l.store.GetMigrationsInfo(ctx)
}
