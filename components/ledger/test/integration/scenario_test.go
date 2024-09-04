//go:build it

package test_suite

import (
	. "github.com/formancehq/ledger/pkg/testserver"
	"github.com/formancehq/stack/ledger/client/models/components"
	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"
	"github.com/formancehq/stack/libs/go-libs/logging"
	. "github.com/formancehq/stack/libs/go-libs/testing/platform/pgtesting"
	"github.com/formancehq/stack/libs/go-libs/testing/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/require"
	"math/big"
)

var _ = Context("Ledger integration tests", func() {
	var (
		db  = UsePostgresDatabase(pgServer)
		ctx = logging.TestingContext()
	)

	testServer := UseNewTestServer(func() Configuration {
		return Configuration{
			PostgresConfiguration: db.GetValue().ConnectionOptions(),
			Output:                GinkgoWriter,
			Debug:                 debug,
		}
	})
	When("Starting the service", func() {
		It("Should be ok", func() {
			info, err := testServer.GetValue().Client().Ledger.V2.GetInfo(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(info.V2ConfigInfoResponse.Version).To(Equal("develop"))
		})
	})
	When("Creating a new ledger", func() {
		var ledgerName = "foo"
		BeforeEach(func() {
			_, err := testServer.GetValue().Client().Ledger.V2.CreateLedger(ctx, ledgerName, &components.V2CreateLedgerRequest{})
			Expect(err).To(BeNil())
		})
		It("Should be ok", func() {})
		When("Creating a new transaction", func() {
			FIt("Should be ok", func() {
				_, err := testServer.GetValue().Client().Ledger.V2.CreateTransaction(ctx, ledgerName, components.V2PostTransaction{
					Postings: []components.V2Posting{
						{
							Amount:      big.NewInt(100),
							Asset:       "USD/2",
							Destination: "bank",
							Source:      "world",
						},
					},
				}, nil, nil)
				Expect(err).To(BeNil())

				bunDB, err := bunconnect.OpenSQLDB(ctx, db.GetValue().ConnectionOptions())
				require.NoError(GinkgoT(), err)

				utils.DumpTables(GinkgoT(), ctx, bunDB, `
				SELECT left(query, 20), calls, total_exec_time, rows, 100.0 * shared_blks_hit /
					   nullif(shared_blks_hit + shared_blks_read, 0) AS hit_percent
				  FROM pg_stat_statements ORDER BY total_exec_time`)
			})
		})
	})
})
