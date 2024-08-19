package v2

import (
	"encoding/json"
	"net/http"

	"github.com/formancehq/payments/internal/api/backend"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/api"
)

func paymentsUpdateMetadata(backend backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO(polo): add span
		ctx := r.Context()

		id, err := models.PaymentIDFromString(paymentID(r))
		if err != nil {
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		var metadata map[string]string
		err = json.NewDecoder(r.Body).Decode(&metadata)
		if err != nil {
			api.BadRequest(w, ErrValidation, err)
			return
		}

		err = backend.PaymentsUpdateMetadata(ctx, id, metadata)
		if err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		api.NoContent(w)
	}
}
