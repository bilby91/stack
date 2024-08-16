package models

import (
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

type AggregatedBalance struct {
	Asset  string   `json:"asset"`
	Amount *big.Int `json:"amount"`
}
