package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/formancehq/stack/libs/go-libs/auth"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/formancehq/stack/libs/go-libs/health"
	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/go-chi/chi/v5/middleware"

	ingester "github.com/formancehq/stack/ee/ingester/internal"
	"github.com/formancehq/stack/ee/ingester/internal/controller"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/go-chi/chi/v5"
)

//go:generate mockgen -source api.go -destination api_generated.go -package api . Backend
type Backend interface {
	ListPipelines(ctx context.Context) (*bunpaginate.Cursor[ingester.Pipeline], error)
	GetPipeline(ctx context.Context, id string) (*ingester.Pipeline, error)
	CreatePipeline(ctx context.Context, pipelineConfiguration ingester.PipelineConfiguration) (*ingester.Pipeline, error)
	DeletePipeline(ctx context.Context, id string) error
	StartPipeline(ctx context.Context, id string) error
	PausePipeline(ctx context.Context, id string) error
	ResumePipeline(ctx context.Context, id string) error
	ResetPipeline(ctx context.Context, id string) error
	StopPipeline(ctx context.Context, id string) error

	ListConnectors(ctx context.Context) (*bunpaginate.Cursor[ingester.Connector], error)
	CreateConnector(ctx context.Context, configuration ingester.ConnectorConfiguration) (*ingester.Connector, error)
	DeleteConnector(ctx context.Context, id string) error
	GetConnector(ctx context.Context, id string) (*ingester.Connector, error)
}

type ErrModuleNotAvailable = controller.ErrModuleNotAvailable
type ErrConnectorNotFound = controller.ErrConnectorNotFound
type ErrPipelineNotFound = controller.ErrPipelineNotFound
type ErrPipelineAlreadyExists = controller.ErrPipelineAlreadyExists
type ErrInvalidStateSwitch = controller.ErrInvalidStateSwitch
type ErrPipelineAlreadyStarted = controller.ErrAlreadyStarted
type ErrInUsePipeline = controller.ErrInUsePipeline
type ErrInvalidConnectorConfiguration = controller.ErrInvalidDriverConfiguration
type ErrConnectorUsed = controller.ErrConnectorUsed

type API struct {
	backend          Backend
	logger           logging.Logger
	healthController *health.HealthController
	authenticator    auth.Authenticator
	serviceInfo      api.ServiceInfo
}

func (a *API) Router() chi.Router {
	ret := chi.NewMux()
	ret.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			r = r.WithContext(logging.ContextWithLogger(r.Context(), a.logger))
			handler.ServeHTTP(w, r)
		})
	})

	ret.Get("/_info", api.InfoHandler(a.serviceInfo))
	ret.Get("/_healthcheck", a.healthController.Check)
	ret.Route("/", func(r chi.Router) {
		r.Use(service.OTLPMiddleware("ingester", a.serviceInfo.Debug))
		r.Use(middleware.RequestLogger(api.NewLogFormatter()))
		r.Use(auth.Middleware(a.authenticator))
		r.Route("/pipelines", func(r chi.Router) {
			r.Get("/", a.listPipelines)
			r.Post("/", a.createPipeline)
			r.Route("/{pipelineID}", func(r chi.Router) {
				r.Get("/", a.readPipeline)
				r.Delete("/", a.deletePipeline)
				r.Post("/start", a.startPipeline)
				r.Post("/stop", a.stopPipeline)
				r.Post("/reset", a.resetPipeline)
				r.Post("/pause", a.pausePipeline)
				r.Post("/resume", a.resumePipeline)
			})
		})
		ret.Route("/connectors", func(r chi.Router) {
			r.Get("/", a.listConnectors)
			r.Post("/", a.createConnector)
			r.Route("/{connectorID}", func(r chi.Router) {
				r.Delete("/", a.deleteConnector)
				r.Get("/", a.getConnector)
			})
		})
	})

	return ret
}

func (a *API) pipelineID(r *http.Request) string {
	return chi.URLParam(r, "pipelineID")
}

func (a *API) connectorID(r *http.Request) string {
	return chi.URLParam(r, "connectorID")
}

func NewAPI(
	backend Backend,
	healthController *health.HealthController,
	authenticator auth.Authenticator,
	logger logging.Logger,
	serviceInfo api.ServiceInfo,
) *API {
	return &API{
		backend:          backend,
		logger:           logger,
		serviceInfo:      serviceInfo,
		authenticator:    authenticator,
		healthController: healthController,
	}
}

func withBody[V any](w http.ResponseWriter, r *http.Request, fn func(v V)) {
	var v V
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		api.BadRequest(w, "VALIDATION", err)
		return
	}

	fn(v)
}
