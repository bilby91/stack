package ledger_test

import (
	"context"
	"fmt"
	ledgercontroller "github.com/formancehq/ledger/internal/controller/ledger"
	. "github.com/formancehq/ledger/internal/storage/ledger"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"math/big"
	"testing"

	"github.com/formancehq/stack/libs/go-libs/time"

	"github.com/pkg/errors"

	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/pointer"

	ledger "github.com/formancehq/ledger/internal"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/formancehq/stack/libs/go-libs/query"
	"github.com/stretchr/testify/require"
)

func TestGetTransactionWithVolumes(t *testing.T) {
	t.Parallel()
	store := newLedgerStore(t)
	now := time.Now()
	ctx := logging.TestingContext()

	tx1Data := ledger.NewTransactionData().
		WithPostings(
			ledger.NewPosting("world", "central_bank", "USD", big.NewInt(100)),
		).
		WithReference("tx1").
		WithDate(now.Add(-3 * time.Hour))
	tx1, err := store.InsertTransaction(ctx, tx1Data)
	require.NoError(t, err)

	tx2Data := ledger.NewTransactionData().
		WithPostings(
			ledger.NewPosting("world", "central_bank", "USD", big.NewInt(100)),
		).
		WithReference("tx2").
		WithDate(now.Add(-2 * time.Hour))
	tx2, err := store.InsertTransaction(ctx, tx2Data)
	require.NoError(t, err)

	tx, err := store.GetTransactionWithVolumes(ctx, ledgercontroller.NewGetTransactionQuery(tx1.ID).
		WithExpandVolumes().
		WithExpandEffectiveVolumes())
	require.NoError(t, err)
	require.Equal(t, tx1Data.Postings, tx.Postings)
	require.Equal(t, tx1Data.Reference, tx.Reference)
	require.Equal(t, tx1Data.Timestamp, tx.Timestamp)
	RequireEqual(t, ledger.AccountsAssetsVolumes{
		"world": {
			"USD": {
				Input:  big.NewInt(0),
				Output: big.NewInt(100),
			},
		},
		"central_bank": {
			"USD": {
				Input:  big.NewInt(100),
				Output: big.NewInt(0),
			},
		},
	}, tx.PostCommitVolumes)
	RequireEqual(t, ledger.AccountsAssetsVolumes{
		"world": {
			"USD": {
				Input:  big.NewInt(0),
				Output: big.NewInt(0),
			},
		},
		"central_bank": {
			"USD": {
				Input:  big.NewInt(0),
				Output: big.NewInt(0),
			},
		},
	}, tx.PreCommitVolumes)

	tx, err = store.GetTransactionWithVolumes(ctx, ledgercontroller.NewGetTransactionQuery(tx2.ID).
		WithExpandVolumes().
		WithExpandEffectiveVolumes())
	require.Equal(t, tx2Data.Postings, tx.Postings)
	require.Equal(t, tx2Data.Reference, tx.Reference)
	require.Equal(t, tx2Data.Timestamp, tx.Timestamp)
	RequireEqual(t, ledger.AccountsAssetsVolumes{
		"world": {
			"USD": {
				Input:  big.NewInt(0),
				Output: big.NewInt(200),
			},
		},
		"central_bank": {
			"USD": {
				Input:  big.NewInt(200),
				Output: big.NewInt(0),
			},
		},
	}, tx.PostCommitVolumes)
	RequireEqual(t, ledger.AccountsAssetsVolumes{
		"world": {
			"USD": {
				Input:  big.NewInt(0),
				Output: big.NewInt(100),
			},
		},
		"central_bank": {
			"USD": {
				Input:  big.NewInt(100),
				Output: big.NewInt(0),
			},
		},
	}, tx.PreCommitVolumes)
}

