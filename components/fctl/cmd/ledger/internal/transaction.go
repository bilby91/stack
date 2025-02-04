package internal

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	formance "github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
)

func TransactionIDOrLastN(ctx context.Context, ledgerClient *formance.Formance, ledger, id string) (*big.Int, error) {
	if strings.HasPrefix(id, "last") {
		id = strings.TrimPrefix(id, "last")
		sub := int64(0)
		if id != "" {
			var err error
			sub, err = strconv.ParseInt(id, 10, 64)
			if err != nil {
				return nil, err
			}
		}
		pageSize := int64(1)
		request := operations.ListTransactionsRequest{
			Ledger:   ledger,
			PageSize: &pageSize,
		}
		response, err := ledgerClient.Ledger.V1.ListTransactions(ctx, request)
		if err != nil {
			return nil, err
		}

		if response.StatusCode >= 300 {
			return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
		}

		if len(response.TransactionsCursorResponse.Cursor.Data) == 0 {
			return nil, errors.New("no transaction found")
		}
		return response.TransactionsCursorResponse.Cursor.Data[0].Txid.Sub(
			response.TransactionsCursorResponse.Cursor.Data[0].Txid,
			big.NewInt(sub),
		), nil
	}

	v, ok := big.NewInt(0).SetString(id, 10)
	if !ok {
		return nil, fmt.Errorf("invalid bigint: %s", id)
	}

	return v, nil
}
