package sqlstorage

import (
	"database/sql"
	"github.com/numary/ledger/pkg/health"
	"github.com/numary/ledger/pkg/storage"
	"go.uber.org/fx"
)

type SQLiteConfig struct {
	Dir    string
	DBName string
}

type PostgresConfig struct {
	ConnString string
}

type ModuleConfig struct {
	StorageDriver  string
	SQLiteConfig   *SQLiteConfig
	PostgresConfig *PostgresConfig
}

func OpenSQLDB(flavor Flavor, dataSourceName string) (*sql.DB, error) {
	c, ok := sqlDrivers[flavor]
	if !ok {
		panic("PostgresSQL driver not found")
	}
	return sql.Open(c.driverName, dataSourceName)
}

func DriverModule(cfg ModuleConfig) fx.Option {
	options := make([]fx.Option, 0)
	options = append(options, fx.Provide(func() Flavor {
		return FlavorFromString(cfg.StorageDriver)
	}))

	switch FlavorFromString(cfg.StorageDriver) {
	case PostgreSQL:
		options = append(options, fx.Provide(func() (*sql.DB, error) {
			return OpenSQLDB(PostgreSQL, cfg.PostgresConfig.ConnString)
		}))
		options = append(options, fx.Provide(func(db *sql.DB) (*cachedDBDriver, error) {
			return NewCachedDBDriver(PostgreSQL.String(), PostgreSQL, func() (DB, error) {
				return NewPostgresDB(db), nil
			}), nil
		}))
		options = append(options, fx.Provide(func(driver *cachedDBDriver) storage.Driver {
			return driver
		}))
		options = append(options, health.ProvideHealthCheck(func(db *sql.DB) health.NamedCheck {
			return health.NewNamedCheck(PostgreSQL.String(), health.CheckFn(db.PingContext))
		}))
	case SQLite:
		options = append(options, fx.Provide(func() (storage.Driver, error) {
			return NewOpenCloseDBDriver(SQLite.String(), SQLite, func() (DB, error) {
				return NewSQLiteDB(cfg.SQLiteConfig.Dir, cfg.SQLiteConfig.DBName), nil
			}), nil
		}))
	default:
		panic("Unsupported driver: " + cfg.StorageDriver)
	}
	options = append(options, fx.Invoke(func(driver storage.Driver, lifecycle fx.Lifecycle) error {
		lifecycle.Append(fx.Hook{
			OnStart: driver.Initialize,
			OnStop:  driver.Close,
		})
		return nil
	}))
	return fx.Options(options...)
}
