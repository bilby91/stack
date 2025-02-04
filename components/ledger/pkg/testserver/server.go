package testserver

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/formancehq/ledger/cmd"
	ledgerclient "github.com/formancehq/stack/ledger/client"
	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"
	"github.com/formancehq/stack/libs/go-libs/httpclient"
	"github.com/formancehq/stack/libs/go-libs/httpserver"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/stretchr/testify/require"
)

type T interface {
	require.TestingT
	TempDir() string
	Cleanup(func())
	Helper()
	Logf(format string, args ...any)
}

type Configuration struct {
	PostgresConfiguration bunconnect.ConnectionOptions
	Output                io.Writer
	Debug                 bool
}

type Server struct {
	configuration Configuration
	t             T
	httpClient    *ledgerclient.Formance
	cancel        func()
	ctx           context.Context
	errorChan     chan error
}

func (s *Server) Start() {
	s.t.Helper()

	tmpDir := s.t.TempDir()
	require.NoError(s.t, os.MkdirAll(tmpDir, 0700))
	s.t.Cleanup(func() {
		_ = os.RemoveAll(tmpDir)
	})

	rootCmd := cmd.NewRootCommand()
	args := []string{
		"serve",
		"--" + cmd.BindFlag, ":0",
		"--" + bunconnect.PostgresURIFlag, s.configuration.PostgresConfiguration.DatabaseSourceName,
		"--" + bunconnect.PostgresMaxOpenConnsFlag, fmt.Sprint(s.configuration.PostgresConfiguration.MaxOpenConns),
		"--" + bunconnect.PostgresConnMaxIdleTimeFlag, fmt.Sprint(s.configuration.PostgresConfiguration.ConnMaxIdleTime),
	}
	if s.configuration.PostgresConfiguration.MaxIdleConns != 0 {
		args = append(
			args,
			"--"+bunconnect.PostgresMaxIdleConnsFlag,
			fmt.Sprint(s.configuration.PostgresConfiguration.MaxIdleConns),
		)
	}
	if s.configuration.PostgresConfiguration.MaxOpenConns != 0 {
		args = append(
			args,
			"--"+bunconnect.PostgresMaxOpenConnsFlag,
			fmt.Sprint(s.configuration.PostgresConfiguration.MaxOpenConns),
		)
	}
	if s.configuration.PostgresConfiguration.ConnMaxIdleTime != 0 {
		args = append(
			args,
			"--"+bunconnect.PostgresConnMaxIdleTimeFlag,
			fmt.Sprint(s.configuration.PostgresConfiguration.ConnMaxIdleTime),
		)
	}
	if s.configuration.Debug {
		args = append(args, "--"+service.DebugFlag)
	}

	s.t.Logf("Starting application with flags: %s", strings.Join(args, " "))
	rootCmd.SetArgs(args)
	rootCmd.SilenceErrors = true
	output := s.configuration.Output
	if output == nil {
		output = io.Discard
	}
	rootCmd.SetOut(output)
	rootCmd.SetErr(output)

	s.ctx = logging.TestingContext()
	s.ctx, s.cancel = context.WithCancel(s.ctx)
	s.ctx = service.ContextWithLifecycle(s.ctx)
	s.ctx = httpserver.ContextWithServerInfo(s.ctx)

	s.errorChan = make(chan error, 1)
	go func() {
		s.errorChan <- rootCmd.ExecuteContext(s.ctx)
	}()

	select {
	case <-service.Ready(s.ctx):
	case err := <-s.errorChan:
		if err != nil {
			require.NoError(s.t, err)
		} else {
			require.Fail(s.t, "unexpected service stop")
		}
	}

	s.httpClient = ledgerclient.New(
		ledgerclient.WithServerURL(httpserver.URL(s.ctx)),
		ledgerclient.WithClient(&http.Client{
			Transport: httpclient.NewDebugHTTPTransport(http.DefaultTransport),
		}),
	)
}

func (s *Server) Stop() {
	s.t.Helper()

	if s.cancel == nil {
		return
	}
	s.cancel()
	s.cancel = nil

	// Wait app to be marked as stopped
	select {
	case <-service.Stopped(s.ctx):
	case <-time.After(5 * time.Second):
		require.Fail(s.t, "service should have been stopped")
	}

	// Ensure the app has been properly shutdown
	select {
	case err := <-s.errorChan:
		require.NoError(s.t, err)
	case <-time.After(5 * time.Second):
		require.Fail(s.t, "service should have been stopped without error")
	}
}

func (s *Server) Client() *ledgerclient.Formance {
	return s.httpClient
}

func (s *Server) Restart() {
	s.t.Helper()

	s.Stop()
	s.Start()
}

func New(t T, configuration Configuration) *Server {
	srv := &Server{
		t:             t,
		configuration: configuration,
	}
	t.Logf("Start testing server")
	srv.Start()
	t.Cleanup(func() {
		t.Logf("Stop testing server")
		srv.Stop()
	})

	return srv
}
