//go:build it

package performance_test

import (
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/testing/docker"
	"github.com/formancehq/stack/libs/go-libs/testing/platform/pgtesting"
	"github.com/formancehq/stack/libs/go-libs/testing/utils"
	"testing"
)

var (
	dockerPool *docker.Pool
	pgServer   *pgtesting.PostgresServer
)

func TestMain(m *testing.M) {
	utils.WithTestMain(func(t *utils.TestingTForMain) int {
		dockerPool = docker.NewPool(t, logging.Testing())
		pgServer = pgtesting.CreatePostgresServer(t, dockerPool, pgtesting.WithPGStatsExtension())

		return m.Run()
	})
}
