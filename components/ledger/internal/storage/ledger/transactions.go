package ledger

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	ledgercontroller "github.com/formancehq/ledger/internal/controller/ledger"
	. "github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/formancehq/stack/libs/go-libs/platform/postgres"
	"math/big"
	"regexp"
	"strings"

	"github.com/formancehq/stack/libs/go-libs/pointer"

	"github.com/formancehq/stack/libs/go-libs/time"

	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"

	ledger "github.com/formancehq/ledger/internal"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/formancehq/stack/libs/go-libs/query"
	"github.com/uptrace/bun"
)

var (
	metadataRegex = regexp.MustCompile("metadata\\[(.+)\\]")
)

type Transaction struct {
	bun.BaseModel `bun:"table:transactions,alias:transactions"`

	Ledger            string            `bun:"ledger,type:varchar"`
	ID                int               `bun:"id,type:numeric"`
	Seq               int               `bun:"seq,scanonly"`
	Timestamp         *time.Time        `bun:"timestamp,type:timestamp without time zone"`
	Reference         string            `bun:"reference,type:varchar,unique,nullzero"`
	Postings          []ledger.Posting  `bun:"postings,type:jsonb"`
	Metadata          metadata.Metadata `bun:"metadata,type:jsonb,default:'{}'"`
	RevertedAt        *time.Time        `bun:"reverted_at"`
	InsertedAt        *time.Time        `bun:"inserted_at"`
	Sources           []string          `bun:"sources,type:jsonb"`
	Destinations      []string          `bun:"destinations,type:jsonb"`
	SourcesArray      []map[string]any  `bun:"sources_arrays,type:jsonb"`
	DestinationsArray []map[string]any  `bun:"destinations_arrays,type:jsonb"`
}

func (t *Transaction) toCore() ledger.Transaction {
	return ledger.Transaction{
		TransactionData: ledger.TransactionData{
			Reference:  t.Reference,
			Metadata:   t.Metadata,
			Timestamp:  *t.Timestamp,
			Postings:   t.Postings,
			InsertedAt: *t.InsertedAt,
		},
		ID:       t.ID,
		Reverted: t.RevertedAt != nil && !t.RevertedAt.IsZero(),
		Seq:      t.Seq,
	}
}

type ExpandedTransaction struct {
	Transaction
	bun.BaseModel `bun:"table:transactions,alias:transactions"`

	PostCommitEffectiveVolumes ledger.AccountsAssetsVolumes `bun:"post_commit_effective_volumes,type:jsonb"`
	PostCommitVolumes          ledger.AccountsAssetsVolumes `bun:"post_commit_volumes,type:jsonb"`
}

func (t ExpandedTransaction) toCore() ledger.ExpandedTransaction {
	var (
		preCommitEffectiveVolumes ledger.AccountsAssetsVolumes
		preCommitVolumes          ledger.AccountsAssetsVolumes
	)
	if t.PostCommitEffectiveVolumes != nil {
		preCommitEffectiveVolumes = t.PostCommitEffectiveVolumes.Copy()
		for _, posting := range t.Postings {
			preCommitEffectiveVolumes.AddOutput(posting.Source, posting.Asset, big.NewInt(0).Neg(posting.Amount))
			preCommitEffectiveVolumes.AddInput(posting.Destination, posting.Asset, big.NewInt(0).Neg(posting.Amount))
		}
	}
	if t.PostCommitVolumes != nil {
		preCommitVolumes = t.PostCommitVolumes.Copy()
		for _, posting := range t.Postings {
			preCommitVolumes.AddOutput(posting.Source, posting.Asset, big.NewInt(0).Neg(posting.Amount))
			preCommitVolumes.AddInput(posting.Destination, posting.Asset, big.NewInt(0).Neg(posting.Amount))
		}
	}
	return ledger.ExpandedTransaction{
		Transaction:                t.Transaction.toCore(),
		PreCommitEffectiveVolumes:  preCommitEffectiveVolumes,
		PostCommitEffectiveVolumes: t.PostCommitEffectiveVolumes,
		PreCommitVolumes:           preCommitVolumes,
		PostCommitVolumes:          t.PostCommitVolumes,
	}
}

type account string

var _ driver.Valuer = account("")

func (m1 account) Value() (driver.Value, error) {
	ret, err := json.Marshal(strings.Split(string(m1), ":"))
	if err != nil {
		return nil, err
	}
	return string(ret), nil
}

