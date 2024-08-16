package v2

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/formancehq/paymentsv3/internal/api/backend"
	"github.com/formancehq/paymentsv3/internal/models"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/google/uuid"
)

type poolsAddAccountRequest struct {
	AccountID string `json:"accountID"`
}

func (c *poolsAddAccountRequest) Validate() error {
	if c.AccountID == "" {
		return errors.New("accountID is required")
	}

	return nil
}

func poolsAddAccount(backend backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO(polo): add span
		ctx := r.Context()

		id, err := uuid.Parse(poolID(r))
		if err != nil {
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		var poolsAddAccountRequest poolsAddAccountRequest
		err = json.NewDecoder(r.Body).Decode(&poolsAddAccountRequest)
		if err != nil {
			api.BadRequest(w, ErrMissingOrInvalidBody, err)
			return
		}

		if err := poolsAddAccountRequest.Validate(); err != nil {
			api.BadRequest(w, ErrValidation, err)
			return
		}

		accountID, err := models.AccountIDFromString(poolsAddAccountRequest.AccountID)
		if err != nil {
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		err = backend.PoolsAddAccount(ctx, id, accountID)
		if err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		api.NoContent(w)
	}
}
