package sqlstorage_test

import (
	"context"
	"fmt"
	"github.com/huandu/go-sqlbuilder"
	"github.com/numary/ledger/internal/pgtesting"
	"github.com/numary/ledger/pkg/core"
	"github.com/numary/ledger/pkg/ledger/query"
	"github.com/numary/ledger/pkg/storage/sqlstorage"
	"github.com/pborman/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func BenchmarkStore(b *testing.B) {

	if testing.Verbose() {
		logrus.StandardLogger().Level = logrus.DebugLevel
	}

	pgServer, err := pgtesting.PostgresServer()
	assert.NoError(b, err)
	defer pgServer.Close()

	type driverConfig struct {
		driver    string
		dbFactory func() (sqlstorage.DB, error)
		flavor    sqlbuilder.Flavor
	}
	var drivers = []driverConfig{
		{
			driver: "sqlite3",
			dbFactory: func() (sqlstorage.DB, error) {
				return sqlstorage.NewSQLiteDB(os.TempDir(), uuid.New()), nil
			},
			flavor: sqlbuilder.SQLite,
		},
		{
			driver: "pgx",
			dbFactory: func() (sqlstorage.DB, error) {
				db, err := sqlstorage.OpenSQLDB(sqlstorage.PostgreSQL, pgServer.ConnString())
				if err != nil {
					return nil, err
				}
				return sqlstorage.NewPostgresDB(db), nil
			},
			flavor: sqlbuilder.PostgreSQL,
		},
	}

	type testingFunction struct {
		name string
		fn   func(b *testing.B, store *sqlstorage.Store)
	}

	for _, driver := range drivers {
		for _, tf := range []testingFunction{
			{
				name: "FindTransactions",
				fn:   testBenchmarkFindTransactions,
			},
			{
				name: "LastTransaction",
				fn:   testBenchmarkLastTransaction,
			},
			{
				name: "AggregateVolumes",
				fn:   testBenchmarkAggregateVolumes,
			},
		} {
			b.Run(fmt.Sprintf("%s/%s", driver.driver, tf.name), func(b *testing.B) {
				ledger := uuid.New()

				db, err := driver.dbFactory()
				if !assert.NoError(b, err) {
					return
				}

				schema, err := db.Schema(context.Background(), uuid.New())
				if !assert.NoError(b, err) {
					return
				}

				store, err := sqlstorage.NewStore(ledger, driver.flavor, schema, func(ctx context.Context) error {
					return db.Close(context.Background())
				})
				assert.NoError(b, err)
				defer store.Close(context.Background())

				_, err = store.Initialize(context.Background())
				assert.NoError(b, err)

				b.ResetTimer()

				tf.fn(b, store)
			})
		}
	}
}

func testBenchmarkFindTransactions(b *testing.B, store *sqlstorage.Store) {
	datas := make([]core.Transaction, 0)
	for i := 0; i < 1000; i++ {
		datas = append(datas, core.Transaction{
			TransactionData: core.TransactionData{
				Postings: []core.Posting{
					{
						Source:      "world",
						Destination: fmt.Sprintf("player%d", i),
						Asset:       "USD",
						Amount:      100,
					},
					{
						Source:      "world",
						Destination: fmt.Sprintf("player%d", i+1),
						Asset:       "USD",
						Amount:      100,
					},
				},
			},
			ID: int64(i),
		})
	}

	_, err := store.SaveTransactions(context.Background(), datas)
	assert.NoError(b, err)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		txs, err := store.FindTransactions(context.Background(), query.Query{
			Limit: 100,
		})
		assert.NoError(b, err)
		if txs.PageSize != 100 {
			b.Errorf("Should have 100 transactions but get %d", txs.PageSize)
		}
	}

}

func testBenchmarkLastTransaction(b *testing.B, store *sqlstorage.Store) {
	datas := make([]core.Transaction, 0)
	count := 1000
	for i := 0; i < count; i++ {
		datas = append(datas, core.Transaction{
			TransactionData: core.TransactionData{
				Postings: []core.Posting{
					{
						Source:      "world",
						Destination: fmt.Sprintf("player%d", i),
						Asset:       "USD",
						Amount:      100,
					},
					{
						Source:      "world",
						Destination: fmt.Sprintf("player%d", i+1),
						Asset:       "USD",
						Amount:      100,
					},
				},
			},
			ID: int64(i),
		})
	}

	_, err := store.SaveTransactions(context.Background(), datas)
	assert.NoError(b, err)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		tx, err := store.LastTransaction(context.Background())
		assert.NoError(b, err)
		assert.Equal(b, int64(count-1), tx.ID)
	}

}

func testBenchmarkAggregateVolumes(b *testing.B, store *sqlstorage.Store) {
	datas := make([]core.Transaction, 0)
	count := 1000
	for i := 0; i < count; i++ {
		datas = append(datas, core.Transaction{
			TransactionData: core.TransactionData{
				Postings: []core.Posting{
					{
						Source:      "world",
						Destination: fmt.Sprintf("player%d", i),
						Asset:       "USD",
						Amount:      100,
					},
					{
						Source:      "world",
						Destination: fmt.Sprintf("player%d", i+1),
						Asset:       "USD",
						Amount:      100,
					},
					{
						Source:      fmt.Sprintf("player%d", i),
						Destination: fmt.Sprintf("player%d", i+1),
						Asset:       "USD",
						Amount:      50,
					},
				},
			},
			ID: int64(i),
		})
	}

	_, err := store.SaveTransactions(context.Background(), datas)
	assert.NoError(b, err)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := store.AggregateVolumes(context.Background(), "world")
		assert.NoError(b, err)
	}

}
