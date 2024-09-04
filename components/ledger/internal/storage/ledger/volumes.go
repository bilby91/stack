package ledger

import (
	"context"
	"fmt"
	ledgercontroller "github.com/formancehq/ledger/internal/controller/ledger"
	"regexp"

	ledger "github.com/formancehq/ledger/internal"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	lquery "github.com/formancehq/stack/libs/go-libs/query"
	"github.com/uptrace/bun"
)

func (s *Store) volumesQueryContext(q ledgercontroller.GetVolumesWithBalancesQuery) (string, []any, bool, error) {

	metadataRegex := regexp.MustCompile("metadata\\[(.+)\\]")
	balanceRegex := regexp.MustCompile("balance\\[(.*)\\]")
	var (
		subQuery string
		args     []any
		err      error
	)

	var useMetadata = false

	if q.Options.QueryBuilder != nil {
		subQuery, args, err = q.Options.QueryBuilder.Build(lquery.ContextFn(func(key, operator string, value any) (string, []any, error) {

			convertOperatorToSQL := func() string {
				switch operator {
				case "$match":
					return "="
				case "$lt":
					return "<"
				case "$gt":
					return ">"
				case "$lte":
					return "<="
				case "$gte":
					return ">="
				}
				panic("unreachable")
			}

			switch {
			case key == "account" || key == "address":
				// TODO: Should allow comparison operator only if segments not used
				if operator != "$match" {
					return "", nil, newErrInvalidQuery(fmt.Sprintf("'%s' column can only be used with $match", key))
				}

				switch address := value.(type) {
				case string:
					return filterAccountAddress(address, "account_address"), nil, nil
				default:
					return "", nil, newErrInvalidQuery("unexpected type %T for column 'address'", address)
				}
			case metadataRegex.Match([]byte(key)):
				if operator != "$match" {
					return "", nil, newErrInvalidQuery("'metadata' column can only be used with $match")
				}
				useMetadata = true
				match := metadataRegex.FindAllStringSubmatch(key, 3)
				key := "metadata"

				return key + " @> ?", []any{map[string]any{
					match[0][1]: value,
				}}, nil
			case key == "metadata":
				if operator != "$exists" {
					return "", nil, newErrInvalidQuery("'metadata' key filter can only be used with $exists")
				}
				useMetadata = true
				key := "metadata"

				return fmt.Sprintf("%s -> ? IS NOT NULL", key), []any{value}, nil
			case balanceRegex.Match([]byte(key)):
				match := balanceRegex.FindAllStringSubmatch(key, 2)
				return fmt.Sprintf(`balance %s ?  and asset = ?`, convertOperatorToSQL()), []any{value, match[0][1]}, nil
			default:
				return "", nil, newErrInvalidQuery("unknown key '%s' when building query", key)
			}
		}))
		if err != nil {
			return "", nil, false, err
		}
	}

	return subQuery, args, useMetadata, nil

}

func (s *Store) buildVolumesWithBalancesQuery(query *bun.SelectQuery, q ledgercontroller.GetVolumesWithBalancesQuery, where string, args []any, useMetadata bool) *bun.SelectQuery {

	filtersForVolumes := q.Options.Options
	dateFilterColumn := "effective_date"

	if filtersForVolumes.UseInsertionDate {
		dateFilterColumn = "insertion_date"
	}

	query = query.
		Column("account_address_array").
		Column("account_address").
		Column("asset").
		ColumnExpr("sum(case when not is_source then amount else 0 end) as input").
		ColumnExpr("sum(case when is_source then amount else 0 end) as output").
		ColumnExpr("sum(case when not is_source then amount else -amount end) as balance").
		TableExpr(s.PrefixWithBucket("moves"))

	if useMetadata {
		query = query.ColumnExpr("accounts.metadata as metadata").
			Join(fmt.Sprintf(`join lateral (	
		select metadata
		from %s a 
		where a.seq = moves.accounts_seq
		) accounts on true`, s.PrefixWithBucket("accounts"))).Group("metadata")
	}

	query = query.
		Where("ledger = ?", s.ledgerName).
		Apply(filterPIT(filtersForVolumes.PIT, dateFilterColumn)).
		Apply(filterOOT(filtersForVolumes.OOT, dateFilterColumn)).
		GroupExpr("account_address, account_address_array, asset")

	globalQuery := query.NewSelect()
	globalQuery = globalQuery.
		With("query", query).
		TableExpr("query")

	if filtersForVolumes.GroupLvl > 0 {
		globalQuery = globalQuery.
			ColumnExpr(fmt.Sprintf(`(array_to_string((string_to_array(account_address, ':'))[1:LEAST(array_length(string_to_array(account_address, ':'),1),%d)],':')) as account`, filtersForVolumes.GroupLvl)).
			Column("asset").
			ColumnExpr("sum(input) as input").
			ColumnExpr("sum(output) as output").
			ColumnExpr("sum(balance) as balance").
			GroupExpr("account, asset")
	} else {
		globalQuery = globalQuery.ColumnExpr("account_address as account, asset, input, output, balance")
	}

	if useMetadata {
		globalQuery = globalQuery.Column("metadata")
	}

	if where != "" {
		globalQuery.Where(where, args...)
	}

	return globalQuery
}

func (s *Store) GetVolumesWithBalances(ctx context.Context, q ledgercontroller.GetVolumesWithBalancesQuery) (*bunpaginate.Cursor[ledger.VolumesWithBalanceByAssetByAccount], error) {
	var (
		where       string
		args        []any
		err         error
		useMetadata bool
	)
	if q.Options.QueryBuilder != nil {
		where, args, useMetadata, err = s.volumesQueryContext(q)
		if err != nil {
			return nil, err
		}
	}

	return paginateWithOffsetWithoutModel[ledgercontroller.PaginatedQueryOptions[ledgercontroller.FiltersForVolumes], ledger.VolumesWithBalanceByAssetByAccount](
		s, ctx, (*bunpaginate.OffsetPaginatedQuery[ledgercontroller.PaginatedQueryOptions[ledgercontroller.FiltersForVolumes]])(&q),
		func(query *bun.SelectQuery) *bun.SelectQuery {
			return s.buildVolumesWithBalancesQuery(query, q, where, args, useMetadata)
		},
	)
}