// Scan - Implement the database/sql scanner interface
func (m1 *account) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	v, err := driver.String.ConvertValue(value)
	if err != nil {
		return err
	}

	array := make([]string, 0)
	switch vv := v.(type) {
	case []uint8:
		err = json.Unmarshal(vv, &array)
	case string:
		err = json.Unmarshal([]byte(vv), &array)
	default:
		panic("not handled type")
	}
	if err != nil {
		return err
	}
	*m1 = account(strings.Join(array, ":"))
	return nil
}

func (s *Store) buildTransactionQuery(p ledgercontroller.PITFilterWithVolumes, query *bun.SelectQuery) *bun.SelectQuery {

	selectMetadata := query.NewSelect().
		TableExpr(s.PrefixWithBucket("transactions_metadata")).
		Where("transactions.seq = transactions_metadata.transactions_seq").
		Order("revision desc").
		Limit(1)

	if p.PIT != nil && !p.PIT.IsZero() {
		selectMetadata = selectMetadata.Where("date <= ?", p.PIT)
	}

	query = query.
		Where("transactions.ledger = ?", s.ledgerName)

	if p.PIT != nil && !p.PIT.IsZero() {
		query = query.
			Where("timestamp <= ?", p.PIT).
			ColumnExpr("transactions.*").
			Column("transactions_metadata.metadata").
			Join(fmt.Sprintf(`left join lateral (%s) as transactions_metadata on true`, selectMetadata.String())).
			ColumnExpr(fmt.Sprintf("case when reverted_at is not null and reverted_at > '%s' then null else reverted_at end", p.PIT.Format(time.DateFormat)))
	} else {
		query = query.Column("transactions.metadata", "transactions.*")
	}

	if p.ExpandEffectiveVolumes {
		query = query.ColumnExpr(s.PrefixWithBucket("get_aggregated_effective_volumes_for_transaction(?, transactions.seq) as post_commit_effective_volumes"), s.ledgerName)
	}
	if p.ExpandVolumes {
		query = query.ColumnExpr(s.PrefixWithBucket("get_aggregated_volumes_for_transaction(?, transactions.seq) as post_commit_volumes"), s.ledgerName)
	}
	return query
}

func (s *Store) transactionQueryContext(qb query.Builder, q ledgercontroller.GetTransactionsQuery) (string, []any, error) {

	return qb.Build(query.ContextFn(func(key, operator string, value any) (string, []any, error) {
		switch {
		case key == "reference" || key == "timestamp":
			return fmt.Sprintf("%s %s ?", key, query.DefaultComparisonOperatorsMapping[operator]), []any{value}, nil
		case key == "reverted":
			if operator != "$match" {
				return "", nil, newErrInvalidQuery("'reverted' column can only be used with $match")
			}
			switch value := value.(type) {
			case bool:
				ret := "reverted_at is"
				if value {
					ret += " not"
				}
				return ret + " null", nil, nil
			default:
				return "", nil, newErrInvalidQuery("'reverted' can only be used with bool value")
			}
		case key == "account":
			// TODO: Should allow comparison operator only if segments not used
			if operator != "$match" {
				return "", nil, newErrInvalidQuery("'account' column can only be used with $match")
			}
			switch address := value.(type) {
			case string:
				return filterAccountAddressOnTransactions(address, true, true), nil, nil
			default:
				return "", nil, newErrInvalidQuery("unexpected type %T for column 'account'", address)
			}
		case key == "source":
			// TODO: Should allow comparison operator only if segments not used
			if operator != "$match" {
				return "", nil, errors.New("'source' column can only be used with $match")
			}
			switch address := value.(type) {
			case string:
				return filterAccountAddressOnTransactions(address, true, false), nil, nil
			default:
				return "", nil, newErrInvalidQuery("unexpected type %T for column 'source'", address)
			}
		case key == "destination":
			// TODO: Should allow comparison operator only if segments not used
			if operator != "$match" {
				return "", nil, errors.New("'destination' column can only be used with $match")
			}
			switch address := value.(type) {
			case string:
				return filterAccountAddressOnTransactions(address, false, true), nil, nil
			default:
				return "", nil, newErrInvalidQuery("unexpected type %T for column 'destination'", address)
			}
		case metadataRegex.Match([]byte(key)):
			if operator != "$match" {
				return "", nil, newErrInvalidQuery("'account' column can only be used with $match")
			}
			match := metadataRegex.FindAllStringSubmatch(key, 3)

			key := "metadata"
			if q.Options.Options.PIT != nil && !q.Options.Options.PIT.IsZero() {
				key = "transactions_metadata.metadata"
			}

			return key + " @> ?", []any{map[string]any{
				match[0][1]: value,
			}}, nil

		case key == "metadata":
			if operator != "$exists" {
				return "", nil, newErrInvalidQuery("'metadata' key filter can only be used with $exists")
			}
			if q.Options.Options.PIT != nil && !q.Options.Options.PIT.IsZero() {
				key = "transactions_metadata.metadata"
			}

			return fmt.Sprintf("%s -> ? IS NOT NULL", key), []any{value}, nil
		default:
			return "", nil, newErrInvalidQuery("unknown key '%s' when building query", key)
		}
	}))
}