func TestGetTransaction(t *testing.T) {
	t.Parallel()
	store := newLedgerStore(t)
	now := time.Now()
	ctx := logging.TestingContext()

	tx1Data := ledger.NewTransactionData().
		WithPostings(
			ledger.NewPosting("world", "central_bank", "USD", big.NewInt(100)),
		).
		WithReference("tx1").
		WithDate(now.Add(-3 * time.Hour))
	tx1, err := store.InsertTransaction(ctx, tx1Data)
	require.NoError(t, err)

	tx2Data := ledger.NewTransactionData().
		WithPostings(
			ledger.NewPosting("world", "central_bank", "USD", big.NewInt(100)),
		).
		WithReference("tx2").
		WithDate(now.Add(-2 * time.Hour))
	_, err = store.InsertTransaction(ctx, tx2Data)
	require.NoError(t, err)

	tx, err := store.GetTransaction(context.Background(), tx1.ID)
	require.NoError(t, err)
	require.Equal(t, tx1.Postings, tx.Postings)
	require.Equal(t, tx1.Reference, tx.Reference)
	require.Equal(t, tx1.Timestamp, tx.Timestamp)
}

func TestGetTransactionByReference(t *testing.T) {
	t.Parallel()
	store := newLedgerStore(t)
	now := time.Now()
	ctx := logging.TestingContext()

	tx1Data := ledger.NewTransactionData().
		WithPostings(
			ledger.NewPosting("world", "central_bank", "USD", big.NewInt(100)),
		).
		WithReference("tx1").
		WithDate(now.Add(-3 * time.Hour))
	tx1, err := store.InsertTransaction(ctx, tx1Data)
	require.NoError(t, err)

	tx2Data := ledger.NewTransactionData().
		WithPostings(
			ledger.NewPosting("world", "central_bank", "USD", big.NewInt(100)),
		).
		WithReference("tx2").
		WithDate(now.Add(-2 * time.Hour))
	_, err = store.InsertTransaction(ctx, tx2Data)
	require.NoError(t, err)

	tx, err := store.GetTransactionByReference(context.Background(), "tx1")
	require.NoError(t, err)
	require.Equal(t, tx1.Postings, tx.Postings)
	require.Equal(t, tx1.Reference, tx.Reference)
	require.Equal(t, tx1.Timestamp, tx.Timestamp)
}

func TestCountTransactions(t *testing.T) {
	t.Parallel()
	store := newLedgerStore(t)

	for i := 0; i < 3; i++ {
		data := ledger.TransactionData{
			Postings: ledger.Postings{
				ledger.NewPosting("world", fmt.Sprintf("account%d", i), "USD", big.NewInt(100)),
			},
			Metadata: metadata.Metadata{},
		}
		_, err := store.InsertTransaction(logging.TestingContext(), data)
		require.NoError(t, err)
	}

	count, err := store.CountTransactions(context.Background(), ledgercontroller.NewGetTransactionsQuery(ledgercontroller.NewPaginatedQueryOptions(ledgercontroller.PITFilterWithVolumes{})))
	require.NoError(t, err, "counting transactions should not fail")
	require.Equal(t, 3, count, "count should be equal")
}

func TestUpdateTransactionsMetadata(t *testing.T) {
	t.Parallel()
	store := newLedgerStore(t)
	now := time.Now()
	ctx := logging.TestingContext()

	tx1Data := ledger.NewTransactionData().
		WithPostings(
			ledger.NewPosting("world", "alice", "USD", big.NewInt(100)),
		).
		WithDate(now.Add(-3 * time.Hour))
	tx1, err := store.InsertTransaction(ctx, tx1Data)
	require.NoError(t, err)

	tx2Data := ledger.NewTransactionData().
		WithPostings(
			ledger.NewPosting("world", "polo", "USD", big.NewInt(200)),
		).
		WithDate(now.Add(-2 * time.Hour))
	tx2, err := store.InsertTransaction(ctx, tx2Data)
	require.NoError(t, err)

	_, err = store.UpdateTransactionMetadata(ctx, tx1.ID, metadata.Metadata{"foo1": "bar2"})
	require.NoError(t, err)

	_, err = store.UpdateTransactionMetadata(ctx, tx2.ID, metadata.Metadata{"foo2": "bar2"})
	require.NoError(t, err)

	tx, err := store.GetTransactionWithVolumes(context.Background(), ledgercontroller.NewGetTransactionQuery(tx1.ID).WithExpandVolumes().WithExpandEffectiveVolumes())
	require.NoError(t, err, "getting transaction should not fail")
	require.Equal(t, tx.Metadata, metadata.Metadata{"foo1": "bar2"}, "metadata should be equal")

	tx, err = store.GetTransactionWithVolumes(context.Background(), ledgercontroller.NewGetTransactionQuery(tx2.ID).WithExpandVolumes().WithExpandEffectiveVolumes())
	require.NoError(t, err, "getting transaction should not fail")
	require.Equal(t, tx.Metadata, metadata.Metadata{"foo2": "bar2"}, "metadata should be equal")
}

