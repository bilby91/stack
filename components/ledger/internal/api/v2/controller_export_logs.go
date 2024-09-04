package v2

import (
	"context"
	"encoding/json"
	ledgercontroller "github.com/formancehq/ledger/internal/controller/ledger"
	"net/http"

	ledger "github.com/formancehq/ledger/internal"
	"github.com/formancehq/ledger/internal/api/backend"
	"github.com/formancehq/stack/libs/go-libs/api"
)

func exportLogs(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/octet-stream")
	if err := backend.LedgerFromContext(r.Context()).Export(r.Context(), ledgercontroller.ExportWriterFn(func(ctx context.Context, log *ledger.ChainedLog) error {
		return enc.Encode(log)
	})); err != nil {
		api.InternalServerError(w, r, err)
		return
	}
}