func (s *Store) buildTransactionListQuery(selectQuery *bun.SelectQuery, q ledgercontroller.PaginatedQueryOptions[ledgercontroller.PITFilterWithVolumes], where string, args []any) *bun.SelectQuery {

	selectQuery = s.buildTransactionQuery(q.Options, selectQuery)
	if where != "" {
		return selectQuery.Where(where, args...)
	}

	return selectQuery
}

func (s *Store) GetTransactions(ctx context.Context, q ledgercontroller.GetTransactionsQuery) (*bunpaginate.Cursor[ledger.ExpandedTransaction], error) {

	var (
		where string
		args  []any
		err   error
	)
	if q.Options.QueryBuilder != nil {
		where, args, err = s.transactionQueryContext(q.Options.QueryBuilder, q)
		if err != nil {
			return nil, err
		}
	}

	transactions, err := paginateWithColumn[ledgercontroller.PaginatedQueryOptions[ledgercontroller.PITFilterWithVolumes], ExpandedTransaction](s, ctx,
		(*bunpaginate.ColumnPaginatedQuery[ledgercontroller.PaginatedQueryOptions[ledgercontroller.PITFilterWithVolumes]])(&q),
		func(query *bun.SelectQuery) *bun.SelectQuery {
			return s.buildTransactionListQuery(query, q.Options, where, args)
		},
	)
	if err != nil {
		return nil, err
	}

	return bunpaginate.MapCursor(transactions, ExpandedTransaction.toCore), nil
}

func (s *Store) CountTransactions(ctx context.Context, q ledgercontroller.GetTransactionsQuery) (int, error) {

	var (
		where string
		args  []any
		err   error
	)

	if q.Options.QueryBuilder != nil {
		where, args, err = s.transactionQueryContext(q.Options.QueryBuilder, q)
		if err != nil {
			return 0, err
		}
	}

	return count[ExpandedTransaction](s, ctx, func(query *bun.SelectQuery) *bun.SelectQuery {
		return s.buildTransactionListQuery(query, q.Options, where, args)
	})
}

func (s *Store) GetTransactionWithVolumes(ctx context.Context, filter ledgercontroller.GetTransactionQuery) (*ledger.ExpandedTransaction, error) {
	ret, err := fetch[*ExpandedTransaction](s, ctx,
		func(query *bun.SelectQuery) *bun.SelectQuery {
			return s.buildTransactionQuery(filter.PITFilterWithVolumes, query).
				Where("transactions.id = ?", filter.ID).
				Limit(1)
		})
	if err != nil {
		return nil, err
	}

	return pointer.For(ret.toCore()), nil
}

func (s *Store) GetTransaction(ctx context.Context, txId int) (*ledger.Transaction, error) {
	tx, err := fetch[*Transaction](s, ctx,
		func(query *bun.SelectQuery) *bun.SelectQuery {
			return query.
				ColumnExpr(`transactions.id, transactions.inserted_at, transactions.reference, transactions.postings, transactions.timestamp, transactions.reverted_at, tm.metadata`).
				Join("left join"+s.PrefixWithBucket("transactions_metadata")+" tm on tm.transactions_seq = transactions.seq").
				Where("transactions.id = ?", txId).
				Where("transactions.ledger = ?", s.ledgerName).
				Order("tm.revision desc").
				Limit(1)
		})
	if err != nil {
		return nil, err
	}

	return pointer.For(tx.toCore()), nil
}

func (s *Store) GetTransactionByReference(ctx context.Context, ref string) (*ledger.ExpandedTransaction, error) {
	ret, err := fetch[*ExpandedTransaction](s, ctx,
		func(query *bun.SelectQuery) *bun.SelectQuery {
			return query.
				ColumnExpr(`transactions.*, tm.metadata`).
				Join("left join "+s.PrefixWithBucket("transactions_metadata")+" tm on tm.transactions_seq = transactions.seq").
				Where("transactions.reference = ?", ref).
				Where("transactions.ledger = ?", s.ledgerName).
				Order("tm.revision desc").
				Limit(1)
		})
	if err != nil {
		return nil, err
	}

	return pointer.For(ret.toCore()), nil
}

