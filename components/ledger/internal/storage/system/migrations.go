package system

import (
	"context"
	"github.com/formancehq/stack/libs/go-libs/platform/postgres"

	"github.com/formancehq/stack/libs/go-libs/migrations"
	"github.com/uptrace/bun"
)

func Migrate(ctx context.Context, db bun.IDB) error {
	migrator := migrations.NewMigrator(migrations.WithSchema(Schema, true))
	migrator.RegisterMigrations(
		migrations.Migration{
			Name: "Init schema",
			UpWithContext: func(ctx context.Context, tx bun.Tx) error {

				_, err := tx.NewCreateTable().
					Model((*Ledger)(nil)).
					Exec(ctx)
				if err != nil {
					return postgres.ResolveError(err)
				}

				_, err = tx.NewCreateTable().
					Model((*configuration)(nil)).
					Exec(ctx)
				return postgres.ResolveError(err)
			},
		},
		migrations.Migration{
			Name: "Add ledger, bucket naming constraints 63 chars",
			UpWithContext: func(ctx context.Context, tx bun.Tx) error {
				_, err := tx.ExecContext(ctx, `
					alter table ledgers
					add column if not exists ledger varchar(63),
					add column if not exists bucket varchar(63);

					alter table ledgers
					alter column ledger type varchar(63),
					alter column bucket type varchar(63);
				`)
				if err != nil {
					return err
				}
				return nil
			},
		},
		migrations.Migration{
			Name: "Add ledger metadata",
			UpWithContext: func(ctx context.Context, tx bun.Tx) error {
				_, err := tx.ExecContext(ctx, `
					alter table ledgers
					add column if not exists metadata jsonb;
				`)
				if err != nil {
					return err
				}
				return nil
			},
		},
		migrations.Migration{
			Name: "Fix empty ledger metadata",
			UpWithContext: func(ctx context.Context, tx bun.Tx) error {
				_, err := tx.ExecContext(ctx, `
					update ledgers
					set metadata = '{}'::jsonb
					where metadata is null;
				`)
				if err != nil {
					return err
				}
				return nil
			},
		},
		migrations.Migration{
			Name: "Add ledger state",
			UpWithContext: func(ctx context.Context, tx bun.Tx) error {
				_, err := tx.ExecContext(ctx, `
					alter table ledgers
					add column if not exists state varchar(255) default 'initializing';

					update ledgers
					set state = 'in-use'
					where state = '';
				`)
				if err != nil {
					return err
				}
				return nil
			},
		},

	)
	return migrator.Up(ctx, db)
}
