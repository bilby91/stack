package v3

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/formancehq/payments/internal/api/backend"
	"github.com/formancehq/stack/libs/go-libs/api"
)

type connectorsInstallRequest struct {
	Provider string `json:"provider"`
}

func (request connectorsInstallRequest) validate() error {
	if request.Provider == "" {
		return errors.New("provider is required")
	}

	return nil
}

func connectorsInstall(backend backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO(polo): add open telemetry span
		ctx := r.Context()

		config, err := io.ReadAll(r.Body)
		if err != nil {
			api.BadRequest(w, ErrMissingOrInvalidBody, err)
			return
		}

		var request connectorsInstallRequest
		if err := json.Unmarshal(config, &request); err != nil {
			api.BadRequest(w, ErrMissingOrInvalidBody, err)
			return
		}

		if err := request.validate(); err != nil {
			api.BadRequest(w, ErrMissingOrInvalidBody, err)
			return
		}

		connectorID, err := backend.ConnectorsInstall(ctx, request.Provider, config)
		if err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		api.Created(w, connectorID)
	}
}
