//go:build it

package performance_test

import (
	"context"
	"fmt"
	ledger "github.com/formancehq/ledger/internal"
	"github.com/formancehq/ledger/internal/controller/ledger/writer"
	"github.com/formancehq/ledger/internal/storage/bucket"
	ledgerstore "github.com/formancehq/ledger/internal/storage/ledger"
	"github.com/formancehq/ledger/pkg/testserver"
	"github.com/formancehq/stack/ledger/client/models/components"
	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"
	"github.com/formancehq/stack/libs/go-libs/bun/bundebug"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"io"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func BenchmarkSequentialWorldToBank(b *testing.B) {
	w := newWriter(b)
	ctx := logging.TestingContext()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := w.CreateTransaction(ctx, writer.Parameters{}, ledger.RunScript{
			Script: ledger.Script{
				Plain: `
send [USD/2 100] (
	source = @world
	destination = @bank
)`,
			},
		})
		require.NoError(b, err)
	}
}

func BenchmarkWrite(b *testing.B) {

	type testCase struct {
		name          string
		scriptFactory func(int) string
	}

	for _, tc := range []testCase{
		{
			name: "world to bank",
			scriptFactory: func(_ int) string {
				return `
send [USD/2 100] (
	source = @world
	destination = @bank
)`
			},
		},
		{
			name: "world to not existing destination",
			scriptFactory: func(id int) string {
				return fmt.Sprintf(`
send [USD/2 100] (
	source = @world
	destination = @dst:%d
)`, id)
			},
		},
		{
			name: "not existing source to not existing destination",
			scriptFactory: func(id int) string {
				return fmt.Sprintf(`
send [USD/2 100] (
	source = @src:%d allowing unbounded overdraft
	destination = @dst:%d
)`, id, id)
			},
		},
	} {
		b.Run(tc.name, func(b *testing.B) {
			runParallelBenchmark(b, tc.scriptFactory)
		})
	}
}

type benchmark func(int) string

type report struct {
	longestTxLock              sync.Mutex
	longestTransactionID       int
	longestTransactionDuration time.Duration
	startOfBench               time.Time
	totalDuration              atomic.Int64
}

func runParallelBenchmark(b *testing.B, fn benchmark) {
	b.Helper()

	cpt := atomic.Int64{}

	type env struct {
		name    string
		factory func(b *testing.B) Env
	}

	for _, envFactory := range []env{
		{
			name:    "testserver",
			factory: NewTestServerEnv,
		},
		{
			name:    "core",
			factory: NewCoreEnv,
		},
	} {
		b.Run(envFactory.name, func(b *testing.B) {
			env := envFactory.factory(b)

			report := &report{
				startOfBench: time.Now(),
			}

			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					script := fn(int(cpt.Add(1)))
					now := time.Now()
					transaction, err := env.Executor().ExecuteScript(script)
					require.NoError(b, err)

					latency := time.Since(now)
					report.totalDuration.Add(latency.Milliseconds())

					report.longestTxLock.Lock()
					if latency > report.longestTransactionDuration {
						report.longestTransactionID = transaction.ID
						report.longestTransactionDuration = latency
					}
					report.longestTxLock.Unlock()
				}
			})
			b.StopTimer()

			b.Logf("Longest transaction: %d (%s)", report.longestTransactionID, report.longestTransactionDuration.String())
			b.ReportMetric((float64(time.Duration(b.N))/float64(time.Since(report.startOfBench)))*float64(time.Second), "t/s")
			b.ReportMetric(float64(report.totalDuration.Load()/int64(b.N)), "ms/transaction")
		})
	}
}

func newWriter(b *testing.B) *writer.Writer {
	b.Helper()

	ctx := logging.TestingContext()

	pgDatabase := pgServer.NewDatabase(b)

	hooks := make([]bun.QueryHook, 0)
	if testing.Verbose() {
		hooks = append(hooks, bundebug.NewQueryHook())
	}

	connectionOptions := pgDatabase.ConnectionOptions()
	connectionOptions.ConnMaxIdleTime = time.Minute
	connectionOptions.MaxOpenConns = 100
	connectionOptions.MaxIdleConns = 100

	bunDB, err := bunconnect.OpenSQLDB(ctx, connectionOptions, hooks...)
	require.NoError(b, err)

	bucket := bucket.New(bunDB, "_default")
	require.NoError(b, bucket.Migrate(ctx))
	require.NoError(b, ledgerstore.Migrate(ctx, bunDB, "_default", "benchmark"))

	ledgerStore := ledgerstore.NewDefaultStoreAdapter(
		ledgerstore.New(bunDB, "_default", "benchmark"),
	)
	machineFactory := writer.NewDefaultMachineFactory(
		writer.NewCachedCompiler(
			writer.NewDefaultCompiler(),
			writer.CacheConfiguration{
				MaxCount: 10,
			},
		),
		ledgerStore,
	)
	return writer.New(ledgerStore, machineFactory)
}

type TransactionExecutor interface {
	ExecuteScript(string) (*ledger.Transaction, error)
}
type TransactionExecutorFn func(string) (*ledger.Transaction, error)

func (fn TransactionExecutorFn) ExecuteScript(script string) (*ledger.Transaction, error) {
	return fn(script)
}

type Env interface {
	Executor() TransactionExecutor
}

type CoreEnv struct {
	w *writer.Writer
}

func (c *CoreEnv) Executor() TransactionExecutor {
	return TransactionExecutorFn(func(plain string) (*ledger.Transaction, error) {
		return c.w.CreateTransaction(context.Background(), writer.Parameters{}, ledger.RunScript{
			Script: ledger.Script{
				Plain: plain,
			},
		})
	})
}

func NewCoreEnv(b *testing.B) Env {
	return &CoreEnv{
		w: newWriter(b),
	}
}

var _ Env = (*CoreEnv)(nil)

type TestServerEnv struct {
	testServer *testserver.Server
}

func (c *TestServerEnv) Executor() TransactionExecutor {
	return TransactionExecutorFn(func(plain string) (*ledger.Transaction, error) {
		ret, err := c.testServer.Client().Ledger.V2.
			CreateTransaction(context.Background(), "_default", components.V2PostTransaction{
				Script: &components.Script{
					Plain: plain,
				},
			}, nil, nil)
		if err != nil {
			return nil, err
		}
		return &ledger.Transaction{
			ID: int(ret.V2CreateTransactionResponse.Data.ID.Int64()),
		}, nil
	})
}

var _ Env = (*TestServerEnv)(nil)

func NewTestServerEnv(b *testing.B) Env {

	connectionOptions := pgServer.NewDatabase(b).ConnectionOptions()
	connectionOptions.MaxOpenConns = 100
	connectionOptions.MaxIdleConns = 100
	connectionOptions.ConnMaxIdleTime = time.Minute

	var output io.Writer = os.Stdout
	if !testing.Verbose() {
		output = io.Discard
	}

	testServer := testserver.New(b, testserver.Configuration{
		PostgresConfiguration: connectionOptions,
		Debug:                 testing.Verbose(),
		Output:                output,
	})
	testServer.Start()

	_, err := testServer.Client().Ledger.V2.
		CreateLedger(context.Background(), "_default", &components.V2CreateLedgerRequest{})
	require.NoError(b, err)

	return &TestServerEnv{
		testServer: testServer,
	}
}
