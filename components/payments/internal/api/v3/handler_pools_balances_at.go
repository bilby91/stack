package v3

import (
	"net/http"
	"time"

	"github.com/formancehq/payments/internal/api/backend"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func poolsBalancesAt(backend backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO(polo): add span
		ctx := r.Context()

		id, err := uuid.Parse(poolID(r))
		if err != nil {
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		atTime := r.URL.Query().Get("at")
		if atTime == "" {
			api.BadRequest(w, ErrValidation, errors.New("missing atTime"))
			return
		}

		at, err := time.Parse(time.RFC3339, atTime)
		if err != nil {
			api.BadRequest(w, ErrValidation, errors.Wrap(err, "invalid atTime"))
			return
		}

		balances, err := backend.PoolsBalancesAt(ctx, id, at)
		if err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		api.Ok(w, balances)
	}
}
