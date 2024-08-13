package backend

import (
	"context"

	"github.com/formancehq/paymentsv3/internal/connectors/plugins"
	"github.com/formancehq/paymentsv3/internal/models"
	"github.com/formancehq/paymentsv3/internal/storage"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
)

//go:generate mockgen -source backend.go -destination backend_generated.go -package backend . Backend
type Backend interface {
	// Connectors
	GetConnectorConfigs() plugins.Configs
	ListConnectors(ctx context.Context, query storage.ListConnectorssQuery) (*bunpaginate.Cursor[models.Connector], error)
}
