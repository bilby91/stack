package controller

import (
	"context"
	"database/sql"
	"encoding/json"
	"sync"

	ingester "github.com/formancehq/stack/ee/ingester/internal"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"

	"github.com/formancehq/stack/ee/ingester/internal/runner"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/pkg/errors"
)

//go:generate mockgen -source controller.go -destination controller_generated.go -package controller . ConfigValidator
type ConfigValidator interface {
	ValidateConfig(connectorName string, rawConnectorConfig json.RawMessage) error
}

type Controller struct {
	mu             sync.Mutex
	inUsePipelines collectionutils.Set[string]

	runner          Runner
	store           Store
	configValidator ConfigValidator
	logger          logging.Logger
}

func (ctrl *Controller) ListConnectors(ctx context.Context) (*bunpaginate.Cursor[ingester.Connector], error) {
	return ctrl.store.ListConnectors(ctx)
}

// CreateConnector can return following errors:
// * ErrInvalidDriverConfiguration
func (ctrl *Controller) CreateConnector(ctx context.Context, configuration ingester.ConnectorConfiguration) (*ingester.Connector, error) {

	if err := ctrl.configValidator.ValidateConfig(configuration.Driver, configuration.Config); err != nil {
		return nil, NewErrInvalidDriverConfiguration(configuration.Driver, err)
	}

	connector := ingester.NewConnector(configuration)
	if err := ctrl.store.CreateConnector(ctx, connector); err != nil {
		return nil, err
	}
	return &connector, nil
}

// DeleteConnector can return following errors:
// ErrConnectorNotFound
func (ctrl *Controller) DeleteConnector(ctx context.Context, id string) error {
	if err := ctrl.store.DeleteConnector(ctx, id); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return NewErrConnectorNotFound(id)
		default:
			return err
		}
	}
	return nil
}

// GetConnector can return following errors:
// ErrConnectorNotFound
func (ctrl *Controller) GetConnector(ctx context.Context, id string) (*ingester.Connector, error) {
	connector, err := ctrl.store.GetConnector(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, NewErrConnectorNotFound(id)
		default:
			return nil, err
		}
	}
	return connector, nil
}

func (ctrl *Controller) ListPipelines(ctx context.Context) (*bunpaginate.Cursor[ingester.Pipeline], error) {
	return ctrl.store.ListPipelines(ctx)
}

func (ctrl *Controller) GetPipeline(ctx context.Context, id string) (*ingester.Pipeline, error) {
	pipeline, err := ctrl.store.GetPipeline(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, runner.NewErrPipelineNotFound(id)
		}
		return nil, err
	}

	return pipeline, nil
}

// PausePipeline can return following errors:
// * ErrPipelineNotFound
// * ErrInvalidStateSwitch
// * ErrInUsePipeline
func (ctrl *Controller) PausePipeline(ctx context.Context, id string) error {
	return ctrl.callAndWaitStateOnPipeline(ctx, id, Pipeline.Pause, func(state ingester.State) bool {
		return state.Label == ingester.StateLabelPause
	})
}

// ResumePipeline can return following errors:
// * ErrPipelineNotFound
// * ErrInvalidStateSwitch
// * ErrInUsePipeline
func (ctrl *Controller) ResumePipeline(ctx context.Context, id string) error {
	return ctrl.callAndWaitStateOnPipeline(ctx, id, Pipeline.Resume, func(c ingester.State) bool {
		return c.Label != ingester.StateLabelPause
	})
}

// ResetPipeline can return following errors:
// * ErrPipelineNotFound
// * ErrInUsePipeline
func (ctrl *Controller) ResetPipeline(ctx context.Context, id string) error {
	return ctrl.callAndWaitStateOnPipeline(ctx, id, Pipeline.Reset, func(state ingester.State) bool {
		return state.Label == ingester.StateLabelInit
	})
}

