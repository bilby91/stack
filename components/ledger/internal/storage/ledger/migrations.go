package ledger

import (
	"bytes"
	"context"
	"embed"
	_ "embed"
	"fmt"
	"github.com/formancehq/ledger/internal/opentelemetry/tracer"
	"github.com/formancehq/stack/libs/go-libs/migrations"
	"github.com/uptrace/bun"
	"text/template"
)

//go:embed migrations
var migrationsDir embed.FS

func getMigrator(bucketName, ledgerName string) *migrations.Migrator {
	migrator := migrations.NewMigrator(
		migrations.WithSchema(bucketName, false),
		migrations.WithTableName(fmt.Sprintf("migrations_%s", ledgerName)),
	)
	migrator.RegisterMigrationsFromFileSystem(migrationsDir, "migrations", func(s string) string {
		buf := bytes.NewBufferString("")

		t := template.Must(template.New("migration").Parse(s))
		if err := t.Execute(buf, map[string]interface{}{
			"Bucket": bucketName,
			"Ledger": ledgerName,
		}); err != nil {
			panic(err)
		}

		return buf.String()
	})

	return migrator
}

func Migrate(ctx context.Context, db bun.IDB, bucketName, ledgerName string) error {
	ctx, span := tracer.Start(ctx, "Migrate ledger")
	defer span.End()

	return getMigrator(bucketName, ledgerName).Up(ctx, db)
}
