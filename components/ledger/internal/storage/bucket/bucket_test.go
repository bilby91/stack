package bucket

import (
	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"
	"testing"

	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestBuckets(t *testing.T) {
	ctx := logging.TestingContext()
	name := uuid.NewString()[:8]

	<-srv.Done()

	pgDatabase := srv.GetValue().NewDatabase(t)
	db, err := bunconnect.OpenSQLDB(ctx, pgDatabase.ConnectionOptions())
	require.NoError(t, err)

	bucket := New(db, name)
	require.NoError(t, bucket.Migrate(ctx))
}
