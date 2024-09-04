package ledger

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	ledger "github.com/formancehq/ledger/internal"
	ledgercontroller "github.com/formancehq/ledger/internal/controller/ledger"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"github.com/formancehq/stack/libs/go-libs/query"
	"github.com/formancehq/stack/libs/go-libs/time"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

type Log struct {
	bun.BaseModel `bun:"logs,alias:logs"`

	Ledger         string     `bun:"ledger,type:varchar"`
	ID             int        `bun:"id,unique,type:numeric"`
	Type           string     `bun:"type,type:log_type"`
	Hash           []byte     `bun:"hash,type:bytea"`
	Date           time.Time  `bun:"date,type:timestamptz"`
	Data           RawMessage `bun:"data,type:jsonb"`
	IdempotencyKey *string    `bun:"idempotency_key,type:varchar(256),unique"`
}

func (log *Log) toCore() *ledger.ChainedLog {

	payload, err := ledger.HydrateLog(ledger.LogTypeFromString(log.Type), log.Data)
	if err != nil {
		panic(errors.Wrap(err, "hydrating log data"))
	}

	return &ledger.ChainedLog{
		Log: ledger.Log{
			Type: ledger.LogTypeFromString(log.Type),
			Data: payload,
			Date: log.Date.UTC(),
			IdempotencyKey: func() string {
				if log.IdempotencyKey != nil {
					return *log.IdempotencyKey
				}
				return ""
			}(),
		},
		ID:   log.ID,
		Hash: log.Hash,
	}
}

type RawMessage json.RawMessage

func (j RawMessage) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return string(j), nil
}

func (s *Store) logsQueryBuilder(q ledgercontroller.PaginatedQueryOptions[any]) func(*bun.SelectQuery) *bun.SelectQuery {
	return func(selectQuery *bun.SelectQuery) *bun.SelectQuery {

		selectQuery = selectQuery.Where("ledger = ?", s.ledgerName)
		if q.QueryBuilder != nil {
			subQuery, args, err := q.QueryBuilder.Build(query.ContextFn(func(key, operator string, value any) (string, []any, error) {
				switch {
				case key == "date":
					return fmt.Sprintf("%s %s ?", key, query.DefaultComparisonOperatorsMapping[operator]), []any{value}, nil
				default:
					return "", nil, fmt.Errorf("unknown key '%s' when building query", key)
				}
			}))
			if err != nil {
				panic(err)
			}
			selectQuery = selectQuery.Where(subQuery, args...)
		}

		return selectQuery
	}
}

func (s *Store) InsertLog(ctx context.Context, log ledger.Log) (*ledger.ChainedLog, error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}

	// we lock logs table as we need than the last log does not change until the transaction commit
	_, err = tx.ExecContext(ctx, "lock table "+s.PrefixWithBucket("logs"))
	if err != nil {
		return nil, errors.Wrap(err, "failed to lock logs")
	}

	lastLog := &Log{}
	var lastCoreLog *ledger.ChainedLog
	if err := tx.NewSelect().
		Model(lastLog).
		ModelTableExpr(s.PrefixWithBucketUsingModel(lastLog)).
		OrderExpr("id desc").
		Where("ledger = ?", s.ledgerName).
		Limit(1).
		Scan(ctx); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(err, "failed to read last log")
		}
	} else {
		lastCoreLog = lastLog.toCore()
	}

	newLog := log.ChainLog(lastCoreLog)

	data, err := json.Marshal(newLog.Data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal log data")
	}

	_, err = tx.
		NewInsert().
		Model(&Log{
			Ledger: s.ledgerName,
			ID:     newLog.ID,
			Type:   newLog.Type.String(),
			Hash:   newLog.Hash,
			Date:   newLog.Date,
			Data:   data,
			IdempotencyKey: func() *string {
				if newLog.IdempotencyKey == "" {
					return nil
				}
				return &newLog.IdempotencyKey
			}(),
		}).
		ModelTableExpr(s.PrefixWithBucketUsingModel(Log{})).
		Exec(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "inserting log")
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "failed to commit transaction")
	}

	return pointer.For(newLog), nil
}

func (s *Store) GetLastLog(ctx context.Context) (*ledger.ChainedLog, error) {
	ret, err := fetch[*Log](s, ctx,
		func(query *bun.SelectQuery) *bun.SelectQuery {
			return query.
				OrderExpr("id desc").
				Where("ledger = ?", s.ledgerName).
				Limit(1)
		})
	if err != nil {
		return nil, err
	}

	return ret.toCore(), nil
}

func (s *Store) GetLogs(ctx context.Context, q ledgercontroller.GetLogsQuery) (*bunpaginate.Cursor[ledger.ChainedLog], error) {
	logs, err := paginateWithColumn[ledgercontroller.PaginatedQueryOptions[any], Log](s, ctx,
		(*bunpaginate.ColumnPaginatedQuery[ledgercontroller.PaginatedQueryOptions[any]])(&q),
		s.logsQueryBuilder(q.Options),
	)
	if err != nil {
		return nil, err
	}

	return bunpaginate.MapCursor(logs, func(from Log) ledger.ChainedLog {
		return *from.toCore()
	}), nil
}

func (s *Store) ReadLogWithIdempotencyKey(ctx context.Context, key string) (*ledger.ChainedLog, error) {
	ret, err := fetch[*Log](s, ctx,
		func(query *bun.SelectQuery) *bun.SelectQuery {
			return query.
				OrderExpr("id desc").
				Limit(1).
				Where("idempotency_key = ?", key).
				Where("ledger = ?", s.ledgerName)
		})
	if err != nil {
		return nil, err
	}

	return ret.toCore(), nil
}
