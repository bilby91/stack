package v3

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/formancehq/paymentsv3/internal/api/backend"
	"github.com/formancehq/paymentsv3/internal/models"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/google/uuid"
)

type bankAccountsCreateRequest struct {
	Name string `json:"name"`

	AccountNumber *string `json:"accountNumber"`
	IBAN          *string `json:"iban"`
	SwiftBicCode  *string `json:"swiftBicCode"`
	Country       *string `json:"country"`

	Metadata map[string]string `json:"metadata"`
}

func (r *bankAccountsCreateRequest) Validate() error {
	if r.AccountNumber == nil && r.IBAN == nil {
		return errors.New("either accountNumber or iban must be provided")
	}

	if r.Name == "" {
		return errors.New("name must be provided")
	}

	return nil
}

func bankAccountsCreate(backend backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO(polo): add open telemetry span
		ctx := r.Context()

		var req bankAccountsCreateRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			api.BadRequest(w, ErrMissingOrInvalidBody, err)
			return
		}

		if err := req.Validate(); err != nil {
			api.BadRequest(w, ErrValidation, err)
			return
		}

		bankAccount := &models.BankAccount{
			ID:            uuid.New(),
			CreatedAt:     time.Now().UTC(),
			Name:          req.Name,
			AccountNumber: req.AccountNumber,
			IBAN:          req.IBAN,
			SwiftBicCode:  req.SwiftBicCode,
			Country:       req.Country,
			Metadata:      req.Metadata,
		}

		err = backend.BankAccountsCreate(ctx, *bankAccount)
		if err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		api.Created(w, bankAccount.ID.String())
	}
}
