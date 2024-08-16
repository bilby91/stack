package services

import (
	"context"

	"github.com/formancehq/paymentsv3/internal/models"
	"github.com/google/uuid"
)

func (s *Service) PoolsAddAccount(ctx context.Context, id uuid.UUID, accountID models.AccountID) error {
	return s.storage.PoolsAddAccount(ctx, id, accountID)
}
