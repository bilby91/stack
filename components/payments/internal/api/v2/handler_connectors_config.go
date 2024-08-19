package v2

import (
	"net/http"

	"github.com/formancehq/payments/internal/api/backend"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/api"
)

func connectorsConfig(backend backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO(polo): add open telemetry span
		ctx := r.Context()

		connectorID, err := models.ConnectorIDFromString(connectorID(r))
		if err != nil {
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		config, err := backend.ConnectorsConfig(ctx, connectorID)
		if err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		api.Ok(w, config)
	}
}
