package services

import (
	"context"

	"github.com/formancehq/paymentsv3/internal/models"
)

func (s *Service) PaymentsGet(ctx context.Context, id models.PaymentID) (*models.Payment, error) {
	return s.storage.PaymentsGet(ctx, id)
}
