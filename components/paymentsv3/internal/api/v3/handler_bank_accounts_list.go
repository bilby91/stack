package v3

import (
	"net/http"

	"github.com/formancehq/paymentsv3/internal/api/backend"
	"github.com/formancehq/paymentsv3/internal/storage"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/formancehq/stack/libs/go-libs/pointer"
)

func bankAccountsList(backend backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO(polo): add span
		ctx := r.Context()

		query, err := bunpaginate.Extract[storage.ListBankAccountsQuery](r, func() (*storage.ListBankAccountsQuery, error) {
			options, err := getPagination(r, storage.BankAccountQuery{})
			if err != nil {
				return nil, err
			}
			return pointer.For(storage.NewListBankAccountsQuery(*options)), nil
		})
		if err != nil {
			api.BadRequest(w, ErrValidation, err)
			return
		}

		cursor, err := backend.BankAccountsList(ctx, *query)
		if err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		api.RenderCursor(w, *cursor)
	}
}
