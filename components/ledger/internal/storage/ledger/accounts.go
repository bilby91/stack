package ledger

import (
	"context"
	"database/sql"
	"fmt"
	ledgercontroller "github.com/formancehq/ledger/internal/controller/ledger"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/formancehq/stack/libs/go-libs/platform/postgres"
	"github.com/pkg/errors"
	"regexp"
	"strings"

	"github.com/formancehq/stack/libs/go-libs/time"

	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"

	ledger "github.com/formancehq/ledger/internal"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"github.com/formancehq/stack/libs/go-libs/query"
	"github.com/uptrace/bun"
)

type Account struct {
	bun.BaseModel `bun:"table:accounts"`

	Ledger        string            `bun:"ledger"`
	Address       string            `bun:"address"`
	AddressArray  []string          `bun:"address_array"`
	Metadata      map[string]string `bun:"metadata,type:jsonb"`
	InsertionDate time.Time         `bun:"insertion_date"`
	UpdatedAt     time.Time         `bun:"updated_at"`
	FirstUsage    time.Time         `bun:"first_usage"`
}

func (s *Store) buildAccountQuery(q ledgercontroller.PITFilterWithVolumes, query *bun.SelectQuery) *bun.SelectQuery {

	query = query.
		Column("accounts.address", "accounts.first_usage").
		Where("accounts.ledger = ?", s.ledgerName).
		Apply(filterPIT(q.PIT, "first_usage")).
		Order("accounts.address")

	if q.PIT != nil && !q.PIT.IsZero() {
		query = query.
			Column("accounts.address").
			ColumnExpr(`coalesce(accounts_metadata.metadata, '{}'::jsonb) as metadata`).
			Join(fmt.Sprintf(`
				left join lateral (
					select metadata, accounts_seq
					from %s
					where accounts_metadata.accounts_seq = accounts.seq and accounts_metadata.date < ?
					order by revision desc 
					limit 1
				) accounts_metadata on true
			`, s.PrefixWithBucket("accounts_metadata")), q.PIT)
	} else {
		query = query.Column("metadata")
	}

	if q.ExpandVolumes {
		query = query.
			ColumnExpr("volumes.*").
			Join("join "+s.PrefixWithBucket("get_account_aggregated_volumes(?, accounts.address, ?)")+"volumes on true", s.ledgerName, q.PIT)
	}

	if q.ExpandEffectiveVolumes {
		query = query.
			ColumnExpr("effective_volumes.*").
			Join("join "+s.PrefixWithBucket("get_account_aggregated_effective_volumes(?, accounts.address, ?)")+" effective_volumes on true", s.ledgerName, q.PIT)
	}

	return query
}

func (s *Store) accountQueryContext(qb query.Builder, q ledgercontroller.GetAccountsQuery) (string, []any, error) {
	metadataRegex := regexp.MustCompile("metadata\\[(.+)\\]")
	balanceRegex := regexp.MustCompile("balance\\[(.*)\\]")

	return qb.Build(query.ContextFn(func(key, operator string, value any) (string, []any, error) {
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
		case key == "address":
			// TODO: Should allow comparison operator only if segments not used
			if operator != "$match" {
				return "", nil, errors.New("'address' column can only be used with $match")
			}
			switch address := value.(type) {
			case string:
				return filterAccountAddress(address, "accounts.address"), nil, nil
			default:
				return "", nil, newErrInvalidQuery("unexpected type %T for column 'address'", address)
			}
		case metadataRegex.Match([]byte(key)):
			if operator != "$match" {
				return "", nil, newErrInvalidQuery("'account' column can only be used with $match")
			}
			match := metadataRegex.FindAllStringSubmatch(key, 3)

			key := "metadata"
			if q.Options.Options.PIT != nil && !q.Options.Options.PIT.IsZero() {
				key = "accounts_metadata.metadata"
			}

			return key + " @> ?", []any{map[string]any{
				match[0][1]: value,
			}}, nil
		case balanceRegex.Match([]byte(key)):
			match := balanceRegex.FindAllStringSubmatch(key, 2)

			return fmt.Sprintf(`(
				select %s
				from %s
				where asset = ? and account_address = accounts.address and ledger = ?
				order by seq desc
				limit 1
			) %s ?`, s.PrefixWithBucket("balance_from_volumes(post_commit_volumes)"), s.PrefixWithBucket("moves"), convertOperatorToSQL()), []any{match[0][1], s.ledgerName, value}, nil
		case key == "balance":
			return fmt.Sprintf(`(
				select %s
				from %s
				where account_address = accounts.address and ledger = ?
				order by seq desc
				limit 1
			) %s ?`, s.PrefixWithBucket("balance_from_volumes(post_commit_volumes)"), s.PrefixWithBucket("moves"), convertOperatorToSQL()), []any{s.ledgerName, value}, nil

		case key == "metadata":
			if operator != "$exists" {
				return "", nil, newErrInvalidQuery("'metadata' key filter can only be used with $exists")
			}
			if q.Options.Options.PIT != nil && !q.Options.Options.PIT.IsZero() {
				key = "accounts_metadata.metadata"
			}

			return fmt.Sprintf("%s -> ? IS NOT NULL", key), []any{value}, nil
		default:
			return "", nil, newErrInvalidQuery("unknown key '%s' when building query", key)
		}
	}))
}

