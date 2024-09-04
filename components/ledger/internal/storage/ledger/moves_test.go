//go:build it

package ledger_test

import (
	ledger "github.com/formancehq/ledger/internal"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/formancehq/stack/libs/go-libs/time"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"math/big"
	"testing"
)

func TestMoves(t *testing.T) {
	t.Parallel()

	store := newLedgerStore(t)
	ctx := logging.TestingContext()

	now := time.Now()
	_, err := store.UpsertAccount(ctx, ledger.Account{
		BaseModel:     bun.BaseModel{},
		Address:       "world",
		Metadata:      metadata.Metadata{},
		FirstUsage:    now,
		InsertionDate: now,
		UpdatedAt:     now,
	})
	require.NoError(t, err)

	_, err = store.UpsertAccount(ctx, ledger.Account{
		BaseModel:     bun.BaseModel{},
		Address:       "bank",
		Metadata:      metadata.Metadata{},
		FirstUsage:    now,
		InsertionDate: now,
		UpdatedAt:     now,
	})
	require.NoError(t, err)

	_, err = store.UpsertAccount(ctx, ledger.Account{
		BaseModel:     bun.BaseModel{},
		Address:       "bank2",
		Metadata:      metadata.Metadata{},
		FirstUsage:    now,
		InsertionDate: now,
		UpdatedAt:     now,
	})
	require.NoError(t, err)

	// Insert first tx
	tx1, err := store.InsertTransaction(ctx, ledger.NewTransactionData().WithPostings(
		ledger.NewPosting("world", "bank", "USD/2", big.NewInt(100)),
	).WithDate(now))
	require.NoError(t, err)

	for _, move := range tx1.GetMoves() {
		require.NoError(t, store.InsertMoves(ctx, move))
	}

	balance, err := store.GetBalance(ctx, "world", "USD/2")
	require.NoError(t, err)
	require.Equal(t, big.NewInt(-100), balance)

	balance, err = store.GetBalance(ctx, "bank", "USD/2")
	require.NoError(t, err)
	require.Equal(t, big.NewInt(100), balance)

	// Insert second tx
	tx2, err := store.InsertTransaction(ctx, ledger.NewTransactionData().WithPostings(
		ledger.NewPosting("world", "bank2", "USD/2", big.NewInt(100)),
	).WithDate(now.Add(time.Minute)))
	require.NoError(t, err)

	for _, move := range tx2.GetMoves() {
		require.NoError(t, store.InsertMoves(ctx, move))
	}

	balance, err = store.GetBalance(ctx, "world", "USD/2")
	require.NoError(t, err)
	require.Equal(t, big.NewInt(-200), balance)

	balance, err = store.GetBalance(ctx, "bank2", "USD/2")
	require.NoError(t, err)
	require.Equal(t, big.NewInt(100), balance)

	// Insert backdated tx
	tx3, err := store.InsertTransaction(ctx, ledger.NewTransactionData().WithPostings(
		ledger.NewPosting("world", "bank", "USD/2", big.NewInt(100)),
	).WithDate(now.Add(30*time.Second)))
	require.NoError(t, err)

	for _, move := range tx3.GetMoves() {
		require.NoError(t, store.InsertMoves(ctx, move))
	}

	balance, err = store.GetBalance(ctx, "world", "USD/2")
	require.NoError(t, err)
	require.Equal(t, big.NewInt(-300), balance)

	balance, err = store.GetBalance(ctx, "bank", "USD/2")
	require.NoError(t, err)
	require.Equal(t, big.NewInt(200), balance)

	//utils.DumpTables(t, ctx, store.GetDB(),
	//	"select * from "+store.PrefixWithBucket("accounts"),
	//	"select * from "+store.PrefixWithBucket("transactions"),
	//	"select * from "+store.PrefixWithBucket("moves")+" order by effective_date, seq",
	//	//"select * from "+store.PrefixWithBucket("moves")+" order by seq",
	//)
}
