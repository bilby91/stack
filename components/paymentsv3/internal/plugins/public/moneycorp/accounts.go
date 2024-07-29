package moneycorp

import (
	"context"
	"encoding/json"
	"time"

	"github.com/formancehq/stack/components/paymentsv3/internal/plugins/models"
)

type accountsState struct {
	LastPage int `json:"lastPage"`
	// Moneycorp does not send the creation date for accounts, but we can still
	// sort by ID created (which is incremental when creating accounts).
	LastIDCreated string `json:"lastIDCreated"`
}

func (p Plugin) fetchAccounts(ctx context.Context, req models.FetchAccountsRequest) (models.FetchAccountsResponse, error) {
	var oldState accountsState
	if req.State != nil {
		if err := json.Unmarshal(req.State, &oldState); err != nil {
			return models.FetchAccountsResponse{}, err
		}
	}

	newState := accountsState{
		LastPage:      oldState.LastPage,
		LastIDCreated: oldState.LastIDCreated,
	}

	var accounts []models.Account
	for page := oldState.LastPage; ; page++ {
		newState.LastPage = page

		pagedAccounts, err := p.client.GetAccounts(ctx, page)
		if err != nil {
			return models.FetchAccountsResponse{}, err
		}

		if len(pagedAccounts) == 0 {
			break
		}

		for _, account := range pagedAccounts {
			if account.ID <= oldState.LastIDCreated {
				continue
			}

			raw, err := json.Marshal(account)
			if err != nil {
				return models.FetchAccountsResponse{}, err
			}

			accounts = append(accounts, models.Account{
				Reference: account.ID,
				// Moneycorp does not send the opening date of the account
				CreatedAt: time.Now().UTC(),
				Name:      &account.Attributes.AccountName,
				Raw:       raw,
			})

			newState.LastIDCreated = account.ID

			if len(pagedAccounts) < p.client.PageSize() {
				break
			}
		}
	}

	payload, err := json.Marshal(newState)
	if err != nil {
		return models.FetchAccountsResponse{}, err
	}

	return models.FetchAccountsResponse{
		Accounts: accounts,
		NewState: payload,
	}, nil
}
