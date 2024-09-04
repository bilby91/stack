package utils

import (
	"context"
	"fmt"
	"github.com/shomali11/xsql"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
)

func DumpTables(t require.TestingT, ctx context.Context, db bun.IDB, queries... string) {
	for _, query := range queries {
		rows, err := db.QueryContext(ctx, query)
		require.NoError(t, err)

		prettied, err := xsql.Pretty(rows)
		require.NoError(t, err)

		fmt.Println(prettied)
	}
}
