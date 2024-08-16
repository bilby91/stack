package services

import (
	"context"

	"github.com/formancehq/paymentsv3/internal/models"
)

func (s *Service) PoolsCreate(ctx context.Context, pool models.Pool) error {
	return s.storage.PoolsUpsert(ctx, pool)
}
