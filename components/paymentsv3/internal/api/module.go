package api

import (
	"github.com/formancehq/paymentsv3/internal/api/backend"
	"github.com/formancehq/paymentsv3/internal/api/services"
	"github.com/formancehq/paymentsv3/internal/connectors/engine"
	"github.com/formancehq/paymentsv3/internal/storage"
	"go.uber.org/fx"
)

func TagVersion() fx.Annotation {
	return fx.ResultTags(`group:"apiVersions"`)
}

func NewModule() fx.Option {
	return fx.Options(
		fx.Provide(fx.Annotate(NewRouter, fx.ParamTags(``, ``, ``, ``, `group:"apiVersions"`))),
		fx.Provide(func(storage storage.Storage, engine engine.Engine) backend.Backend {
			return services.New(storage, engine)
		}),
	)
}
