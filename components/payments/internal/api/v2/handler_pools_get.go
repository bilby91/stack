package v2

import (
	"encoding/json"
	"net/http"

	"github.com/formancehq/payments/internal/api/backend"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/google/uuid"
)

func poolsGet(backend backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO(polo): add span
		ctx := r.Context()

		id, err := uuid.Parse(poolID(r))
		if err != nil {
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		pool, err := backend.PoolsGet(ctx, id)
		if err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		data := &poolResponse{
			ID:   pool.ID.String(),
			Name: pool.Name,
		}

		accounts := make([]string, len(pool.PoolAccounts))
		for i := range pool.PoolAccounts {
			accounts[i] = pool.PoolAccounts[i].AccountID.String()
		}
		data.Accounts = accounts

		err = json.NewEncoder(w).Encode(api.BaseResponse[poolResponse]{
			Data: data,
		})
		if err != nil {
			api.InternalServerError(w, r, err)
			return
		}
	}
}
