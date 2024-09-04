//go:build it

package driver_test

import (
	"fmt"
	"github.com/formancehq/ledger/internal/storage/driver"
	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"
	"github.com/formancehq/stack/libs/go-libs/bun/bundebug"
	"github.com/formancehq/stack/libs/go-libs/testing/docker"
	"github.com/formancehq/stack/libs/go-libs/testing/platform/pgtesting"
	"github.com/uptrace/bun"
	"testing"

	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/stretchr/testify/require"
)

// todo: restore
//func TestErrorOnOutdatedBucket(t *testing.T) {
//	t.Parallel()
//
//	ctx := logging.TestingContext()
//	d := newStorageDriver(t)
//
//	name := uuid.NewString()
//
//	b, err := d.OpenBucket(name)
//	require.NoError(t, err)
//
//	upToDate, err := b.IsUpToDate(ctx)
//	require.NoError(t, err)
//	require.False(t, upToDate)
//}

// todo: restore
//func TestGetLedgerFromAlternateBucket(t *testing.T) {
//	t.Parallel()
//
//	d := newStorageDriver(t)
//	ctx := logging.TestingContext()
//
//	ledgerName := "ledger0"
//	bucketName := "bucket0"
//
//	bucket, err := d.CreateBucket(ctx, bucketName)
//	require.NoError(t, err)
//
//	_, err = bucket.GetLedgerStore(ctx, ledgerName)
//	require.NoError(t, err)
//}

func TestUpgradeAllBuckets(t *testing.T) {
	t.Parallel()

	d := newStorageDriver(t)
	ctx := logging.TestingContext()

	count := 30

	for i := 0; i < count; i++ {
		name := fmt.Sprintf("ledger%d", i)
		_, err := d.CreateBucket(ctx, name)
		require.NoError(t, err)
	}

	require.NoError(t, d.UpgradeAllBuckets(ctx))
}

func newStorageDriver(t docker.T) *driver.Driver {
	t.Helper()

	ctx := logging.TestingContext()
	pgServer := pgtesting.CreatePostgresServer(t, docker.NewPool(t, logging.Testing()))
	pgDatabase := pgServer.NewDatabase(t)

	hooks := make([]bun.QueryHook, 0)
	if testing.Verbose() {
		hooks = append(hooks, bundebug.NewQueryHook())
	}
	db, err := bunconnect.OpenSQLDB(ctx, pgDatabase.ConnectionOptions(), hooks...)
	require.NoError(t, err)

	d := driver.New(db)

	require.NoError(t, d.Initialize(logging.TestingContext()))

	return d
}
