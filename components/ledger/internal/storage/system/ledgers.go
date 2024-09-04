package system

import (
	"context"
	ledger "github.com/formancehq/ledger/internal"
	system "github.com/formancehq/ledger/internal/controller/system"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/formancehq/stack/libs/go-libs/platform/postgres"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"github.com/formancehq/stack/libs/go-libs/time"

	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"

	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

const (
	StateInitializing = "initializing"
	StateInUse        = "in-use"
)

type Ledger struct {
	bun.BaseModel `bun:"_system.ledgers,alias:ledgers"`

	Name     string            `bun:"ledger,type:varchar(255),pk" json:"name"` // Primary key
	AddedAt  time.Time         `bun:"addedat,type:timestamp" json:"addedAt"`
	Bucket   string            `bun:"bucket,type:varchar(255)" json:"bucket"`
	Metadata map[string]string `bun:"metadata,type:jsonb" json:"metadata"`
	State    string            `bun:"state,type:varchar(255)" json:"-"`
}

func (l Ledger) toCore() ledger.Ledger {
	return ledger.Ledger{
		Name: l.Name,
		Configuration: ledger.Configuration{
			Bucket:   l.Bucket,
			Metadata: l.Metadata,
		},
		AddedAt: l.AddedAt,
		State:   l.State,
	}
}

func (s *Store) ListLedgers(ctx context.Context, q system.ListLedgersQuery) (*bunpaginate.Cursor[ledger.Ledger], error) {
	query := s.db.NewSelect().
		Model(&Ledger{}).
		Column("ledger", "bucket", "addedat", "metadata", "state").
		Order("addedat asc")

	cursor, err := bunpaginate.UsingOffset[system.PaginatedQueryOptions, Ledger](ctx, query, bunpaginate.OffsetPaginatedQuery[system.PaginatedQueryOptions](q))
	if err != nil {
		return nil, err
	}

	return bunpaginate.MapCursor(cursor, Ledger.toCore), nil
}

func (s *Store) DeleteLedger(ctx context.Context, name string) error {
	_, err := s.db.NewDelete().
		Model((*Ledger)(nil)).
		Where("ledger = ?", name).
		Exec(ctx)

	return errors.Wrap(postgres.ResolveError(err), "delete ledger from system store")
}

func (s *Store) RegisterLedger(ctx context.Context, l *ledger.Ledger) (bool, error) {
	return RegisterLedger(ctx, s.db, l)
}

func (s *Store) GetLedger(ctx context.Context, name string) (*ledger.Ledger, error) {
	ret := &Ledger{}
	if err := s.db.NewSelect().
		Model(ret).
		Column("ledger", "bucket", "addedat", "metadata", "state").
		Where("ledger = ?", name).
		Scan(ctx); err != nil {
		return nil, postgres.ResolveError(err)
	}

	return pointer.For(ret.toCore()), nil
}

func (s *Store) UpdateLedgerMetadata(ctx context.Context, name string, m metadata.Metadata) error {
	_, err := s.db.NewUpdate().
		Model(&Ledger{}).
		Set("metadata = metadata || ?", m).
		Where("ledger = ?", name).
		Exec(ctx)
	return err
}

func (s *Store) UpdateLedgerState(ctx context.Context, name string, state string) error {
	_, err := s.db.NewUpdate().
		Model(&Ledger{}).
		Set("state = ?", state).
		Where("ledger = ?", name).
		Exec(ctx)
	return err
}

func (s *Store) DeleteLedgerMetadata(ctx context.Context, name string, key string) error {
	_, err := s.db.NewUpdate().
		Model(&Ledger{}).
		Set("metadata = metadata - ?", key).
		Where("ledger = ?", name).
		Exec(ctx)
	return err
}

func RegisterLedger(ctx context.Context, db bun.IDB, l *ledger.Ledger) (bool, error) {
	if l.Metadata == nil {
		l.Metadata = metadata.Metadata{}
	}
	ret, err := db.NewInsert().
		Model(&Ledger{
			BaseModel: bun.BaseModel{},
			Name:      l.Name,
			AddedAt:   l.AddedAt,
			Bucket:    l.Bucket,
			Metadata:  l.Metadata,
			State:     l.State,
		}).
		Ignore().
		Exec(ctx)
	if err != nil {
		return false, postgres.ResolveError(err)
	}

	affected, err := ret.RowsAffected()
	if err != nil {
		return false, postgres.ResolveError(err)
	}

	return affected > 0, nil
}
