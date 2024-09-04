package v1

import (
	"context"
	_ "embed"
	ledger "github.com/formancehq/ledger/internal"
	systemcontroller "github.com/formancehq/ledger/internal/controller/system"
	"net/http"

	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"

	"github.com/formancehq/stack/libs/go-libs/collectionutils"

	"github.com/formancehq/ledger/internal/api/backend"
	sharedapi "github.com/formancehq/stack/libs/go-libs/api"
)

type ConfigInfo struct {
	Server  string        `json:"server"`
	Version string        `json:"version"`
	Config  *LedgerConfig `json:"config"`
}

type LedgerConfig struct {
	LedgerStorage *LedgerStorage `json:"storage"`
}

type LedgerStorage struct {
	Driver  string   `json:"driver"`
	Ledgers []string `json:"ledgers"`
}

func getInfo(backend backend.Backend) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		ledgerNames := make([]string, 0)
		if err := bunpaginate.Iterate(r.Context(), systemcontroller.NewListLedgersQuery(100),
			func(ctx context.Context, q systemcontroller.ListLedgersQuery) (*bunpaginate.Cursor[ledger.Ledger], error) {
				return backend.ListLedgers(ctx, q)
			},
			func(cursor *bunpaginate.Cursor[ledger.Ledger]) error {
				ledgerNames = append(ledgerNames, collectionutils.Map(cursor.Data, func(from ledger.Ledger) string {
					return from.Name
				})...)
				return nil
			},
		); err != nil {
			sharedapi.InternalServerError(w, r, err)
			return
		}

		sharedapi.Ok(w, ConfigInfo{
			Server:  "ledger",
			Version: backend.GetVersion(),
			Config: &LedgerConfig{
				LedgerStorage: &LedgerStorage{
					Driver:  "postgres",
					Ledgers: ledgerNames,
				},
			},
		})
	}
}
