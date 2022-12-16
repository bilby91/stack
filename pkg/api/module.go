package api

import (
	"context"
	"net/http"

	"github.com/formancehq/go-libs/sharedapi"
	sharedhealth "github.com/formancehq/go-libs/sharedhealth/pkg"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.uber.org/fx"
)

func CreateRootRouter() *mux.Router {
	rootRouter := mux.NewRouter()
	rootRouter.Use(otelmux.Middleware("auth"))
	rootRouter.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			handler.ServeHTTP(w, r)
		})
	})
	return rootRouter
}

func addInfoRoute(router *mux.Router, serviceInfo sharedapi.ServiceInfo) {
	router.Path("/_info").Methods(http.MethodGet).HandlerFunc(sharedapi.InfoHandler(serviceInfo))
}

func Module(addr string, serviceInfo sharedapi.ServiceInfo) fx.Option {
	return fx.Options(
		sharedhealth.ProvideHealthCheck(delegatedOIDCServerAvailability),
		sharedhealth.Module(),
		fx.Supply(serviceInfo),
		fx.Provide(CreateRootRouter),
		fx.Invoke(func(lc fx.Lifecycle, r *mux.Router, healthController *sharedhealth.HealthController) {
			finalRouter := mux.NewRouter()
			finalRouter.Path("/_healthcheck").HandlerFunc(healthController.Check)
			finalRouter.PathPrefix("/").Handler(r)
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					return StartServer(ctx, addr, finalRouter)
				},
			})
		}),
		fx.Invoke(
			addInfoRoute,
			addClientRoutes,
			addScopeRoutes,
			addUserRoutes,
		),
	)
}
