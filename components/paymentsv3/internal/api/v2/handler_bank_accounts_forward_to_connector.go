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

type bankAccountsForwardToConnectorRequest struct {
	ConnectorID string `json:"connectorID"`
}

func (f *bankAccountsForwardToConnectorRequest) Validate() error {
	if f.ConnectorID == "" {
		return errors.New("connectorID must be provided")
	}

	return nil
}

func bankAccountsForwardToConnector(backend backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO(polo): add span
		ctx := r.Context()

		id, err := uuid.Parse(bankAccountID(r))
		if err != nil {
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		var req bankAccountsForwardToConnectorRequest
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			api.BadRequest(w, ErrMissingOrInvalidBody, err)
			return
		}

		err = req.Validate()
		if err != nil {
			api.BadRequest(w, ErrMissingOrInvalidBody, err)
			return
		}

		connectorID, err := models.ConnectorIDFromString(req.ConnectorID)
		if err != nil {
			api.BadRequest(w, ErrValidation, err)
			return
		}

		bankAccount, err := backend.BankAccountsForwardToConnector(ctx, id, connectorID)
		if err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		data := &bankAccountResponse{
			ID:        bankAccount.ID.String(),
			Name:      bankAccount.Name,
			CreatedAt: bankAccount.CreatedAt,
			Metadata:  bankAccount.Metadata,
		}

		if bankAccount.IBAN != nil {
			data.Iban = *bankAccount.IBAN
		}

		if bankAccount.AccountNumber != nil {
			data.AccountNumber = *bankAccount.AccountNumber
		}

		if bankAccount.SwiftBicCode != nil {
			data.SwiftBicCode = *bankAccount.SwiftBicCode
		}

		if bankAccount.Country != nil {
			data.Country = *bankAccount.Country
		}

		for _, relatedAccount := range bankAccount.RelatedAccounts {
			data.RelatedAccounts = append(data.RelatedAccounts, &bankAccountRelatedAccountsResponse{
				ID:          "",
				CreatedAt:   relatedAccount.CreatedAt,
				AccountID:   relatedAccount.AccountID.String(),
				ConnectorID: relatedAccount.ConnectorID.String(),
				Provider:    relatedAccount.ConnectorID.Provider,
			})
		}

		err = json.NewEncoder(w).Encode(api.BaseResponse[bankAccountResponse]{
			Data: data,
		})
		if err != nil {
			api.InternalServerError(w, r, err)
			return
		}
	}
}
