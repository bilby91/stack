package services

import (
	"context"

	"github.com/formancehq/paymentsv3/internal/models"
	"github.com/formancehq/paymentsv3/internal/storage"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
)

func (s *Service) PoolsList(ctx context.Context, query storage.ListPoolsQuery) (*bunpaginate.Cursor[models.Pool], error) {
	return s.storage.PoolsList(ctx, query)
}
