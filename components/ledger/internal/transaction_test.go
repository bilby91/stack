package ledger

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReverseTransaction(t *testing.T) {
	tx := NewTransactionData().WithPostings(
		NewPosting("world", "users:001", "COIN", big.NewInt(100)),
		NewPosting("users:001", "payments:001", "COIN", big.NewInt(100)),
	)

	expected := NewTransactionData().WithPostings(
		NewPosting("payments:001", "users:001", "COIN", big.NewInt(100)),
		NewPosting("users:001", "world", "COIN", big.NewInt(100)),
	)

	require.Equal(t, expected, tx.Reverse(false))
}