func (s *Store) GetLastTransaction(ctx context.Context) (*ledger.ExpandedTransaction, error) {
	ret, err := fetch[*ExpandedTransaction](s, ctx,
		func(query *bun.SelectQuery) *bun.SelectQuery {
			return query.
				ColumnExpr(`transactions.*, tm.metadata`).
				Join("left join "+s.PrefixWithBucket("transactions_metadata")+" tm on tm.transactions_seq = transactions.seq").
				Order("transactions.seq desc", "tm.revision desc").
				Where("transactions.ledger = ?", s.ledgerName).
				Limit(1)
		})
	if err != nil {
		return nil, err
	}

	return pointer.For(ret.toCore()), nil
}

func (s *Store) InsertTransaction(ctx context.Context, data ledger.TransactionData) (*ledger.Transaction, error) {
	sources := Map(data.Postings, ledger.Posting.GetSource)
	destinations := Map(data.Postings, ledger.Posting.GetDestination)
	tx := &Transaction{
		Ledger:   s.ledgerName,
		Postings: data.Postings,
		Metadata: data.Metadata,
		Timestamp: func() *time.Time {
			if data.Timestamp.IsZero() {
				return nil
			}
			return &data.Timestamp
		}(),
		Reference: data.Reference,
		InsertedAt: func() *time.Time {
			if data.InsertedAt.IsZero() {
				return nil
			}
			return &data.InsertedAt
		}(),
		Sources:           sources,
		Destinations:      destinations,
		SourcesArray:      Map(sources, convertAddrToIndexedJSONB),
		DestinationsArray: Map(destinations, convertAddrToIndexedJSONB),
	}
	_, err := s.db.NewInsert().
		Model(tx).
		ModelTableExpr(s.PrefixWithBucket("transactions")).
		Value("id", "nextval(?)", s.PrefixWithBucket(fmt.Sprintf(`"%s_transaction_id"`, s.ledgerName))).
		Returning("*").
		Exec(ctx)
	if err != nil {
		return nil, postgres.ResolveError(err)
	}

	return pointer.For(tx.toCore()), nil
}

func (s *Store) RevertTransaction(ctx context.Context, id int) (*ledger.Transaction, bool, error) {
	ret := &Transaction{}
	now := time.Now()
	sqlResult, err := s.db.NewUpdate().
		Model(ret).
		ModelTableExpr(s.PrefixWithBucket("transactions")).
		Where("id = ?", id).
		Where("reverted_at is null").
		Where("ledger = ?", s.ledgerName).
		Set("reverted_at = ?", now).
		Set("updated_at = ?", now).
		Returning("*").
		Exec(ctx)
	if err != nil {
		return nil, false, postgres.ResolveError(err)
	}

	rowsAffected, err := sqlResult.RowsAffected()
	if err != nil {
		return nil, false, postgres.ResolveError(err)
	}

	if rowsAffected == 0 {
		return pointer.For(ret.toCore()), false, nil
	}

	return pointer.For(ret.toCore()), true, nil
}

func (s *Store) UpdateTransactionMetadata(ctx context.Context, transactionID int, m metadata.Metadata) (*ledger.Transaction, error) {
	tx := &Transaction{}
	_, err := s.db.NewUpdate().
		Model(tx).
		ModelTableExpr(s.PrefixWithBucket("transactions")).
		Where("id = ?", transactionID).
		Where("ledger = ?", s.ledgerName).
		Set("metadata = metadata || ?", m).
		Set("updated_at = ?", time.Now()).
		Returning("*").
		Exec(ctx)
	if err != nil {
		return nil, err
	}
	return pointer.For(tx.toCore()), nil
}

func (s *Store) DeleteTransactionMetadata(ctx context.Context, id int, key string) (*ledger.Transaction, error) {
	ret := &Transaction{}
	_, err := s.db.NewUpdate().
		Model(ret).
		ModelTableExpr(s.PrefixWithBucketUsingModel(Transaction{})).
		Set("metadata = metadata - ?", key).
		Where("id = ?", id).
		Where("ledger = ?", s.ledgerName).
		Returning("*").
		Exec(ctx)
	if err != nil {
		return nil, postgres.ResolveError(err)
	}

	return pointer.For(ret.toCore()), nil
}

func convertAddrToIndexedJSONB(addr string) map[string]any {
	ret := map[string]any{}
	parts := strings.Split(addr, ":")
	for i := range parts {
		ret[fmt.Sprint(i)] = parts[i]
	}
	ret[fmt.Sprint(len(parts))] = nil

	return ret
}
