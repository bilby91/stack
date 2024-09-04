package v2

import (
	"fmt"
	ledgercontroller "github.com/formancehq/ledger/internal/controller/ledger"
	"net/http"

	"github.com/formancehq/ledger/internal/api/backend"
	sharedapi "github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
)

func getLogs(w http.ResponseWriter, r *http.Request) {
	l := backend.LedgerFromContext(r.Context())

	query := ledgercontroller.GetLogsQuery{}

	if r.URL.Query().Get(QueryKeyCursor) != "" {
		err := bunpaginate.UnmarshalCursor(r.URL.Query().Get(QueryKeyCursor), &query)
		if err != nil {
			sharedapi.BadRequest(w, ErrValidation, fmt.Errorf("invalid '%s' query param", QueryKeyCursor))
			return
		}
	} else {
		var err error

		pageSize, err := bunpaginate.GetPageSize(r)
		if err != nil {
			sharedapi.BadRequest(w, ErrValidation, err)
			return
		}

		qb, err := getQueryBuilder(r)
		if err != nil {
			sharedapi.BadRequest(w, ErrValidation, err)
			return
		}

		query = ledgercontroller.NewGetLogsQuery(ledgercontroller.PaginatedQueryOptions[any]{
			QueryBuilder: qb,
			PageSize:     pageSize,
		})
	}

	cursor, err := l.GetLogs(r.Context(), query)
	if err != nil {
		sharedapi.InternalServerError(w, r, err)
		return
	}

	sharedapi.RenderCursor(w, *cursor)
}
