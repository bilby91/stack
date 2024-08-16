package v3

import (
	"net/http"

	"github.com/formancehq/paymentsv3/internal/api/backend"
	"github.com/formancehq/paymentsv3/internal/models"
	"github.com/formancehq/stack/libs/go-libs/api"
)

func paymentsGet(backend backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO(polo): add span
		ctx := r.Context()

		id, err := models.PaymentIDFromString(paymentID(r))
		if err != nil {
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		payment, err := backend.PaymentsGet(ctx, id)
		if err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		api.Ok(w, payment)
	}
}
