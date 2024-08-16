package v2

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/formancehq/paymentsv3/internal/api/backend"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/contextutil"
)

func connectorsInstall(backend backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO(polo): add open telemetry span
		ctx := r.Context()

		provider := strings.ToLower(connectorProvider(r))
		if provider == "" {
			api.BadRequest(w, ErrValidation, errors.New("provider is required"))
			return
		}

		config, err := io.ReadAll(r.Body)
		if err != nil {
			api.BadRequest(w, ErrMissingOrInvalidBody, err)
			return
		}

		// Detach the context to avoid cancellation of the installation process
		// leading to a partial installation
		ctx, _ = contextutil.Detached(ctx)
		connectorID, err := backend.ConnectorsInstall(ctx, provider, config)
		if err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		api.Created(w, connectorID)
	}
}
