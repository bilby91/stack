package v3

import (
	"net/http"

	"github.com/formancehq/paymentsv3/internal/api/backend"
)

func installConnector(backend backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