func (s *Store) buildAccountListQuery(selectQuery *bun.SelectQuery, q ledgercontroller.GetAccountsQuery, where string, args []any) *bun.SelectQuery {
	selectQuery = s.buildAccountQuery(q.Options.Options, selectQuery)

	if where != "" {
		return selectQuery.Where(where, args...)
	}

	return selectQuery
}

func (s *Store) GetAccountsWithVolumes(ctx context.Context, q ledgercontroller.GetAccountsQuery) (*bunpaginate.Cursor[ledger.ExpandedAccount], error) {
	var (
		where string
		args  []any
		err   error
	)
	if q.Options.QueryBuilder != nil {
		where, args, err = s.accountQueryContext(q.Options.QueryBuilder, q)
		if err != nil {
			return nil, err
		}
	}

	return paginateWithOffset[ledgercontroller.PaginatedQueryOptions[ledgercontroller.PITFilterWithVolumes], ledger.ExpandedAccount](s, ctx,
		(*bunpaginate.OffsetPaginatedQuery[ledgercontroller.PaginatedQueryOptions[ledgercontroller.PITFilterWithVolumes]])(&q),
		func(query *bun.SelectQuery) *bun.SelectQuery {
			return s.buildAccountListQuery(query, q, where, args)
		},
	)
}

func (s *Store) GetAccount(ctx context.Context, address string) (*ledger.Account, error) {
	account, err := fetch[*ledger.Account](s, ctx, func(query *bun.SelectQuery) *bun.SelectQuery {
		return query.
			ColumnExpr("accounts.*").
			ColumnExpr("coalesce(accounts_metadata.metadata, '{}'::jsonb) as metadata").
			Join(fmt.Sprintf("left join %s on accounts_metadata.accounts_seq = accounts.seq", s.PrefixWithBucket("accounts_metadata"))).
			Where("accounts.address = ?", address).
			Where("accounts.ledger = ?", s.ledgerName).
			Order("revision desc").
			Limit(1)
	})
	if err != nil {
		if postgres.IsNotFoundError(err) {
			return pointer.For(ledger.NewAccount(address)), nil
		}
		return nil, err
	}

	return account, nil
}

