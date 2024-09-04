package v2

import (
	"github.com/formancehq/ledger/internal/controller/ledger/writer"
	"net/http"
	"strings"

	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"

	"github.com/pkg/errors"
)

const (
	MaxPageSize     = bunpaginate.MaxPageSize
	DefaultPageSize = bunpaginate.QueryDefaultPageSize

	QueryKeyCursor = "cursor"
)

var (
	ErrInvalidBalanceOperator = errors.New(
		"invalid parameter 'balanceOperator', should be one of 'e, ne, gt, gte, lt, lte'")
	ErrInvalidStartTime = errors.New("invalid 'startTime' query param")
	ErrInvalidEndTime   = errors.New("invalid 'endTime' query param")
)

func getCommandParameters(r *http.Request) writer.Parameters {
	dryRunAsString := r.URL.Query().Get("dryRun")
	dryRun := strings.ToUpper(dryRunAsString) == "YES" || strings.ToUpper(dryRunAsString) == "TRUE" || dryRunAsString == "1"

	return writer.Parameters{
		DryRun:         dryRun,
		IdempotencyKey: api.IdempotencyKeyFromRequest(r),
	}
}
