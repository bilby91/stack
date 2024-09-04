//go:build it

package ledger_test

import (
	"context"
	"database/sql"
	ledgercontroller "github.com/formancehq/ledger/internal/controller/ledger"
	. "github.com/formancehq/ledger/internal/storage/ledger"
	"github.com/formancehq/stack/libs/go-libs/testing/utils"
	"math/big"
	"testing"

	"github.com/formancehq/stack/libs/go-libs/time"

	"github.com/formancehq/stack/libs/go-libs/logging"

	ledger "github.com/formancehq/ledger/internal"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/formancehq/stack/libs/go-libs/query"
	"github.com/stretchr/testify/require"
)

func TestGetAccounts(t *testing.T) {
	t.Parallel()
	store := newLedgerStore(t)
	now := time.Now()
	ctx := logging.TestingContext()

	_, err := store.InsertTransaction(ctx, ledger.NewTransactionData().
		WithPostings(ledger.NewPosting("world", "account:1", "USD", big.NewInt(100))).
		WithDate(now).
		WithInsertedAt(now))
	require.NoError(t, err)

	require.NoError(t, store.UpdateAccountMetadata(ctx, "account:1", map[string]string{
		"category": "4",
	}))

	require.NoError(t, store.UpdateAccountMetadata(ctx, "account:1", map[string]string{
		"category": "1",
	}))
	require.NoError(t, store.UpdateAccountMetadata(ctx, "account:2", map[string]string{
		"category": "2",
	}))
	require.NoError(t, store.UpdateAccountMetadata(ctx, "account:3", map[string]string{
		"category": "3",
	}))
	require.NoError(t, store.UpdateAccountMetadata(ctx, "orders:1", map[string]string{
		"foo": "bar",
	}))
	require.NoError(t, store.UpdateAccountMetadata(ctx, "orders:2", map[string]string{
		"foo": "bar",
	}))

	_, err = store.InsertTransaction(ctx, ledger.NewTransactionData().
		WithPostings(ledger.NewPosting("world", "account:1", "USD", big.NewInt(100))).
		WithDate(now.Add(4*time.Minute)).
		WithInsertedAt(now.Add(100*time.Millisecond)))
	require.NoError(t, err)

	_, err = store.InsertTransaction(ctx, ledger.NewTransactionData().
		WithPostings(ledger.NewPosting("account:1", "bank", "USD", big.NewInt(50))).
		WithDate(now.Add(3*time.Minute)).
		WithInsertedAt(now.Add(200*time.Millisecond)))
	require.NoError(t, err)

	_, err = store.InsertTransaction(ctx, ledger.NewTransactionData().
		WithPostings(ledger.NewPosting("world", "account:1", "USD", big.NewInt(0))).
		WithDate(now.Add(-time.Minute)).
		WithInsertedAt(now.Add(200*time.Millisecond)))
	require.NoError(t, err)

	t.Run("list all", func(t *testing.T) {
		t.Parallel()
		accounts, err := store.GetAccountsWithVolumes(ctx, ledgercontroller.NewGetAccountsQuery(ledgercontroller.NewPaginatedQueryOptions(ledgercontroller.PITFilterWithVolumes{})))
		require.NoError(t, err)
		require.Len(t, accounts.Data, 7)
	})

	t.Run("list using metadata", func(t *testing.T) {
		t.Parallel()
		accounts, err := store.GetAccountsWithVolumes(ctx, ledgercontroller.NewGetAccountsQuery(ledgercontroller.NewPaginatedQueryOptions(ledgercontroller.PITFilterWithVolumes{}).
			WithQueryBuilder(query.Match("metadata[category]", "1")),
		))
		require.NoError(t, err)
		require.Len(t, accounts.Data, 1)
	})

	t.Run("list before date", func(t *testing.T) {
		t.Parallel()
		accounts, err := store.GetAccountsWithVolumes(ctx, ledgercontroller.NewGetAccountsQuery(ledgercontroller.NewPaginatedQueryOptions(ledgercontroller.PITFilterWithVolumes{
			PITFilter: ledgercontroller.PITFilter{
				PIT: &now,
			},
		})))
		require.NoError(t, err)
		require.Len(t, accounts.Data, 2)
	})

	t.Run("list with volumes", func(t *testing.T) {
		t.Parallel()

		accounts, err := store.GetAccountsWithVolumes(ctx, ledgercontroller.NewGetAccountsQuery(ledgercontroller.NewPaginatedQueryOptions(ledgercontroller.PITFilterWithVolumes{
			ExpandVolumes: true,
		}).WithQueryBuilder(query.Match("address", "account:1"))))
		require.NoError(t, err)
		require.Len(t, accounts.Data, 1)
		require.Equal(t, ledger.VolumesByAssets{
			"USD": ledger.NewVolumesInt64(200, 50),
		}, accounts.Data[0].Volumes)
	})

	t.Run("list with volumes using PIT", func(t *testing.T) {
		t.Parallel()

		accounts, err := store.GetAccountsWithVolumes(ctx, ledgercontroller.NewGetAccountsQuery(ledgercontroller.NewPaginatedQueryOptions(ledgercontroller.PITFilterWithVolumes{
			PITFilter: ledgercontroller.PITFilter{
				PIT: &now,
			},
			ExpandVolumes: true,
		}).WithQueryBuilder(query.Match("address", "account:1"))))
		require.NoError(t, err)
		require.Len(t, accounts.Data, 1)
		require.Equal(t, ledger.VolumesByAssets{
			"USD": ledger.NewVolumesInt64(100, 0),
		}, accounts.Data[0].Volumes)
	})

	t.Run("list with effective volumes", func(t *testing.T) {
		t.Parallel()

		accounts, err := store.GetAccountsWithVolumes(ctx, ledgercontroller.NewGetAccountsQuery(ledgercontroller.NewPaginatedQueryOptions(ledgercontroller.PITFilterWithVolumes{
			ExpandEffectiveVolumes: true,
		}).WithQueryBuilder(query.Match("address", "account:1"))))
		require.NoError(t, err)
		require.Len(t, accounts.Data, 1)
		require.Equal(t, ledger.VolumesByAssets{
			"USD": ledger.NewVolumesInt64(200, 50),
		}, accounts.Data[0].EffectiveVolumes)
	})

	t.Run("list with effective volumes using PIT", func(t *testing.T) {
		t.Parallel()
		accounts, err := store.GetAccountsWithVolumes(ctx, ledgercontroller.NewGetAccountsQuery(ledgercontroller.NewPaginatedQueryOptions(ledgercontroller.PITFilterWithVolumes{
			PITFilter: ledgercontroller.PITFilter{
				PIT: &now,
			},
			ExpandEffectiveVolumes: true,
		}).WithQueryBuilder(query.Match("address", "account:1"))))
		require.NoError(t, err)
		require.Len(t, accounts.Data, 1)
		require.Equal(t, ledger.VolumesByAssets{
			"USD": ledger.NewVolumesInt64(100, 0),
		}, accounts.Data[0].EffectiveVolumes)
	})

	t.Run("list using filter on address", func(t *testing.T) {
		t.Parallel()
		accounts, err := store.GetAccountsWithVolumes(ctx, ledgercontroller.NewGetAccountsQuery(ledgercontroller.NewPaginatedQueryOptions(ledgercontroller.PITFilterWithVolumes{}).
			WithQueryBuilder(query.Match("address", "account:")),
		))
		require.NoError(t, err)
		require.Len(t, accounts.Data, 3)
	})
	t.Run("list using filter on multiple address", func(t *testing.T) {
		t.Parallel()
		accounts, err := store.GetAccountsWithVolumes(ctx, ledgercontroller.NewGetAccountsQuery(ledgercontroller.NewPaginatedQueryOptions(ledgercontroller.PITFilterWithVolumes{}).
			WithQueryBuilder(
				query.Or(
					query.Match("address", "account:1"),
					query.Match("address", "orders:"),
				),
			),
		))
		require.NoError(t, err)
		require.Len(t, accounts.Data, 3)
	})
	t.Run("list using filter on balances", func(t *testing.T) {
		t.Parallel()
		accounts, err := store.GetAccountsWithVolumes(ctx, ledgercontroller.NewGetAccountsQuery(ledgercontroller.NewPaginatedQueryOptions(ledgercontroller.PITFilterWithVolumes{}).
			WithQueryBuilder(query.Lt("balance[USD]", 0)),
		))
		require.NoError(t, err)
		require.Len(t, accounts.Data, 1) // world

		accounts, err = store.GetAccountsWithVolumes(ctx, ledgercontroller.NewGetAccountsQuery(ledgercontroller.NewPaginatedQueryOptions(ledgercontroller.PITFilterWithVolumes{}).
			WithQueryBuilder(query.Gt("balance[USD]", 0)),
		))
		require.NoError(t, err)
		require.Len(t, accounts.Data, 2)
		require.Equal(t, "account:1", accounts.Data[0].Account.Address)
		require.Equal(t, "bank", accounts.Data[1].Account.Address)
	})

	t.Run("list using filter on exists metadata", func(t *testing.T) {
		t.Parallel()
		accounts, err := store.GetAccountsWithVolumes(ctx, ledgercontroller.NewGetAccountsQuery(ledgercontroller.NewPaginatedQueryOptions(ledgercontroller.PITFilterWithVolumes{}).
			WithQueryBuilder(query.Exists("metadata", "foo")),
		))
		require.NoError(t, err)
		require.Len(t, accounts.Data, 2)

		accounts, err = store.GetAccountsWithVolumes(ctx, ledgercontroller.NewGetAccountsQuery(ledgercontroller.NewPaginatedQueryOptions(ledgercontroller.PITFilterWithVolumes{}).
			WithQueryBuilder(query.Exists("metadata", "category")),
		))
		require.NoError(t, err)
		require.Len(t, accounts.Data, 3)
	})

	t.Run("list using filter invalid field", func(t *testing.T) {
		t.Parallel()
		_, err := store.GetAccountsWithVolumes(ctx, ledgercontroller.NewGetAccountsQuery(ledgercontroller.NewPaginatedQueryOptions(ledgercontroller.PITFilterWithVolumes{}).
			WithQueryBuilder(query.Lt("invalid", 0)),
		))
		require.Error(t, err)
		require.True(t, IsErrInvalidQuery(err))
	})
}

