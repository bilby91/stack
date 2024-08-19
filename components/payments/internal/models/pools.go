package models

import (
	"time"

	"github.com/google/uuid"
)

type PoolAccounts struct {
	PoolID    uuid.UUID `json:"poolID"`
	AccountID AccountID `json:"accountID"`
}

type Pool struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`

	PoolAccounts []PoolAccounts `json:"poolAccounts"`
}
