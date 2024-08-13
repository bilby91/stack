package services

import (
	"context"

	"github.com/formancehq/paymentsv3/internal/models"
	"github.com/formancehq/paymentsv3/internal/storage"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
)

func (s *Service) ListConnectors(ctx context.Context, query storage.ListConnectorssQuery) (*bunpaginate.Cursor[models.Connector], error) {
	cursor, err := s.storage.ListConnectors(ctx, query)
	return cursor, newStorageError(err, "failed to list connectors")
}
