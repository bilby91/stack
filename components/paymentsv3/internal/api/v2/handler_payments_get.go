package v2

import (
	"encoding/json"
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

		data := paymentResponse{
			ID:            payment.ID.String(),
			Reference:     payment.Reference,
			Type:          payment.Type.String(),
			Provider:      payment.ConnectorID.Provider,
			ConnectorID:   payment.ConnectorID.String(),
			Status:        payment.Status.String(),
			Amount:        payment.Amount,
			InitialAmount: payment.InitialAmount,
			Scheme:        payment.Scheme.String(),
			Asset:         payment.Asset,
			CreatedAt:     payment.CreatedAt,
			Metadata:      payment.Metadata,
		}

		if payment.SourceAccountID != nil {
			data.SourceAccountID = payment.SourceAccountID.String()
		}

		if payment.DestinationAccountID != nil {
			data.DestinationAccountID = payment.DestinationAccountID.String()
		}

		data.Adjustments = make([]paymentAdjustment, len(payment.Adjustments))
		for i := range payment.Adjustments {
			data.Adjustments[i] = paymentAdjustment{
				Reference: payment.Adjustments[i].ID.Reference,
				CreatedAt: payment.Adjustments[i].CreatedAt,
				Status:    payment.Adjustments[i].Status.String(),
				Amount:    payment.Adjustments[i].Amount,
				Raw:       payment.Adjustments[i].Raw,
			}
		}

		err = json.NewEncoder(w).Encode(api.BaseResponse[paymentResponse]{
			Data: &data,
		})
		if err != nil {
			api.InternalServerError(w, r, err)
			return
		}
	}
}
