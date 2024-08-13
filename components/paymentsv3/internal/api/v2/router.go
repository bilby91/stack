package v2

import (
	"github.com/formancehq/paymentsv3/internal/api/backend"
	"github.com/formancehq/stack/libs/go-libs/auth"
	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/go-chi/chi/v5"
)

func newRouter(backend backend.Backend, a auth.Auth) *chi.Mux {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(auth.Middleware(a))
		r.Use(service.OTLPMiddleware("payments"))
		r.Route("/connectors", func(r chi.Router) {
			r.Get("/", listConnectors(backend))
			r.Get("/configs", getConnectorsConfigs(backend))
		})
	})

	return r
}