func TestUpdateAccountsMetadata(t *testing.T) {
	t.Parallel()
	store := newLedgerStore(t)

	metadata := metadata.Metadata{
		"foo": "bar",
	}

	ctx := logging.TestingContext()

	require.NoError(t, store.UpdateAccountMetadata(ctx, "bank", metadata))

	account, err := store.GetAccountWithVolumes(context.Background(), ledgercontroller.NewGetAccountQuery("bank"))
	require.NoError(t, err, "account retrieval should not fail")

	require.Equal(t, "bank", account.Address, "account address should match")
	require.Equal(t, metadata, account.Metadata, "account metadata should match")
}

func TestGetAccount(t *testing.T) {
	t.Parallel()

	store := newLedgerStore(t)
	now := time.Now()
	ctx := logging.TestingContext()

	_, err := store.InsertTransaction(ctx, ledger.NewTransactionData().WithPostings(
		ledger.NewPosting("world", "multi", "USD/2", big.NewInt(100)),
	).WithDate(now))
	require.NoError(t, err)

	require.NoError(t, store.UpdateAccountMetadata(ctx, "multi", metadata.Metadata{
		"category": "gold",
	}))

	_, err = store.InsertTransaction(ctx, ledger.NewTransactionData().WithPostings(
		ledger.NewPosting("world", "multi", "USD/2", big.NewInt(0)),
	).WithDate(now.Add(-time.Minute)))
	require.NoError(t, err)

	t.Run("find account", func(t *testing.T) {
		t.Parallel()
		account, err := store.GetAccountWithVolumes(ctx, ledgercontroller.NewGetAccountQuery("multi"))
		require.NoError(t, err)
		require.Equal(t, ledger.ExpandedAccount{
			Account: ledger.Account{
				Address: "multi",
				Metadata: metadata.Metadata{
					"category": "gold",
				},
				FirstUsage: now.Add(-time.Minute),
			},
		}, *account)

		account, err = store.GetAccountWithVolumes(ctx, ledgercontroller.NewGetAccountQuery("world"))
		require.NoError(t, err)
		require.Equal(t, ledger.ExpandedAccount{
			Account: ledger.Account{
				Address:    "world",
				Metadata:   metadata.Metadata{},
				FirstUsage: now.Add(-time.Minute),
			},
		}, *account)
	})

	t.Run("find account in past", func(t *testing.T) {
		t.Parallel()
		account, err := store.GetAccountWithVolumes(ctx, ledgercontroller.NewGetAccountQuery("multi").WithPIT(now.Add(-30*time.Second)))
		require.NoError(t, err)
		require.Equal(t, ledger.ExpandedAccount{
			Account: ledger.Account{
				Address:    "multi",
				Metadata:   metadata.Metadata{},
				FirstUsage: now.Add(-time.Minute),
			},
		}, *account)
	})

	t.Run("find account with volumes", func(t *testing.T) {
		t.Parallel()
		account, err := store.GetAccountWithVolumes(ctx, ledgercontroller.NewGetAccountQuery("multi").
			WithExpandVolumes())
		require.NoError(t, err)
		require.Equal(t, ledger.ExpandedAccount{
			Account: ledger.Account{
				Address: "multi",
				Metadata: metadata.Metadata{
					"category": "gold",
				},
				FirstUsage: now.Add(-time.Minute),
			},
			Volumes: ledger.VolumesByAssets{
				"USD/2": ledger.NewVolumesInt64(100, 0),
			},
		}, *account)
	})

	t.Run("find account with effective volumes", func(t *testing.T) {
		t.Parallel()
		account, err := store.GetAccountWithVolumes(ctx, ledgercontroller.NewGetAccountQuery("multi").
			WithExpandEffectiveVolumes())
		require.NoError(t, err)
		require.Equal(t, ledger.ExpandedAccount{
			Account: ledger.Account{
				Address: "multi",
				Metadata: metadata.Metadata{
					"category": "gold",
				},
				FirstUsage: now.Add(-time.Minute),
			},
			EffectiveVolumes: ledger.VolumesByAssets{
				"USD/2": ledger.NewVolumesInt64(100, 0),
			},
		}, *account)
	})

	t.Run("find account using pit", func(t *testing.T) {
		t.Parallel()
		account, err := store.GetAccountWithVolumes(ctx, ledgercontroller.NewGetAccountQuery("multi").WithPIT(now))
		require.NoError(t, err)
		require.Equal(t, ledger.ExpandedAccount{
			Account: ledger.Account{
				Address:    "multi",
				Metadata:   metadata.Metadata{},
				FirstUsage: now.Add(-time.Minute),
			},
		}, *account)
	})

	t.Run("not existent account", func(t *testing.T) {
		t.Parallel()
		_, err := store.GetAccountWithVolumes(ctx, ledgercontroller.NewGetAccountQuery("account_not_existing"))
		require.Error(t, err)
	})

}

