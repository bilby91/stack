package v2

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/formancehq/paymentsv3/internal/api/backend"
	"github.com/formancehq/paymentsv3/internal/models"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/google/uuid"
)

type createPoolRequest struct {
	Name       string   `json:"name"`
	AccountIDs []string `json:"accountIDs"`
}

type poolResponse struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Accounts []string `json:"accounts"`
}

func poolsCreate(backend backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO(polo): add span
		ctx := r.Context()

		var createPoolRequest createPoolRequest
		err := json.NewDecoder(r.Body).Decode(&createPoolRequest)
		if err != nil {
			api.BadRequest(w, ErrMissingOrInvalidBody, err)
			return
		}

		pool := models.Pool{
			ID:        uuid.New(),
			Name:      createPoolRequest.Name,
			CreatedAt: time.Now().UTC(),
		}

		accounts := make([]models.PoolAccounts, len(createPoolRequest.AccountIDs))
		for i, accountID := range createPoolRequest.AccountIDs {
			aID, err := models.AccountIDFromString(accountID)
			if err != nil {
				api.BadRequest(w, ErrValidation, err)
				return
			}

			accounts[i] = models.PoolAccounts{
				PoolID:    pool.ID,
				AccountID: aID,
			}
		}
		pool.PoolAccounts = accounts

		err = backend.PoolsCreate(ctx, pool)
		if err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		data := &poolResponse{
			ID:       pool.ID.String(),
			Name:     pool.Name,
			Accounts: createPoolRequest.AccountIDs,
		}

		err = json.NewEncoder(w).Encode(api.BaseResponse[poolResponse]{
			Data: data,
		})
		if err != nil {
			api.InternalServerError(w, r, err)
			return
		}
	}
}
