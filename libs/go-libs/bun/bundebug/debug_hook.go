package bundebug

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/formancehq/stack/libs/go-libs/logging"

	"github.com/uptrace/bun"
)

type QueryHook struct{}

var _ bun.QueryHook = (*QueryHook)(nil)

func NewQueryHook() *QueryHook {
	return &QueryHook{}
}

func (h *QueryHook) BeforeQuery(
	ctx context.Context, event *bun.QueryEvent,
) context.Context {
	//// todo: maybe use a value in the context to avoid this dirty check
	//if !strings.HasPrefix(event.Query, "select pid") {
	//	rows, err := event.DB.QueryContext(ctx, `select pid, mode, relname, reltype from pg_locks join pg_class on pg_class.oid = pg_locks.relation`)
	//	if err != nil {
	//		panic(err)
	//	}
	//
	//	prettied, err := xsql.Pretty(rows)
	//	if err != nil {
	//		panic(err)
	//	}
	//
	//	_, _ = logging.FromContext(ctx).
	//		Writer().
	//		Write([]byte(prettied))
	//}

	return ctx
}

func (h *QueryHook) AfterQuery(ctx context.Context, event *bun.QueryEvent) {
	dur := time.Since(event.StartTime)

	fields := map[string]any{
		"component": "bun",
		"operation": event.Operation(),
		"duration":  fmt.Sprintf("%s", dur.Round(time.Microsecond)),
	}

	if event.Err != nil {
		fields["err"] = event.Err.Error()
	}

	queryLines := strings.SplitN(event.Query, "\n", 2)
	query := queryLines[0]
	if len(queryLines) > 1 {
		query = query + "..."
	}

	logging.FromContext(ctx).WithFields(fields).Debug(query)
}
