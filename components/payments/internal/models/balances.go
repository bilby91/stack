package models

import (
	"encoding/json"
	"math/big"
	"time"
)

type Balance struct {
	AccountID     AccountID `json:"accountID"`
	CreatedAt     time.Time `json:"createdAt"`
	LastUpdatedAt time.Time `json:"lastUpdatedAt"`

	Asset   string   `json:"asset"`
	Balance *big.Int `json:"balance"`
}

func (b Balance) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		AccountID     string    `json:"accountID"`
		CreatedAt     time.Time `json:"createdAt"`
		LastUpdatedAt time.Time `json:"lastUpdatedAt"`

		Asset   string   `json:"asset"`
		Balance *big.Int `json:"balance"`
	}{
		AccountID:     b.AccountID.String(),
		CreatedAt:     b.CreatedAt,
		LastUpdatedAt: b.LastUpdatedAt,
		Asset:         b.Asset,
		Balance:       b.Balance,
	})
}

func (b *Balance) UnmarshalJSON(data []byte) error {
	var aux struct {
		AccountID     string    `json:"accountID"`
		CreatedAt     time.Time `json:"createdAt"`
		LastUpdatedAt time.Time `json:"lastUpdatedAt"`
		Asset         string    `json:"asset"`
		Balance       *big.Int  `json:"balance"`
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	accountID, err := AccountIDFromString(aux.AccountID)
	if err != nil {
		return err
	}

	b.AccountID = accountID
	b.CreatedAt = aux.CreatedAt
	b.LastUpdatedAt = aux.LastUpdatedAt
	b.Asset = aux.Asset
	b.Balance = aux.Balance

	return nil
}

type AggregatedBalance struct {
	Asset  string   `json:"asset"`
	Amount *big.Int `json:"amount"`
}
