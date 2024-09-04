package ledger

import (
	"context"
	"github.com/pkg/errors"
)

type Stats struct {
	Transactions int `json:"transactions"`
	Accounts     int `json:"accounts"`
}

func (l *Controller) Stats(ctx context.Context) (Stats, error) {
	var stats Stats

	transactions, err := l.store.CountTransactions(ctx, NewGetTransactionsQuery(NewPaginatedQueryOptions(PITFilterWithVolumes{})))
	if err != nil {
		return stats, errors.Wrap(err, "counting transactions")
	}

	accounts, err := l.store.CountAccounts(ctx, NewGetAccountsQuery(NewPaginatedQueryOptions(PITFilterWithVolumes{})))
	if err != nil {
		return stats, errors.Wrap(err, "counting accounts")
	}

	return Stats{
		Transactions: transactions,
		Accounts:     accounts,
	}, nil
}
