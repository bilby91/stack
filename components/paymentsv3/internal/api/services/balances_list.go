package services

import (
	"context"

	"github.com/formancehq/paymentsv3/internal/models"
	"github.com/formancehq/paymentsv3/internal/storage"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
)

func (s *Service) BalancesList(ctx context.Context, query storage.ListBalancesQuery) (*bunpaginate.Cursor[models.Balance], error) {
	return s.storage.BalancesList(ctx, query)
}
