package storage

import (
	systemcontroller "github.com/formancehq/ledger/internal/controller/system"
	"github.com/formancehq/ledger/internal/storage/driver"
	"github.com/formancehq/ledger/internal/storage/system"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
)

func NewFXModule(autoUpgrade bool) fx.Option {
	return fx.Options(
		driver.NewFXModule(autoUpgrade),
		fx.Provide(func(db *bun.DB) *system.Store {
			return system.NewStore(db)
		}),
		fx.Provide(func(store *system.Store) systemcontroller.Store {
			return store
		}),
	)
}
