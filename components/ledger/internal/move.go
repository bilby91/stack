package ledger

import (
	"github.com/formancehq/stack/libs/go-libs/time"
	"math/big"
)

type Move struct {
	IsSource   bool
	Account    string
	Amount     *big.Int
	Asset      string
	InsertedAt    time.Time
	EffectiveDate time.Time
	TransactionSeq int
}

func (m Move) GetAsset() string {
	return m.Asset
}

func (m Move) GetAccount() string {
	return m.Account
}
