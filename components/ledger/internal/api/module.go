package api

import (
	_ "embed"
	"github.com/formancehq/ledger/internal/controller/system"
	"github.com/formancehq/stack/libs/go-libs/httpserver"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/formancehq/ledger/internal/api/backend"
	"github.com/formancehq/ledger/internal/opentelemetry/metrics"
	"github.com/formancehq/stack/libs/go-libs/auth"
	"github.com/formancehq/stack/libs/go-libs/health"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	"go.uber.org/fx"
)

type Config struct {
	Version string
	Debug   bool
	Bind    string
}

func Module(cfg Config) fx.Option {
	return fx.Options(
		fx.Provide(func(
			backend backend.Backend,
			healthController *health.HealthController,
			globalMetricsRegistry metrics.GlobalRegistry,
			authenticator auth.Authenticator,
			logger logging.Logger,
		) chi.Router {
			return NewRouter(
				backend,
				healthController,
				globalMetricsRegistry,
				authenticator,
				logger,
				cfg.Debug,
			)
		}),
		fx.Provide(func(systemController *system.Controller) backend.Backend {
			return backend.NewDefaultBackend(systemController, cfg.Version)
		}),
		fx.Provide(fx.Annotate(noop.NewMeterProvider, fx.As(new(metric.MeterProvider)))),
		fx.Decorate(fx.Annotate(func(meterProvider metric.MeterProvider) (metrics.GlobalRegistry, error) {
			return metrics.RegisterGlobalRegistry(meterProvider)
		}, fx.As(new(metrics.GlobalRegistry)))),
		health.Module(),
		fx.Invoke(func(lc fx.Lifecycle, h chi.Router, logger logging.Logger) {

			// todo: get middlewares used by the data ingester
			wrappedRouter := chi.NewRouter()
			wrappedRouter.Use(func(handler http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					r = r.WithContext(logging.ContextWithLogger(r.Context(), logger))
					handler.ServeHTTP(w, r)
				})
			})
			wrappedRouter.Mount("/", h)

			lc.Append(httpserver.NewHook(wrappedRouter, httpserver.WithAddress(cfg.Bind)))
		}),
	)
}
