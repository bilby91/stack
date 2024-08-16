package v2

import (
	"encoding/json"
	"math/big"
	"net/http"
	"time"

	"github.com/formancehq/paymentsv3/internal/api/backend"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// NOTE: in order to maintain previous version compatibility, we need to keep the
// same response structure as the previous version of the API
type poolBalancesResponse struct {
	Balances []*poolBalanceResponse `json:"balances"`
}

type poolBalanceResponse struct {
	Amount *big.Int `json:"amount"`
	Asset  string   `json:"asset"`
}

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

		data := &poolBalancesResponse{
			Balances: make([]*poolBalanceResponse, len(balances)),
		}

		for i := range balances {
			data.Balances[i] = &poolBalanceResponse{
				Amount: balances[i].Amount,
				Asset:  balances[i].Asset,
			}
		}

		err = json.NewEncoder(w).Encode(api.BaseResponse[poolBalancesResponse]{
			Data: data,
		})
		if err != nil {
			api.InternalServerError(w, r, err)
			return
		}
	}
}