func TestDeleteTransactionsMetadata(t *testing.T) {
	t.Parallel()
	store := newLedgerStore(t)
	now := time.Now()
	ctx := logging.TestingContext()

	tx1Data := ledger.NewTransactionData().
		WithPostings(
			ledger.NewPosting("world", "alice", "USD", big.NewInt(100)),
		).
		WithDate(now.Add(-3 * time.Hour))
	tx1, err := store.InsertTransaction(ctx, tx1Data)
	require.NoError(t, err)

	tx1, err = store.UpdateTransactionMetadata(ctx, tx1.ID, metadata.Metadata{"foo1": "bar1", "foo2": "bar2"})
	require.NoError(t, err)

	tx, err := store.GetTransaction(context.Background(), tx1.ID)
	require.NoError(t, err)
	require.Equal(t, tx.Metadata, metadata.Metadata{"foo1": "bar1", "foo2": "bar2"})

	tx1, err = store.DeleteTransactionMetadata(ctx, tx1.ID, "foo1")
	require.NoError(t, err)

	tx, err = store.GetTransaction(context.Background(), tx1.ID)
	require.NoError(t, err)
	require.Equal(t, metadata.Metadata{"foo2": "bar2"}, tx.Metadata)
}

func TestInsertTransactionInPast(t *testing.T) {
	t.Parallel()

	store := newLedgerStore(t)
	now := time.Now()
	ctx := logging.TestingContext()

	tx1Data := ledger.NewTransactionData().WithPostings(
		ledger.NewPosting("world", "bank", "USD/2", big.NewInt(100)),
	).WithDate(now)
	_, err := store.InsertTransaction(ctx, tx1Data)
	require.NoError(t, err)

	tx2Data := ledger.NewTransactionData().WithPostings(
		ledger.NewPosting("bank", "user1", "USD/2", big.NewInt(50)),
	).WithDate(now.Add(time.Hour))
	tx2, err := store.InsertTransaction(ctx, tx2Data)
	require.NoError(t, err)

	// Insert in past must modify pre/post commit volumes of tx2
	tx3Data := ledger.NewTransactionData().WithPostings(
		ledger.NewPosting("bank", "user2", "USD/2", big.NewInt(50)),
	).WithDate(now.Add(30 * time.Minute))
	_, err = store.InsertTransaction(ctx, tx3Data)
	require.NoError(t, err)

	// Insert before the oldest tx must update first_usage of involved accounts
	tx4Data := ledger.NewTransactionData().WithPostings(
		ledger.NewPosting("world", "bank", "USD/2", big.NewInt(100)),
	).WithDate(now.Add(-time.Minute))
	tx4, err := store.InsertTransaction(ctx, tx4Data)
	require.NoError(t, err)

	tx2FromDatabase, err := store.GetTransactionWithVolumes(ctx, ledgercontroller.NewGetTransactionQuery(tx2.ID).WithExpandVolumes().WithExpandEffectiveVolumes())
	require.NoError(t, err)

	RequireEqual(t, ledger.AccountsAssetsVolumes{
		"bank": {
			"USD/2": ledger.NewVolumesInt64(200, 50),
		},
		"user1": {
			"USD/2": ledger.NewVolumesInt64(0, 0),
		},
	}, tx2FromDatabase.PreCommitEffectiveVolumes)
	RequireEqual(t, ledger.AccountsAssetsVolumes{
		"bank": {
			"USD/2": ledger.NewVolumesInt64(200, 100),
		},
		"user1": {
			"USD/2": ledger.NewVolumesInt64(50, 0),
		},
	}, tx2FromDatabase.PostCommitEffectiveVolumes)

	account, err := store.GetAccount(ctx, "bank")
	require.NoError(t, err)
	require.Equal(t, tx4.Timestamp, account.FirstUsage)
}

