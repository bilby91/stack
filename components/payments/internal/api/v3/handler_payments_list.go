package v3

import (
	"net/http"

	"github.com/formancehq/payments/internal/api/backend"
	"github.com/formancehq/payments/internal/storage"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/formancehq/stack/libs/go-libs/pointer"
)

func paymentsList(backend backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO(polo): add span
		ctx := r.Context()

		query, err := bunpaginate.Extract[storage.ListPaymentsQuery](r, func() (*storage.ListPaymentsQuery, error) {
			options, err := getPagination(r, storage.PaymentQuery{})
			if err != nil {
				return nil, err
			}
			return pointer.For(storage.NewListPaymentsQuery(*options)), nil
		})
		if err != nil {
			api.BadRequest(w, ErrValidation, err)
			return
		}

		cursor, err := backend.PaymentsList(ctx, *query)
		if err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		api.RenderCursor(w, *cursor)
	}
}
