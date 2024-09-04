package v2

import (
	ledgercontroller "github.com/formancehq/ledger/internal/controller/ledger"
	"net/http"

	"github.com/formancehq/ledger/internal/api/backend"
	sharedapi "github.com/formancehq/stack/libs/go-libs/api"

	"github.com/formancehq/stack/libs/go-libs/pointer"

	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
)

func getVolumesWithBalances(w http.ResponseWriter, r *http.Request) {

	l := backend.LedgerFromContext(r.Context())

	query, err := bunpaginate.Extract[ledgercontroller.GetVolumesWithBalancesQuery](r, func() (*ledgercontroller.GetVolumesWithBalancesQuery, error) {
		options, err := getPaginatedQueryOptionsOfFiltersForVolumes(r)
		if err != nil {
			return nil, err
		}

		getVolumesWithBalancesQuery := ledgercontroller.NewGetVolumesWithBalancesQuery(*options)
		return pointer.For(getVolumesWithBalancesQuery), nil

	})

	if err != nil {
		sharedapi.BadRequest(w, ErrValidation, err)
		return
	}

	cursor, err := l.GetVolumesWithBalances(r.Context(), *query)

	if err != nil {
		switch {
		//case ledger.IsErrInvalidQuery(err):
		//	sharedapi.BadRequest(w, ErrValidation, err)
		default:
			sharedapi.InternalServerError(w, r, err)
		}
		return
	}

	sharedapi.RenderCursor(w, *cursor)

}
