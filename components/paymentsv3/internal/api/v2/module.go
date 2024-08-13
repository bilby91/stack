package v2

import (
	"github.com/formancehq/paymentsv3/internal/api"
	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Options(
		fx.Supply(fx.Annotate(api.Version{
			Version: 2,
			Builder: newRouter,
		}, api.TagVersion())),
	)
}
