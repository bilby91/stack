package ledger

import (
	"context"
	"fmt"
	"github.com/formancehq/ledger/internal/storage/bucket"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/migrations"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
	"reflect"
	"strings"
)

type Store struct {
	bucketName string
	ledgerName string
	db         bun.IDB
}

func (s *Store) Name() string {
	return s.ledgerName
}

func (s *Store) GetDB() bun.IDB {
	return s.db
}

func (s *Store) DiscoverBunTable(v any) string {
	if reflect.TypeOf(v).Kind() == reflect.Ptr {
		return s.DiscoverBunTable(reflect.ValueOf(v).Elem().Interface())
	}
	field, ok := reflect.TypeOf(v).FieldByName("BaseModel")
	if !ok {
		return ""
	}
	bunTag := field.Tag.Get("bun")
	tablePart := strings.SplitN(bunTag, ",", 2)[0]

	return strings.TrimPrefix(tablePart, "table:")
}

func (s *Store) PrefixWithBucketUsingModel(v any) string {
	table := s.DiscoverBunTable(v)
	if table == "" {
		return ""
	}
	return s.PrefixWithBucket(table)
}

func (s *Store) PrefixWithBucket(v string) string {
	return fmt.Sprintf(`"%s".%s`, s.bucketName, v)
}

func (s *Store) WithDB(db bun.IDB) *Store {
	return &Store{
		bucketName: s.bucketName,
		ledgerName: s.ledgerName,
		db:         db,
	}
}

func (s *Store) IsUpToDate(ctx context.Context) (bool, error) {
	bucketUpToDate, err := bucket.New(s.db, s.bucketName).IsUpToDate(ctx)
	if err != nil {
		return false, errors.Wrap(err, "failed to check if bucket is up to date")
	}
	if !bucketUpToDate {
		logging.FromContext(ctx).Errorf("bucket %s is not up to date", s.bucketName)
		return false, nil
	}

	ret, err := getMigrator(s.bucketName, s.ledgerName).IsUpToDate(ctx, s.db)
	if err != nil && errors.Is(err, migrations.ErrMissingVersionTable) {
		logging.FromContext(ctx).Errorf("ledger %s is not up to date", s.ledgerName)
		return false, nil
	}
	return ret, err
}

func New(
	db bun.IDB,
	bucketName, ledgerName string,
) *Store {
	return &Store{
		db:         db,
		bucketName: bucketName,
		ledgerName: ledgerName,
	}
}
