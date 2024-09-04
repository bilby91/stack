package ledger

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/formancehq/ledger/internal/bus"
	"github.com/formancehq/ledger/internal/controller/ledger/writer"
	"github.com/formancehq/ledger/internal/opentelemetry/metrics"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"go.uber.org/fx"
)

type ModuleConfiguration struct {
	NSCacheConfiguration writer.CacheConfiguration
}

func NewFXModule(configuration ModuleConfiguration) fx.Option {
	return fx.Options(
		fx.Provide(func(
			storageDriver StorageDriver,
			publisher message.Publisher,
			metricsRegistry metrics.GlobalRegistry,
			logger logging.Logger,
		) *Resolver {
			options := []option{
				WithMessagePublisher(publisher),
				WithMetricsRegistry(metricsRegistry),
			}
			if configuration.NSCacheConfiguration.MaxCount != 0 {
				options = append(options, WithCompiler(writer.NewCachedCompiler(
					writer.NewDefaultCompiler(),
					configuration.NSCacheConfiguration,
				)))
			}
			return NewResolver(storageDriver, options...)
		}),
		fx.Provide(fx.Annotate(bus.NewNoOpMonitor, fx.As(new(bus.Monitor)))),
		fx.Provide(fx.Annotate(metrics.NewNoOpRegistry, fx.As(new(metrics.GlobalRegistry)))),
	)
}
