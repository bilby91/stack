package v2

import (
	"net/http"

	"github.com/formancehq/paymentsv3/internal/api/backend"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/google/uuid"
)

func poolsDelete(backend backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO(polo): add span
		ctx := r.Context()

		id, err := uuid.Parse(poolID(r))
		if err != nil {
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		err = backend.PoolsDelete(ctx, id)
		if err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		api.NoContent(w)
	}
}