func (s *Store) GetAccountWithVolumes(ctx context.Context, q ledgercontroller.GetAccountQuery) (*ledger.ExpandedAccount, error) {
	account, err := fetch[*ledger.ExpandedAccount](s, ctx, func(query *bun.SelectQuery) *bun.SelectQuery {
		query = s.buildAccountQuery(q.PITFilterWithVolumes, query).
			Where("accounts.address = ?", q.Addr).
			Limit(1)

		return query
	})
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (s *Store) CountAccounts(ctx context.Context, q ledgercontroller.GetAccountsQuery) (int, error) {
	var (
		where string
		args  []any
		err   error
	)
	if q.Options.QueryBuilder != nil {
		where, args, err = s.accountQueryContext(q.Options.QueryBuilder, q)
		if err != nil {
			return 0, err
		}
	}

	return count[ledger.Account](s, ctx, func(query *bun.SelectQuery) *bun.SelectQuery {
		return s.buildAccountListQuery(query, q, where, args)
	})
}

func (s *Store) UpdateAccountMetadata(ctx context.Context, account string, m metadata.Metadata) error {
	_, err := s.db.NewInsert().
		Model(&Account{
			Ledger:        s.ledgerName,
			Address:       account,
			AddressArray:  strings.Split(account, ":"),
			Metadata:      m,
			InsertionDate: time.Now(),
			UpdatedAt:     time.Now(),
		}).
		ModelTableExpr(s.PrefixWithBucketUsingModel(Account{})).
		On("CONFLICT (ledger, address) DO UPDATE").
		Set("metadata = excluded.metadata || ?", m).
		Exec(ctx)
	return postgres.ResolveError(err)
}

func (s *Store) DeleteAccountMetadata(ctx context.Context, account, key string) error {
	_, err := s.db.NewUpdate().
		ModelTableExpr(s.PrefixWithBucketUsingModel(Account{})).
		Set("metadata = metadata - ?", key).
		Where("address = ?", account).
		Where("ledger = ?", s.ledgerName).
		Exec(ctx)
	return postgres.ResolveError(err)
}

func (s *Store) UpsertAccount(ctx context.Context, account ledger.Account) (bool, error) {

	model := &Account{
		BaseModel:     bun.BaseModel{},
		Ledger:        s.ledgerName,
		Address:       account.Address,
		AddressArray:  strings.Split(account.Address, ":"),
		InsertionDate: account.InsertionDate,
		UpdatedAt:     account.UpdatedAt,
		Metadata:      account.Metadata,
		FirstUsage:    account.FirstUsage,
	}

	//result, err := s.db.NewInsert().
	//	Model(model).
	//	ModelTableExpr(s.PrefixWithBucketUsingModel(Account{})).
	//	On("conflict (ledger, address) do update").
	//	Set("first_usage = case when ? < excluded.first_usage then ? else excluded.first_usage end", account.FirstUsage, account.FirstUsage).
	//	Set("updated_at = ?", account.UpdatedAt).
	//	Set("metadata = excluded.metadata || ?", account.Metadata).
	//	Where("? < accounts.first_usage or not accounts.metadata @> coalesce(?, '{}'::jsonb)", account.FirstUsage, account.Metadata).
	//	Returning("ctid, xmin, xmax").
	//	Exec(ctx)
	//if err != nil {
	//	return false, err
	//}
	//rowsAffected, err := result.RowsAffected()
	//if err != nil {
	//	return false, err
	//}
	//if rowsAffected == 0 {
	//	return false, nil
	//}
	//
	//return true, nil

	var (
		rowsAffected int64
		err          error
	)
	err = s.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		var result sql.Result
		result, err = s.db.NewInsert().
			Model(model).
			ModelTableExpr(s.PrefixWithBucketUsingModel(Account{})).
			On("conflict (ledger, address) do update").
			Set("first_usage = case when ? < excluded.first_usage then ? else excluded.first_usage end", account.FirstUsage, account.FirstUsage).
			Set("updated_at = ?", account.UpdatedAt).
			Set("metadata = excluded.metadata || ?", account.Metadata).
			Where("? < accounts.first_usage or not accounts.metadata @> coalesce(?, '{}'::jsonb)", account.FirstUsage, account.Metadata).
			Exec(ctx)
		if err != nil {
			return err
		}
		rowsAffected, err = result.RowsAffected()
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			// by roll-backing the transaction, we release the lock, allowing a concurrent transaction
			// to use the table
			return tx.Rollback()
		}

		return nil
	})
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}

func (s *Store) LockAccounts(ctx context.Context, accounts ...string) error {
	rows, err := s.db.QueryContext(ctx, strings.Join(collectionutils.Map(accounts, func(account string) string {
		// todo: add bucket name in the advisory lock number computation
		return fmt.Sprintf(`select pg_advisory_xact_lock(hashtext('%s'))`, account)
	}), ";"))
	if err != nil {
		return errors.Wrap(err, "failed to lock accounts")
	}

	return rows.Close()
}
