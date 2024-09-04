package ledger

import (
	"context"
	ledger "github.com/formancehq/ledger/internal"
	"github.com/formancehq/ledger/internal/controller/ledger/writer"
	"github.com/pkg/errors"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/formancehq/ledger/internal/opentelemetry/metrics"
)

type option func(r *Resolver)

func WithMessagePublisher(publisher message.Publisher) option {
	return func(r *Resolver) {
		r.publisher = publisher
	}
}

func WithMetricsRegistry(registry metrics.GlobalRegistry) option {
	return func(r *Resolver) {
		r.metricsRegistry = registry
	}
}

func WithCompiler(compiler writer.Compiler) option {
	return func(r *Resolver) {
		r.compiler = compiler
	}
}

var defaultOptions = []option{
	WithMetricsRegistry(metrics.NewNoOpRegistry()),
	WithCompiler(writer.NewDefaultCompiler()),
}

type Resolver struct {
	storageDriver   StorageDriver
	metricsRegistry metrics.GlobalRegistry
	compiler        writer.Compiler
	publisher       message.Publisher
}

func NewResolver(storageDriver StorageDriver, options ...option) *Resolver {
	r := &Resolver{
		storageDriver: storageDriver,
	}
	for _, opt := range append(defaultOptions, options...) {
		opt(r)
	}

	return r
}

func (r *Resolver) GetLedger(ctx context.Context, name string) (*Controller, error) {
	if name == "" {
		return nil, errors.New("empty name")
	}

	store, err := r.storageDriver.OpenLedger(ctx, name)
	if err != nil {
		return nil, err
	}

	// todo: add only once
	//r.metricsRegistry.ActiveLedgers().Add(ctx, +1)

	return New(
		name,
		store,
		r.publisher,
		writer.NewDefaultMachineFactory(r.compiler, store),
	), nil
}

func (r *Resolver) CreateLedger(ctx context.Context, name string, configuration ledger.Configuration) error {
	if name == "" {
		return errors.New("empty name")
	}

	return r.storageDriver.CreateLedger(ctx, name, configuration)
}
