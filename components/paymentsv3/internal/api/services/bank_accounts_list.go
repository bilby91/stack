package services

import (
	"context"

	"github.com/formancehq/paymentsv3/internal/models"
	"github.com/formancehq/paymentsv3/internal/storage"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
)

func (s *Service) BankAccountsList(ctx context.Context, query storage.ListBankAccountsQuery) (*bunpaginate.Cursor[models.BankAccount], error) {
	return s.storage.BankAccountsList(ctx, query)
}
