package services

import (
	"context"

	"github.com/formancehq/paymentsv3/internal/models"
	"github.com/formancehq/paymentsv3/internal/storage"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
)

func (s *Service) PaymentsList(ctx context.Context, query storage.ListPaymentsQuery) (*bunpaginate.Cursor[models.Payment], error) {
	return s.storage.PaymentsList(ctx, query)
}
