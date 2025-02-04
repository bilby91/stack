package postgres

import (
	"github.com/uptrace/bun"

	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"

	"github.com/formancehq/webhooks/pkg/storage"
	"go.uber.org/fx"
)

func NewModule(connectionOptions bunconnect.ConnectionOptions, debug bool) fx.Option {
	return fx.Options(
		bunconnect.Module(connectionOptions, debug),
		fx.Provide(func(db *bun.DB) (storage.Store, error) {
			return NewStore(db)
		}),
	)
}
