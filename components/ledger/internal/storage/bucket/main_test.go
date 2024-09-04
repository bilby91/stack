package bucket

import (
	"github.com/formancehq/stack/libs/go-libs/testing/docker"
	. "github.com/formancehq/stack/libs/go-libs/testing/utils"
	"testing"

	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/testing/platform/pgtesting"
)

var (
	srv = NewDeferred[*pgtesting.PostgresServer]()
)

func TestMain(m *testing.M) {
	WithTestMain(func(t *TestingTForMain) int {
		srv.LoadAsync(func() *pgtesting.PostgresServer {
			return pgtesting.CreatePostgresServer(t, docker.NewPool(t, logging.Testing()))
		})

		return m.Run()
	})
}