func TestInsertTransactionInPastInOneBatch(t *testing.T) {
	t.Parallel()
	store := newLedgerStore(t)
	now := time.Now()
	ctx := logging.TestingContext()

	tx1Data := ledger.NewTransactionData().WithPostings(
		ledger.NewPosting("world", "bank", "USD/2", big.NewInt(100)),
	).WithDate(now)
	_, err := store.InsertTransaction(ctx, tx1Data)
	require.NoError(t, err)

	tx2Data := ledger.NewTransactionData().WithPostings(
		ledger.NewPosting("bank", "user1", "USD/2", big.NewInt(50)),
	).WithDate(now.Add(time.Hour))
	tx2, err := store.InsertTransaction(ctx, tx2Data)
	require.NoError(t, err)

	// Insert in past must modify pre/post commit volumes of tx2
	tx3Data := ledger.NewTransactionData().WithPostings(
		ledger.NewPosting("bank", "user2", "USD/2", big.NewInt(50)),
	).WithDate(now.Add(30 * time.Minute))
	_, err = store.InsertTransaction(ctx, tx3Data)
	require.NoError(t, err)

	tx2FromDatabase, err := store.GetTransactionWithVolumes(context.Background(), ledgercontroller.NewGetTransactionQuery(tx2.ID).WithExpandVolumes().WithExpandEffectiveVolumes())
	require.NoError(t, err)

	RequireEqual(t, ledger.AccountsAssetsVolumes{
		"bank": {
			"USD/2": ledger.NewVolumesInt64(100, 50),
		},
		"user1": {
			"USD/2": ledger.NewVolumesInt64(0, 0),
		},
	}, tx2FromDatabase.PreCommitEffectiveVolumes)
	RequireEqual(t, ledger.AccountsAssetsVolumes{
		"bank": {
			"USD/2": ledger.NewVolumesInt64(100, 100),
		},
		"user1": {
			"USD/2": ledger.NewVolumesInt64(50, 0),
		},
	}, tx2FromDatabase.PostCommitEffectiveVolumes)
}

func TestInsertTwoTransactionAtSameDateInSameBatch(t *testing.T) {
	t.Parallel()
	store := newLedgerStore(t)
	now := time.Now()
	ctx := logging.TestingContext()

	tx1Data := ledger.NewTransactionData().WithPostings(
		ledger.NewPosting("world", "bank", "USD/2", big.NewInt(100)),
	).WithDate(now.Add(-time.Hour))
	_, err := store.InsertTransaction(ctx, tx1Data)
	require.NoError(t, err)

	tx2Data := ledger.NewTransactionData().WithPostings(
		ledger.NewPosting("bank", "user1", "USD/2", big.NewInt(10)),
	).WithDate(now)
	tx2, err := store.InsertTransaction(ctx, tx2Data)
	require.NoError(t, err)

	tx3Data := ledger.NewTransactionData().WithPostings(
		ledger.NewPosting("bank", "user2", "USD/2", big.NewInt(10)),
	).WithDate(now)
	tx3, err := store.InsertTransaction(ctx, tx3Data)
	require.NoError(t, err)

	tx2FromDatabase, err := store.GetTransactionWithVolumes(context.Background(), ledgercontroller.NewGetTransactionQuery(tx2.ID).WithExpandVolumes().WithExpandEffectiveVolumes())
	require.NoError(t, err)

	RequireEqual(t, ledger.AccountsAssetsVolumes{
		"bank": {
			"USD/2": ledger.NewVolumesInt64(100, 10),
		},
		"user1": {
			"USD/2": ledger.NewVolumesInt64(10, 0),
		},
	}, tx2FromDatabase.PostCommitVolumes)
	RequireEqual(t, ledger.AccountsAssetsVolumes{
		"bank": {
			"USD/2": ledger.NewVolumesInt64(100, 0),
		},
		"user1": {
			"USD/2": ledger.NewVolumesInt64(0, 0),
		},
	}, tx2FromDatabase.PreCommitVolumes)

	tx3FromDatabase, err := store.GetTransactionWithVolumes(context.Background(), ledgercontroller.NewGetTransactionQuery(tx3.ID).WithExpandVolumes().WithExpandEffectiveVolumes())
	require.NoError(t, err)

	RequireEqual(t, ledger.AccountsAssetsVolumes{
		"bank": {
			"USD/2": ledger.NewVolumesInt64(100, 10),
		},
		"user2": {
			"USD/2": ledger.NewVolumesInt64(0, 0),
		},
	}, tx3FromDatabase.PreCommitVolumes)
	RequireEqual(t, ledger.AccountsAssetsVolumes{
		"bank": {
			"USD/2": ledger.NewVolumesInt64(100, 20),
		},
		"user2": {
			"USD/2": ledger.NewVolumesInt64(10, 0),
		},
	}, tx3FromDatabase.PostCommitVolumes)
}

