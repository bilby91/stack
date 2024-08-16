package v3

import (
	"net/http"

	"github.com/formancehq/paymentsv3/internal/api/backend"
	"github.com/formancehq/paymentsv3/internal/models"
	"github.com/formancehq/stack/libs/go-libs/api"
)

func connectorsReset(backend backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO(polo): add open telemetry span
		ctx := r.Context()

		connectorID, err := models.ConnectorIDFromString(connectorID(r))
		if err != nil {
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		if err := backend.ConnectorsReset(ctx, connectorID); err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		api.NoContent(w)
	}
}
