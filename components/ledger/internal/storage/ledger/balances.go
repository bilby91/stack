package ledger

import (
	"context"
	"errors"
	"fmt"
	ledgercontroller "github.com/formancehq/ledger/internal/controller/ledger"
	"github.com/formancehq/stack/libs/go-libs/platform/postgres"
	"math/big"

	ledger "github.com/formancehq/ledger/internal"
	"github.com/formancehq/stack/libs/go-libs/query"
	"github.com/uptrace/bun"
)

type Balances struct {
	bun.BaseModel `bun:"balances"`

	Ledger  string   `bun:"ledger,type:varchar"`
	Account string   `bun:"account,type:varchar"`
	Asset   string   `bun:"asset,type:varchar"`
	Balance *big.Int `bun:"balance,type:numeric"`
}

func (s *Store) GetAggregatedBalances(ctx context.Context, q ledgercontroller.GetAggregatedBalanceQuery) (ledger.BalancesByAssets, error) {

	var (
		needMetadata bool
		subQuery     string
		args         []any
		err          error
	)
	if q.QueryBuilder != nil {
		subQuery, args, err = q.QueryBuilder.Build(query.ContextFn(func(key, operator string, value any) (string, []any, error) {
			switch {
			case key == "address":
				// TODO: Should allow comparison operator only if segments not used
				if operator != "$match" {
					return "", nil, newErrInvalidQuery("'address' column can only be used with $match")
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
				match := metadataRegex.FindAllStringSubmatch(key, 3)
				needMetadata = true
				key := "accounts.metadata"
				if q.PIT != nil {
					key = "am.metadata"
				}

				return key + " @> ?", []any{map[string]any{
					match[0][1]: value,
				}}, nil

			case key == "metadata":
				if operator != "$exists" {
					return "", nil, newErrInvalidQuery("'metadata' key filter can only be used with $exists")
				}
				needMetadata = true
				key := "accounts.metadata"
				if q.PIT != nil && !q.PIT.IsZero() {
					key = "am.metadata"
				}

				return fmt.Sprintf("%s -> ? IS NOT NULL", key), []any{value}, nil
			default:
				return "", nil, newErrInvalidQuery("unknown key '%s' when building query", key)
			}
		}))
		if err != nil {
			return nil, err
		}
	}

	type Temp struct {
		Aggregated ledger.VolumesByAssets `bun:"aggregated,type:jsonb"`
	}
	ret, err := fetch[*Temp](s, ctx,
		func(selectQuery *bun.SelectQuery) *bun.SelectQuery {
			pitColumn := "effective_date"
			if q.UseInsertionDate {
				pitColumn = "insertion_date"
			}
			moves := s.db.
				NewSelect().
				TableExpr(s.PrefixWithBucket("moves")).
				ColumnExpr("distinct on (moves.account_address, moves.asset) moves.*").
				Order("account_address", "asset").
				Where("moves.ledger = ?", s.ledgerName).
				Apply(filterPIT(q.PIT, pitColumn))

			if q.UseInsertionDate {
				moves = moves.Order("moves.insertion_date desc")
			} else {
				moves = moves.Order("moves.effective_date desc")
			}
			moves = moves.Order("seq desc")

			if needMetadata {
				if q.PIT != nil {
					moves = moves.Join(fmt.Sprintf(`join lateral (	
						select metadata
						from %s am 
						where am.accounts_seq = moves.accounts_seq and (? is null or date <= ?)
						order by revision desc 
						limit 1
					) am on true`, s.PrefixWithBucket("accounts_metadata")), q.PIT, q.PIT)
				} else {
					moves = moves.Join(fmt.Sprintf(`join lateral (	
						select metadata
						from %s a 
						where a.seq = moves.accounts_seq
					) accounts on true`, s.PrefixWithBucket("accounts")), q.PIT)
				}
			}
			if subQuery != "" {
				moves = moves.Where(subQuery, args...)
			}

			volumesColumn := "post_commit_effective_volumes"
			if q.UseInsertionDate {
				volumesColumn = "post_commit_volumes"
			}

			asJsonb := selectQuery.NewSelect().
				TableExpr("moves").
				ColumnExpr(fmt.Sprintf(s.PrefixWithBucket("volumes_to_jsonb((moves.asset, (sum((moves.%s).inputs), sum((moves.%s).outputs))::"+s.PrefixWithBucket("volumes")+"))")+" as aggregated", volumesColumn, volumesColumn)).
				Group("moves.asset")

			return selectQuery.
				With("moves", moves).
				With("data", asJsonb).
				TableExpr("data").
				ColumnExpr(s.PrefixWithBucket("aggregate_objects(data.aggregated)") + " as aggregated")
		})
	if err != nil && !errors.Is(err, postgres.ErrNotFound) {
		return nil, err
	}
	if errors.Is(err, postgres.ErrNotFound) {
		return ledger.BalancesByAssets{}, nil
	}

	return ret.Aggregated.Balances(), nil
}

func (s *Store) GetBalance(ctx context.Context, address, asset string) (*big.Int, error) {
	type Temp struct {
		Balance *big.Int `bun:"balance,type:numeric"`
	}
	v, err := fetch[*Temp](s, ctx, func(query *bun.SelectQuery) *bun.SelectQuery {
		return query.TableExpr(s.PrefixWithBucket("get_account_balance(?, ?, ?)")+" as balance", s.ledgerName, address, asset)
	})
	if err != nil {
		return nil, err
	}

	return v.Balance, nil
}

func (s *Store) AddToBalance(ctx context.Context, addr, asset string, amount *big.Int) (*big.Int, error) {
	r := &Balances{
		Ledger:  s.ledgerName,
		Account: addr,
		Asset:   asset,
		Balance: amount,
	}
	_, err := s.db.NewInsert().
		Model(r).
		ModelTableExpr(s.PrefixWithBucket("balances")).
		On("conflict (ledger, account, asset) do update").
		Set("balance = balances.balance + excluded.balance").
		Returning("balance").
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return r.Balance, err
}
