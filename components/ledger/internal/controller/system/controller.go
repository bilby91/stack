package system

import (
	"context"
	ledger "github.com/formancehq/ledger/internal"
	ledgercontroller "github.com/formancehq/ledger/internal/controller/ledger"
	"github.com/formancehq/ledger/internal/opentelemetry/tracer"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
)

type Controller struct {
	resolver *ledgercontroller.Resolver
	store    Store
}

func (c *Controller) GetLedgerController(ctx context.Context, name string) (*ledgercontroller.Controller, error) {
	return tracer.Trace(ctx, "GetLedgerController", func(ctx context.Context) (*ledgercontroller.Controller, error) {
		return c.resolver.GetLedger(ctx, name)
	})
}

func (c *Controller) CreateLedger(ctx context.Context, name string, configuration ledger.Configuration) error {
	return tracer.SkipResult(tracer.Trace(ctx, "CreateLedger", tracer.NoResult(func(ctx context.Context) error {
		return c.resolver.CreateLedger(ctx, name, configuration)
	})))
}

func (c *Controller) GetLedger(ctx context.Context, name string) (*ledger.Ledger, error) {
	return tracer.Trace(ctx, "GetLedger", func(ctx context.Context) (*ledger.Ledger, error) {
		return c.store.GetLedger(ctx, name)
	})
}

func (c *Controller) ListLedgers(ctx context.Context, query ListLedgersQuery) (*bunpaginate.Cursor[ledger.Ledger], error) {
	return tracer.Trace(ctx, "ListLedgers", func(ctx context.Context) (*bunpaginate.Cursor[ledger.Ledger], error) {
		return c.store.ListLedgers(ctx, query)
	})
}

func (c *Controller) UpdateLedgerMetadata(ctx context.Context, name string, m map[string]string) error {
	return tracer.SkipResult(tracer.Trace(ctx, "UpdateLedgerMetadata", tracer.NoResult(func(ctx context.Context) error {
		return c.store.UpdateLedgerMetadata(ctx, name, m)
	})))
}

func (c *Controller) DeleteLedgerMetadata(ctx context.Context, param string, key string) error {
	return tracer.SkipResult(tracer.Trace(ctx, "DeleteLedgerMetadata", tracer.NoResult(func(ctx context.Context) error {
		return c.store.DeleteLedgerMetadata(ctx, param, key)
	})))
}

func NewController(resolver *ledgercontroller.Resolver, store Store) *Controller {
	return &Controller{
		resolver: resolver,
		store:    store,
	}
}
