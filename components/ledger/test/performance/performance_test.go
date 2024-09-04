//go:build it

package performance_test

import (
	"bytes"
	"fmt"
	"github.com/formancehq/ledger/pkg/testserver"
	"github.com/formancehq/stack/ledger/client/models/components"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"github.com/formancehq/stack/libs/go-libs/time"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"math/big"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

func BenchmarkWorstCase(b *testing.B) {

	db := pgServer.NewDatabase(b)

	ctx := logging.TestingContext()

	ledgerName := uuid.NewString()
	connectionOptions := db.ConnectionOptions()
	connectionOptions.MaxOpenConns = 20
	connectionOptions.MaxIdleConns = 20
	connectionOptions.ConnMaxIdleTime = time.Minute

	testServer := testserver.New(b, testserver.Configuration{
		PostgresConfiguration: connectionOptions,
		Debug:                 testing.Verbose(),
		Output:                os.Stdout,
	})
	testServer.Start()
	defer testServer.Stop()

	_, err := testServer.Client().Ledger.V2.CreateLedger(ctx, ledgerName, &components.V2CreateLedgerRequest{})
	require.NoError(b, err)

	totalDuration := atomic.Int64{}
	runtime.GC()
	b.ResetTimer()
	startOfBench := time.Now()
	counter := atomic.Int64{}
	longestTxLock := sync.Mutex{}
	longestTransactionID := big.NewInt(0)
	longestTransactionDuration := time.Duration(0)
	b.RunParallel(func(pb *testing.PB) {
		buf := bytes.NewBufferString("")
		for pb.Next() {
			buf.Reset()
			id := counter.Add(1)
			now := time.Now()

			// todo: check why the generated sdk does not have the same signature as the global sdk
			transactionResponse, err := testServer.Client().Ledger.V2.CreateTransaction(ctx, ledgerName, components.V2PostTransaction{
				Timestamp: nil,
				Postings:  nil,
				Script: &components.Script{
					Plain: `vars {
	account $account
}

send [USD/2 100] (
	source = @world
	destination = $account
)`,
					Vars: map[string]any{
						"account": fmt.Sprintf("accounts:%d", id),
					},
				},

				Reference: nil,
				Metadata:  nil,
			}, pointer.For(false), nil)
			if err != nil {
				continue
			}

			latency := time.Since(now).Milliseconds()
			totalDuration.Add(latency)

			longestTxLock.Lock()
			if time.Millisecond*time.Duration(latency) > longestTransactionDuration {
				longestTransactionID = transactionResponse.V2CreateTransactionResponse.Data.ID
				longestTransactionDuration = time.Duration(latency) * time.Millisecond
			}
			longestTxLock.Unlock()
		}
	})

	b.StopTimer()
	b.Logf("Longest transaction: %d (%s)", longestTransactionID, longestTransactionDuration.String())
	b.ReportMetric((float64(time.Duration(b.N))/float64(time.Since(startOfBench)))*float64(time.Second), "t/s")
	b.ReportMetric(float64(totalDuration.Load()/int64(b.N)), "ms/transaction")

	runtime.GC()
}
