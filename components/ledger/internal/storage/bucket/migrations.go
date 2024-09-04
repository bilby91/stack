package bucket

import (
	"bytes"
	"context"
	"embed"
	_ "embed"
	"github.com/formancehq/ledger/internal/opentelemetry/tracer"
	"github.com/formancehq/stack/libs/go-libs/migrations"
	"github.com/uptrace/bun"
	"text/template"
)

//go:embed migrations
var migrationsDir embed.FS

func getMigrator(name string) *migrations.Migrator {
	migrator := migrations.NewMigrator(migrations.WithSchema(name, true))
	migrator.RegisterMigrationsFromFileSystem(migrationsDir, "migrations", func(s string) string {
		buf := bytes.NewBufferString("")

		t := template.Must(template.New("migration").Parse(s))
		if err := t.Execute(buf, map[string]interface{}{
			"Bucket": name,
		}); err != nil {
			panic(err)
		}

		return buf.String()
	})

	return migrator
}

func Migrate(ctx context.Context, db bun.IDB, name string) error {
	ctx, span := tracer.Start(ctx, "Migrate bucket")
	defer span.End()

	return getMigrator(name).Up(ctx, db)
}
