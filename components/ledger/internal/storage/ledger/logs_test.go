//go:build it

package ledger_test

import (
	"context"
	"fmt"
	ledgercontroller "github.com/formancehq/ledger/internal/controller/ledger"
	"github.com/formancehq/stack/libs/go-libs/platform/postgres"
	"math/big"
	"testing"

	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/formancehq/stack/libs/go-libs/time"

	"github.com/formancehq/stack/libs/go-libs/logging"

	ledger "github.com/formancehq/ledger/internal"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/formancehq/stack/libs/go-libs/query"
	"github.com/stretchr/testify/require"
)

func TestGetLastLog(t *testing.T) {
	t.Parallel()
	store := newLedgerStore(t)
	now := time.Now()
	ctx := logging.TestingContext()

	lastLog, err := store.GetLastLog(context.Background())
	require.Error(t, err)
	require.True(t, postgres.IsNotFoundError(err))
	require.Nil(t, lastLog)
	tx1 := ledger.ExpandedTransaction{
		Transaction: ledger.Transaction{
			TransactionData: ledger.TransactionData{
				Postings: []ledger.Posting{
					{
						Source:      "world",
						Destination: "central_bank",
						Amount:      big.NewInt(100),
						Asset:       "USD",
					},
				},
				Reference: "tx1",
				Timestamp: now.Add(-3 * time.Hour),
			},
		},
		PostCommitVolumes: ledger.AccountsAssetsVolumes{
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
		},
		PreCommitVolumes: ledger.AccountsAssetsVolumes{
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
		},
	}

	logTx := ledger.NewTransactionLog(tx1.Transaction, map[string]metadata.Metadata{})
	_, err = store.InsertLog(ctx, logTx)
	require.NoError(t, err)

	lastLog, err = store.GetLastLog(ctx)
	require.NoError(t, err)
	require.NotNil(t, lastLog)

	require.Equal(t, tx1.Postings, lastLog.Data.(ledger.NewTransactionLogPayload).Transaction.Postings)
	require.Equal(t, tx1.Reference, lastLog.Data.(ledger.NewTransactionLogPayload).Transaction.Reference)
	require.Equal(t, tx1.Timestamp, lastLog.Data.(ledger.NewTransactionLogPayload).Transaction.Timestamp)
}

func TestReadLogWithIdempotencyKey(t *testing.T) {
	t.Parallel()

	store := newLedgerStore(t)
	ctx := logging.TestingContext()

	logTx := ledger.NewTransactionLog(
		ledger.NewTransaction().
			WithPostings(
				ledger.NewPosting("world", "bank", "USD", big.NewInt(100)),
			),
		map[string]metadata.Metadata{},
	)
	log := logTx.WithIdempotencyKey("test")
	chainedLog, err := store.InsertLog(ctx, log)
	require.NoError(t, err)

	lastLog, err := store.ReadLogWithIdempotencyKey(context.Background(), "test")
	require.NoError(t, err)
	require.NotNil(t, lastLog)
	require.Equal(t, *chainedLog, *lastLog)
}

