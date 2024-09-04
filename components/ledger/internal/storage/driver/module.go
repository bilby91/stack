package driver

import (
	"context"
	ledgercontroller "github.com/formancehq/ledger/internal/controller/ledger"
	"github.com/uptrace/bun"

	"github.com/formancehq/stack/libs/go-libs/logging"
	"go.uber.org/fx"
)

type PostgresConfig struct {
	ConnString string
}

type ModuleConfiguration struct {
}

func NewFXModule(autoUpgrade bool) fx.Option {
	return fx.Options(
		fx.Provide(func(db *bun.DB) (*Driver, error) {
			return New(db), nil
		}),
		fx.Provide(fx.Annotate(NewControllerStorageDriverAdapter, fx.As(new(ledgercontroller.StorageDriver)))),
		fx.Invoke(func(driver *Driver, lifecycle fx.Lifecycle, logger logging.Logger) error {
			lifecycle.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					logger.Infof("Initializing database...")
					return driver.Initialize(ctx)
				},
			})
			return nil
		}),
		fx.Invoke(func(lc fx.Lifecycle, driver *Driver) {
			if autoUpgrade {
				lc.Append(fx.Hook{
					OnStart: driver.UpgradeAllBuckets,
				})
			}
		}),
	)
}