func TestGetAccountWithVolumes(t *testing.T) {
	t.Parallel()
	store := newLedgerStore(t)
	ctx := logging.TestingContext()
	now := time.Now()

	bigInt, _ := big.NewInt(0).SetString("999999999999999999999999999999999999999999999999999999999999999999999999999999999999999", 10)

	_, err := store.InsertTransaction(ctx, ledger.NewTransactionData().WithPostings(
		ledger.NewPosting("world", "multi", "USD/2", bigInt),
	).WithDate(now))
	require.NoError(t, err)

	accountWithVolumes, err := store.GetAccountWithVolumes(ctx,
		ledgercontroller.NewGetAccountQuery("multi").WithExpandVolumes())
	require.NoError(t, err)
	require.Equal(t, &ledger.ExpandedAccount{
		Account: ledger.Account{
			Address:    "multi",
			Metadata:   metadata.Metadata{},
			FirstUsage: now,
		},
		Volumes: map[string]*ledger.Volumes{
			"USD/2": ledger.NewEmptyVolumes().WithInput(bigInt),
		},
	}, accountWithVolumes)
}

func TestUpdateAccountMetadata(t *testing.T) {
	t.Parallel()
	store := newLedgerStore(t)
	ctx := logging.TestingContext()

	require.NoError(t, store.UpdateAccountMetadata(ctx, "central_bank", metadata.Metadata{
		"foo": "bar",
	}))

	account, err := store.GetAccountWithVolumes(ctx, ledgercontroller.NewGetAccountQuery("central_bank"))
	require.NoError(t, err)
	require.EqualValues(t, "bar", account.Metadata["foo"])
}