func TestGetLogs(t *testing.T) {
	t.Parallel()
	store := newLedgerStore(t)
	now := time.Now()
	ctx := logging.TestingContext()

	tx1 := ledger.ExpandedTransaction{
		Transaction: ledger.Transaction{
			TransactionData: ledger.TransactionData{
				Postings: []ledger.Posting{
					{
						Source:      "world",
						Destination: "central_bank",
						Amount:      big.NewInt(100),
						Asset:       "USD",
					},
				},
				Reference: "tx1",
				Timestamp: now.Add(-3 * time.Hour),
			},
		},
		PostCommitVolumes: ledger.AccountsAssetsVolumes{
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
		},
		PreCommitVolumes: ledger.AccountsAssetsVolumes{
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
		},
	}
	tx2 := ledger.ExpandedTransaction{
		Transaction: ledger.Transaction{
			ID: 1,
			TransactionData: ledger.TransactionData{
				Postings: []ledger.Posting{
					{
						Source:      "world",
						Destination: "central_bank",
						Amount:      big.NewInt(100),
						Asset:       "USD",
					},
				},
				Reference: "tx2",
				Timestamp: now.Add(-2 * time.Hour),
			},
		},
		PostCommitVolumes: ledger.AccountsAssetsVolumes{
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
		},
		PreCommitVolumes: ledger.AccountsAssetsVolumes{
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
		},
	}
	tx3 := ledger.ExpandedTransaction{
		Transaction: ledger.Transaction{
			ID: 2,
			TransactionData: ledger.TransactionData{
				Postings: []ledger.Posting{
					{
						Source:      "central_bank",
						Destination: "users:1",
						Amount:      big.NewInt(1),
						Asset:       "USD",
					},
				},
				Reference: "tx3",
				Metadata: metadata.Metadata{
					"priority": "high",
				},
				Timestamp: now.Add(-1 * time.Hour),
			},
		},
		PreCommitVolumes: ledger.AccountsAssetsVolumes{
			"central_bank": {
				"USD": {
					Input:  big.NewInt(200),
					Output: big.NewInt(0),
				},
			},
			"users:1": {
				"USD": {
					Input:  big.NewInt(0),
					Output: big.NewInt(0),
				},
			},
		},
		PostCommitVolumes: ledger.AccountsAssetsVolumes{
			"central_bank": {
				"USD": {
					Input:  big.NewInt(200),
					Output: big.NewInt(1),
				},
			},
			"users:1": {
				"USD": {
					Input:  big.NewInt(1),
					Output: big.NewInt(0),
				},
			},
		},
	}

	for _, tx := range []ledger.ExpandedTransaction{tx1, tx2, tx3} {
		newLog := ledger.NewTransactionLog(tx.Transaction, map[string]metadata.Metadata{}).
			WithDate(tx.Timestamp)

		_, err := store.InsertLog(ctx, newLog)
		require.NoError(t, err)
	}

	cursor, err := store.GetLogs(context.Background(), ledgercontroller.NewGetLogsQuery(ledgercontroller.NewPaginatedQueryOptions[any](nil)))
	require.NoError(t, err)
	require.Equal(t, bunpaginate.QueryDefaultPageSize, cursor.PageSize)

	require.Equal(t, 3, len(cursor.Data))
	require.EqualValues(t, 2, cursor.Data[0].ID)
	require.Equal(t, tx3.Postings, cursor.Data[0].Data.(ledger.NewTransactionLogPayload).Transaction.Postings)
	require.Equal(t, tx3.Reference, cursor.Data[0].Data.(ledger.NewTransactionLogPayload).Transaction.Reference)
	require.Equal(t, tx3.Timestamp, cursor.Data[0].Data.(ledger.NewTransactionLogPayload).Transaction.Timestamp)

	cursor, err = store.GetLogs(context.Background(), ledgercontroller.NewGetLogsQuery(ledgercontroller.NewPaginatedQueryOptions[any](nil).WithPageSize(1)))
	require.NoError(t, err)
	// Should get only the first log.
	require.Equal(t, 1, cursor.PageSize)
	require.EqualValues(t, 2, cursor.Data[0].ID)

	cursor, err = store.GetLogs(context.Background(), ledgercontroller.NewGetLogsQuery(ledgercontroller.NewPaginatedQueryOptions[any](nil).
		WithQueryBuilder(query.And(
			query.Gte("date", now.Add(-2*time.Hour)),
			query.Lt("date", now.Add(-time.Hour)),
		)).
		WithPageSize(10),
	))
	require.NoError(t, err)
	require.Equal(t, 10, cursor.PageSize)
	// Should get only the second log, as StartTime is inclusive and EndTime exclusive.
	require.Len(t, cursor.Data, 1)
	require.EqualValues(t, 1, cursor.Data[0].ID)
}

func TestGetBalance(t *testing.T) {
	t.Parallel()
	store := newLedgerStore(t)
	ctx := logging.TestingContext()

	const (
		batchNumber = 100
		batchSize   = 10
		input       = 100
		output      = 10
	)

	for i := 0; i < batchNumber; i++ {
		for j := 0; j < batchSize; j++ {
			_, err := store.InsertTransaction(ctx, ledger.NewTransactionData().WithPostings(
				ledger.NewPosting("world", fmt.Sprintf("account:%d", j), "EUR/2", big.NewInt(input)),
				ledger.NewPosting(fmt.Sprintf("account:%d", j), "starbucks", "EUR/2", big.NewInt(output)),
			))
			require.NoError(t, err)
		}
	}

	balance, err := store.GetBalance(context.Background(), "account:1", "EUR/2")
	require.NoError(t, err)
	require.Equal(t, big.NewInt((input-output)*batchNumber), balance)
}
