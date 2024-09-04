package v1

import (
	ledger "github.com/formancehq/ledger/internal"
	"github.com/formancehq/stack/libs/go-libs/platform/postgres"
	"net/http"

	"github.com/formancehq/ledger/internal/api/backend"
	sharedapi "github.com/formancehq/stack/libs/go-libs/api"
	"github.com/go-chi/chi/v5"
)

func autoCreateMiddleware(backend backend.Backend) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ledgerName := chi.URLParam(r, "ledger")
			if _, err := backend.GetLedger(r.Context(), ledgerName); err != nil {
				if !postgres.IsNotFoundError(err) {
					sharedapi.InternalServerError(w, r, err)
					return
				}

				if err := backend.CreateLedger(r.Context(), ledgerName, ledger.Configuration{
					Bucket: ledgerName,
				}); err != nil {
					sharedapi.InternalServerError(w, r, err)
					return
				}
			}

			handler.ServeHTTP(w, r)
		})
	}
}