func TestCountAccounts(t *testing.T) {
	t.Parallel()

	store := newLedgerStore(t)
	ctx := logging.TestingContext()

	_, err := store.InsertTransaction(ctx, ledger.NewTransactionData().WithPostings(
		ledger.NewPosting("world", "central_bank", "USD/2", big.NewInt(100)),
	))
	require.NoError(t, err)

	countAccounts, err := store.CountAccounts(ctx, ledgercontroller.NewGetAccountsQuery(ledgercontroller.NewPaginatedQueryOptions(ledgercontroller.PITFilterWithVolumes{})))
	require.NoError(t, err)
	require.EqualValues(t, 2, countAccounts) // world + central_bank
}

func TestUpsertAccount(t *testing.T) {
	t.Parallel()

	store := newLedgerStore(t)
	ctx := logging.TestingContext()

	now := time.Now()

	account := ledger.Account{
		Address:       "foo",
		FirstUsage:    now,
		InsertionDate: now,
		UpdatedAt:     now,
	}

	// initial insert
	_, err := store.UpsertAccount(ctx, account)
	require.NoError(t, err)

	//accountFromDB, err := store.GetAccount(ctx, account.Address)
	//require.NoError(t, err)
	//require.Equal(t, account, *accountFromDB)

	//// update metadata and check database
	//account.Metadata = metadata.Metadata{
	//	"foo": "bar",
	//}
	//
	//_, err = store.UpsertAccount(ctx, account)
	//require.NoError(t, err)
	//
	//utils.DumpTables(t, ctx, store.GetDB(), "select * from "+store.PrefixWithBucket("accounts"))
	//
	//accountFromDB, err = store.GetAccount(ctx, account.Address)
	//require.NoError(t, err)
	//require.Equal(t, account, *accountFromDB)
	//
	//// update first_usage and check database
	//account.FirstUsage = now.Add(-time.Minute)
	//
	//_, err = store.UpsertAccount(ctx, account)
	//require.NoError(t, err)
	//
	//accountFromDB, err = store.GetAccount(ctx, account.Address)
	//require.NoError(t, err)
	//require.Equal(t, account, *accountFromDB)

	// upsert with no modification
	tx, err := store.GetDB().BeginTx(ctx, &sql.TxOptions{})
	require.NoError(t, err)
	defer func() {
		require.NoError(t, tx.Rollback())
	}()

	store = store.WithDB(tx)

	utils.DumpTables(t, ctx, tx,
		//"select * from "+store.PrefixWithBucket("accounts"),
		//`SELECT query FROM pg_locks l JOIN pg_stat_activity a ON l.pid = a.pid`,
		//`select * from pg_class`,
		`select pid, mode, relname from pg_locks join pg_class on pg_class.oid = pg_locks.relation`,
	)

	upserted, err := store.UpsertAccount(ctx, account)
	require.NoError(t, err)
	require.False(t, upserted)

	utils.DumpTables(t, ctx, tx,
		//"select * from "+store.PrefixWithBucket("accounts"),
		//`SELECT query FROM pg_locks l JOIN pg_stat_activity a ON l.pid = a.pid`,
		//`select * from pg_class`,
		`select pid, mode, relname, reltype from pg_locks join pg_class on pg_class.oid = pg_locks.relation`,
		`select * from pg_class where relname = 'accounts_seq_seq'`,
		`select * from pg_authid where oid = 10`,
		`select * from pg_indexes where schemaname = '`+store.Name()+`'`,
	)
}
