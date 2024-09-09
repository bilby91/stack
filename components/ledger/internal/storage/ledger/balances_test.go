//go:build it

package ledger_test

import (
	ledgercontroller "github.com/formancehq/ledger/internal/controller/ledger"
	"math/big"
	"testing"

	"github.com/formancehq/stack/libs/go-libs/time"

	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/pointer"

	ledger "github.com/formancehq/ledger/internal"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/formancehq/stack/libs/go-libs/query"
	"github.com/stretchr/testify/require"
)

func TestGetBalancesAggregated(t *testing.T) {
	t.Parallel()
	store := newLedgerStore(t)
	now := time.Now()
	ctx := logging.TestingContext()

	bigInt, _ := big.NewInt(0).SetString("999999999999999999999999999999999999999999999999999999999999999999999999999999999", 10)
	smallInt := big.NewInt(199)

	tx1 := ledger.NewTransactionData().
		WithPostings(
			ledger.NewPosting("world", "users:1", "USD", bigInt),
			ledger.NewPosting("world", "users:2", "USD", smallInt),
		).
		WithDate(now).
		WithInsertedAt(now)
	_, err := store.InsertTransaction(ctx, tx1)
	require.NoError(t, err)

	tx2 := ledger.NewTransactionData().
		WithPostings(
			ledger.NewPosting("world", "users:1", "USD", bigInt),
			ledger.NewPosting("world", "users:2", "USD", smallInt),
			ledger.NewPosting("world", "xxx", "EUR", smallInt),
		).
		WithDate(now.Add(-time.Minute)).
		WithInsertedAt(now.Add(time.Minute))
	_, err = store.InsertTransaction(ctx, tx2)
	require.NoError(t, err)

	require.NoError(t, store.UpdateAccountMetadata(ctx, "users:1", metadata.Metadata{
		"category": "premium",
	}))

	require.NoError(t, store.UpdateAccountMetadata(ctx, "users:2", metadata.Metadata{
		"category": "premium",
	}))

	require.NoError(t, store.DeleteAccountMetadata(ctx, "users:2", "category"))

	require.NoError(t, store.UpdateAccountMetadata(ctx, "users:1", metadata.Metadata{
		"category": "premium",
	}))

	require.NoError(t, store.UpdateAccountMetadata(ctx, "users:2", metadata.Metadata{
		"category": "2",
	}))

	require.NoError(t, store.UpdateAccountMetadata(ctx, "world", metadata.Metadata{
		"world": "bar",
	}))

	t.Run("aggregate on all", func(t *testing.T) {
		t.Parallel()
		cursor, err := store.GetAggregatedBalances(ctx, ledgercontroller.NewGetAggregatedBalancesQuery(ledgercontroller.PITFilter{}, nil, false))
		require.NoError(t, err)
		RequireEqual(t, ledger.BalancesByAssets{
			"USD": big.NewInt(0),
			"EUR": big.NewInt(0),
		}, cursor)
	})
	t.Run("filter on address", func(t *testing.T) {
		t.Parallel()
		ret, err := store.GetAggregatedBalances(ctx, ledgercontroller.NewGetAggregatedBalancesQuery(ledgercontroller.PITFilter{},
			query.Match("address", "users:"), false))
		require.NoError(t, err)
		require.Equal(t, ledger.BalancesByAssets{
			"USD": big.NewInt(0).Add(
				big.NewInt(0).Mul(bigInt, big.NewInt(2)),
				big.NewInt(0).Mul(smallInt, big.NewInt(2)),
			),
		}, ret)
	})
	t.Run("using pit on effective date", func(t *testing.T) {
		t.Parallel()
		ret, err := store.GetAggregatedBalances(ctx, ledgercontroller.NewGetAggregatedBalancesQuery(ledgercontroller.PITFilter{
			PIT: pointer.For(now.Add(-time.Second)),
		}, query.Match("address", "users:"), false))
		require.NoError(t, err)
		require.Equal(t, ledger.BalancesByAssets{
			"USD": big.NewInt(0).Add(
				bigInt,
				smallInt,
			),
		}, ret)
	})
	t.Run("using pit on insertion date", func(t *testing.T) {
		t.Parallel()
		ret, err := store.GetAggregatedBalances(ctx, ledgercontroller.NewGetAggregatedBalancesQuery(ledgercontroller.PITFilter{
			PIT: pointer.For(now),
		}, query.Match("address", "users:"), true))
		require.NoError(t, err)
		require.Equal(t, ledger.BalancesByAssets{
			"USD": big.NewInt(0).Add(
				bigInt,
				smallInt,
			),
		}, ret)
	})
	t.Run("using a metadata and pit", func(t *testing.T) {
		t.Parallel()
		ret, err := store.GetAggregatedBalances(ctx, ledgercontroller.NewGetAggregatedBalancesQuery(ledgercontroller.PITFilter{
			PIT: pointer.For(now.Add(time.Minute)),
		}, query.Match("metadata[category]", "premium"), false))
		require.NoError(t, err)
		require.Equal(t, ledger.BalancesByAssets{
			"USD": big.NewInt(0).Add(
				big.NewInt(0).Mul(bigInt, big.NewInt(2)),
				big.NewInt(0),
			),
		}, ret)
	})
	t.Run("using a metadata without pit", func(t *testing.T) {
		t.Parallel()
		ret, err := store.GetAggregatedBalances(ctx, ledgercontroller.NewGetAggregatedBalancesQuery(ledgercontroller.PITFilter{},
			query.Match("metadata[category]", "premium"), false))
		require.NoError(t, err)
		require.Equal(t, ledger.BalancesByAssets{
			"USD": big.NewInt(0).Mul(bigInt, big.NewInt(2)),
		}, ret)
	})
	t.Run("when no matching", func(t *testing.T) {
		t.Parallel()
		ret, err := store.GetAggregatedBalances(ctx, ledgercontroller.NewGetAggregatedBalancesQuery(ledgercontroller.PITFilter{},
			query.Match("metadata[category]", "guest"), false))
		require.NoError(t, err)
		require.Equal(t, ledger.BalancesByAssets{}, ret)
	})

	t.Run("using a filter exist on metadata", func(t *testing.T) {
		t.Parallel()
		ret, err := store.GetAggregatedBalances(ctx, ledgercontroller.NewGetAggregatedBalancesQuery(ledgercontroller.PITFilter{}, query.Exists("metadata", "category"), false))
		require.NoError(t, err)
		require.Equal(t, ledger.BalancesByAssets{
			"USD": big.NewInt(0).Add(
				big.NewInt(0).Mul(bigInt, big.NewInt(2)),
				big.NewInt(0).Mul(smallInt, big.NewInt(2)),
			),
		}, ret)
	})
}

func TestAddToBalance(t *testing.T) {
	t.Parallel()

	store := newLedgerStore(t)
	ctx := logging.TestingContext()

	balance, err := store.AddToBalance(ctx, "world", "USD/2", big.NewInt(-100))
	require.NoError(t, err)
	require.Equal(t, int64(-100), balance.Int64())

	balance, err = store.AddToBalance(ctx, "world", "USD/2", big.NewInt(50))
	require.NoError(t, err)
	require.Equal(t, int64(-50), balance.Int64())
}
