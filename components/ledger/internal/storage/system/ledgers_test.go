//go:build it

package system

import (
	"fmt"
	ledger "github.com/formancehq/ledger/internal"
	systemcontroller "github.com/formancehq/ledger/internal/controller/system"
	"github.com/formancehq/stack/libs/go-libs/bun/bundebug"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/uptrace/bun"
	"testing"

	"github.com/google/uuid"

	"github.com/formancehq/stack/libs/go-libs/time"

	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"

	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/stretchr/testify/require"
)

func newSystemStore(t *testing.T) *Store {
	t.Helper()

	ctx := logging.TestingContext()

	hooks := make([]bun.QueryHook, 0)
	if testing.Verbose() {
		hooks = append(hooks, bundebug.NewQueryHook())
	}

	pgServer := srv.NewDatabase(t)
	db, err := bunconnect.OpenSQLDB(ctx, pgServer.ConnectionOptions(), hooks...)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, db.Close())
	})

	store := NewStore(db)

	require.NoError(t, Migrate(ctx, store.DB()))

	return store
}

func TestListLedgers(t *testing.T) {
	ctx := logging.TestingContext()
	store := newSystemStore(t)

	ledgers := make([]ledger.Ledger, 0)
	pageSize := uint64(2)
	count := uint64(10)
	now := time.Now()
	for i := uint64(0); i < count; i++ {
		m := metadata.Metadata{}
		if i%2 == 0 {
			m["foo"] = "bar"
		}
		ledger := ledger.Ledger{
			Name:    fmt.Sprintf("ledger%d", i),
			AddedAt: now.Add(time.Duration(i) * time.Second),
			Configuration: ledger.Configuration{
				Metadata: m,
			},
		}
		ledgers = append(ledgers, ledger)
		_, err := store.RegisterLedger(ctx, &ledger)
		require.NoError(t, err)
	}

	cursor, err := store.ListLedgers(ctx, systemcontroller.NewListLedgersQuery(pageSize))
	require.NoError(t, err)
	require.Len(t, cursor.Data, int(pageSize))
	require.Equal(t, ledgers[:pageSize], cursor.Data)

	for i := pageSize; i < count; i += pageSize {
		query := systemcontroller.ListLedgersQuery{}
		require.NoError(t, bunpaginate.UnmarshalCursor(cursor.Next, &query))

		cursor, err = store.ListLedgers(ctx, query)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 2)
		require.Equal(t, ledgers[i:i+pageSize], cursor.Data)
	}
}

func TestUpdateLedgerMetadata(t *testing.T) {
	ctx := logging.TestingContext()
	store := newSystemStore(t)

	ledger := &ledger.Ledger{
		Name:    uuid.NewString(),
		AddedAt: time.Now(),
	}
	_, err := store.RegisterLedger(ctx, ledger)
	require.NoError(t, err)

	addedMetadata := metadata.Metadata{
		"foo": "bar",
	}
	err = store.UpdateLedgerMetadata(ctx, ledger.Name, addedMetadata)
	require.NoError(t, err)

	ledgerFromDB, err := store.GetLedger(ctx, ledger.Name)
	require.NoError(t, err)
	require.Equal(t, addedMetadata, ledgerFromDB.Metadata)
}

func TestDeleteLedgerMetadata(t *testing.T) {
	ctx := logging.TestingContext()
	store := newSystemStore(t)

	ledger := &ledger.Ledger{
		Name:    uuid.NewString(),
		AddedAt: time.Now(),
		Configuration: ledger.Configuration{
			Metadata: map[string]string{
				"foo": "bar",
			},
		},
	}
	_, err := store.RegisterLedger(ctx, ledger)
	require.NoError(t, err)

	err = store.DeleteLedgerMetadata(ctx, ledger.Name, "foo")
	require.NoError(t, err)

	ledgerFromDB, err := store.GetLedger(ctx, ledger.Name)
	require.NoError(t, err)
	require.Equal(t, metadata.Metadata{}, ledgerFromDB.Metadata)
}