// StopPipeline can return following errors:
// * ErrPipelineNotFound
// * ErrInUsePipeline
func (ctrl *Controller) StopPipeline(ctx context.Context, id string) error {
	originalError := ctrl.callAndWaitStateOnPipeline(ctx, id, Pipeline.Stop, func(state ingester.State) bool {
		return state.Label == ingester.StateLabelStop
	})
	// The method Controller.callAndWaitStateOnPipeline can return ErrPipelineNotFound because
	// the pipeline is not running, but the pipeline can exist in database.
	// So, we check its existence and map error if relevant
	if originalError != nil {
		if errors.Is(originalError, ErrPipelineNotFound("")) {
			pipeline, err := ctrl.store.GetPipeline(ctx, id)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return originalError
				}
				return err
			}
			return runner.NewErrInvalidStateSwitch(id, pipeline.State.Label, ingester.StateLabelStop)
		}
		return originalError
	}
	return nil
}

// DeletePipeline can return following errors:
// * ErrPipelineNotFound
// * ErrInUsePipeline
// * ErrConnectorUsed
// The method will stop the pipeline if it is actually running,
// then it will delete the pair from database
func (ctrl *Controller) DeletePipeline(ctx context.Context, id string) error {
	return ctrl.withPipelineLocked(id, func() error {
		if p, ok := ctrl.runner.GetPipeline(id); ok {
			if err := p.Stop(); err != nil {
				return errors.Wrap(err, "stopping pipeline")
			}
		}

		if err := ctrl.store.DeletePipeline(ctx, id); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return runner.NewErrPipelineNotFound(id)
			}
			return err
		}
		return nil
	})
}

// CreatePipeline can return following errors:
// * ErrModuleNotAvailable
// * ErrConnectorNotFound
// * ErrInUsePipeline
func (ctrl *Controller) CreatePipeline(ctx context.Context, pipelineConfiguration ingester.PipelineConfiguration) (*ingester.Pipeline, error) {
	ctrl.mu.Lock()
	defer ctrl.mu.Unlock()

	pipeline := ingester.NewPipeline(pipelineConfiguration, ingester.NewInitState())

	err := ctrl.store.CreatePipeline(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	if _, err := ctrl.runner.StartPipeline(ctx, pipeline); err != nil {
		switch {
		case errors.Is(err, runner.ErrConnectorNotFound("")):
			return nil, NewErrConnectorNotFound(pipelineConfiguration.ConnectorID)
		default:
			return nil, err
		}
	}

	return &pipeline, nil
}

// StartPipeline can return following errors:
// * ErrPipelineNotFound
// * ErrAlreadyStarted
// * ErrInUsePipeline
func (ctrl *Controller) StartPipeline(ctx context.Context, id string) error {
	return ctrl.withPipelineLocked(id, func() error {
		pipeline, err := ctrl.store.GetPipeline(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				return runner.NewErrPipelineNotFound(id)
			default:
				return err
			}
		}
		_, err = ctrl.runner.StartPipeline(ctx, *pipeline)
		return err
	})
}

func (ctrl *Controller) markPipelineInUse(id string) (func(), error) {
	ctrl.mu.Lock()
	if ctrl.inUsePipelines.Contains(id) {
		ctrl.mu.Unlock()
		return nil, NewErrInUsePipeline(id)
	}
	ctrl.inUsePipelines.Put(id)
	ctrl.mu.Unlock()

	return func() {
		ctrl.mu.Lock()
		ctrl.inUsePipelines.Remove(id)
		ctrl.mu.Unlock()
	}, nil
}

func (ctrl *Controller) withPipelineLocked(id string, fn func() error) error {
	release, err := ctrl.markPipelineInUse(id)
	if err != nil {
		return err
	}
	defer release()

	return fn()
}

func (ctrl *Controller) callAndWaitStateOnPipeline(
	ctx context.Context,
	id string,
	fn func(pipeline Pipeline) error,
	changeFilters ...runner.ChangerFilter[ingester.State],
) error {
	return ctrl.withPipelineLocked(id, func() error {
		p, ok := ctrl.runner.GetPipeline(id)
		if !ok {
			return runner.NewErrPipelineNotFound(id)
		}
		stateListener, cancelStateListener := p.GetActiveState().Listen(changeFilters...)
		defer cancelStateListener()

		if err := fn(p); err != nil {
			return err
		}

		select {
		case <-stateListener:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	})
}

func New(runner Runner, store Store, configValidator ConfigValidator, logger logging.Logger) *Controller {
	return &Controller{
		runner:          runner,
		store:           store,
		inUsePipelines:  collectionutils.NewSet[string](),
		logger:          logger,
		configValidator: configValidator,
	}
}