func TestGetTransactions(t *testing.T) {
	t.Parallel()

	store := newLedgerStore(t)
	now := time.Now()
	ctx := logging.TestingContext()

	tx1Data := ledger.NewTransactionData().
		WithPostings(
			ledger.NewPosting("world", "alice", "USD", big.NewInt(100)),
		).
		WithMetadata(metadata.Metadata{"category": "1"}).
		WithDate(now.Add(-3 * time.Hour))
	tx1, err := store.InsertTransaction(ctx, tx1Data)
	require.NoError(t, err)

	tx2Data := ledger.NewTransactionData().
		WithPostings(
			ledger.NewPosting("world", "bob", "USD", big.NewInt(100)),
		).
		WithMetadata(metadata.Metadata{"category": "2"}).
		WithDate(now.Add(-2 * time.Hour))
	tx2, err := store.InsertTransaction(ctx, tx2Data)
	require.NoError(t, err)

	tx3Data := ledger.NewTransactionData().
		WithPostings(
			ledger.NewPosting("world", "users:marley", "USD", big.NewInt(100)),
		).
		WithMetadata(metadata.Metadata{"category": "3"}).
		WithDate(now.Add(-time.Hour))
	tx3, err := store.InsertTransaction(ctx, tx3Data)
	require.NoError(t, err)

	tx3AfterRevert, hasBeenReverted, err := store.RevertTransaction(ctx, tx3.ID)
	require.NoError(t, err)
	require.True(t, hasBeenReverted)

	tx4, err := store.InsertTransaction(ctx, tx3Data.Reverse(false).WithDate(now))
	require.NoError(t, err)

	tx3AfterRevert, err = store.UpdateTransactionMetadata(ctx, tx3AfterRevert.ID, metadata.Metadata{
		"additional_metadata": "true",
	})

	tx5Data := ledger.NewTransactionData().
		WithPostings(
			ledger.NewPosting("users:marley", "sellers:amazon", "USD", big.NewInt(100)),
		).
		WithDate(now)
	tx5, err := store.InsertTransaction(ctx, tx5Data)
	require.NoError(t, err)

	type testCase struct {
		name        string
		query       ledgercontroller.PaginatedQueryOptions[ledgercontroller.PITFilterWithVolumes]
		expected    []ledger.Transaction
		expectError error
	}
	testCases := []testCase{
		{
			name:     "nominal",
			query:    ledgercontroller.NewPaginatedQueryOptions(ledgercontroller.PITFilterWithVolumes{}),
			expected: []ledger.Transaction{*tx5, *tx4, *tx3AfterRevert, *tx2, *tx1},
		},
		{
			name: "address filter",
			query: ledgercontroller.NewPaginatedQueryOptions(ledgercontroller.PITFilterWithVolumes{}).
				WithQueryBuilder(query.Match("account", "bob")),
			expected: []ledger.Transaction{*tx2},
		},
		{
			name: "address filter using segments matching two addresses by individual segments",
			query: ledgercontroller.NewPaginatedQueryOptions(ledgercontroller.PITFilterWithVolumes{}).
				WithQueryBuilder(query.Match("account", "users:amazon")),
			expected: []ledger.Transaction{},
		},
		{
			name: "address filter using segment",
			query: ledgercontroller.NewPaginatedQueryOptions(ledgercontroller.PITFilterWithVolumes{}).
				WithQueryBuilder(query.Match("account", "users:")),
			expected: []ledger.Transaction{*tx5, *tx4, *tx3AfterRevert},
		},
		{
			name: "filter using metadata",
			query: ledgercontroller.NewPaginatedQueryOptions(ledgercontroller.PITFilterWithVolumes{}).
				WithQueryBuilder(query.Match("metadata[category]", "2")),
			expected: []ledger.Transaction{*tx2},
		},
		{
			name: "using point in time",
			query: ledgercontroller.NewPaginatedQueryOptions(ledgercontroller.PITFilterWithVolumes{
				PITFilter: ledgercontroller.PITFilter{
					PIT: pointer.For(now.Add(-time.Hour)),
				},
			}),
			expected: []ledger.Transaction{*tx3, *tx2, *tx1},
		},
		{
			name: "filter using invalid key",
			query: ledgercontroller.NewPaginatedQueryOptions(ledgercontroller.PITFilterWithVolumes{}).
				WithQueryBuilder(query.Match("invalid", "2")),
			expectError: &ErrInvalidQuery{},
		},
		{
			name: "reverted transactions",
			query: ledgercontroller.NewPaginatedQueryOptions(ledgercontroller.PITFilterWithVolumes{}).
				WithQueryBuilder(query.Match("reverted", true)),
			expected: []ledger.Transaction{*tx3AfterRevert},
		},
		{
			name: "filter using exists metadata",
			query: ledgercontroller.NewPaginatedQueryOptions(ledgercontroller.PITFilterWithVolumes{}).
				WithQueryBuilder(query.Exists("metadata", "category")),
			expected: []ledger.Transaction{*tx3AfterRevert, *tx2, *tx1},
		},
		{
			name: "filter using exists metadata2",
			query: ledgercontroller.NewPaginatedQueryOptions(ledgercontroller.PITFilterWithVolumes{}).
				WithQueryBuilder(query.Not(query.Exists("metadata", "category"))),
			expected: []ledger.Transaction{*tx5, *tx4},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tc.query.Options.ExpandVolumes = true
			tc.query.Options.ExpandEffectiveVolumes = false
			cursor, err := store.GetTransactions(ctx, ledgercontroller.NewGetTransactionsQuery(tc.query))
			if tc.expectError != nil {
				require.True(t, errors.Is(err, tc.expectError))
			} else {
				require.NoError(t, err)
				require.Len(t, cursor.Data, len(tc.expected))
				RequireEqual(t, tc.expected, collectionutils.Map(cursor.Data, ledger.ExpandedTransaction.Base))

				count, err := store.CountTransactions(ctx, ledgercontroller.NewGetTransactionsQuery(tc.query))
				require.NoError(t, err)

				require.EqualValues(t, len(tc.expected), count)
			}
		})
	}
}

func TestGetLastTransaction(t *testing.T) {
	t.Parallel()
	store := newLedgerStore(t)
	ctx := logging.TestingContext()

	tx1Data := ledger.NewTransactionData().
		WithPostings(
			ledger.NewPosting("world", "alice", "USD", big.NewInt(100)),
		)
	_, err := store.InsertTransaction(ctx, tx1Data)
	require.NoError(t, err)

	tx2Data := ledger.NewTransactionData().
		WithPostings(
			ledger.NewPosting("world", "bob", "USD", big.NewInt(100)),
		)
	_, err = store.InsertTransaction(ctx, tx2Data)
	require.NoError(t, err)

	tx3Data := ledger.NewTransactionData().
		WithPostings(
			ledger.NewPosting("world", "users:marley", "USD", big.NewInt(100)),
		)
	tx3, err := store.InsertTransaction(ctx, tx3Data)
	require.NoError(t, err)

	tx, err := store.GetLastTransaction(ctx)
	require.NoError(t, err)
	require.Equal(t, *tx3, tx.Transaction)
}
