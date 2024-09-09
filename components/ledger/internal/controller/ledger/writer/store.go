package writer

import (
	"context"
	"database/sql"
	ledger "github.com/formancehq/ledger/internal"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/uptrace/bun"
	"math/big"
)

//go:generate mockgen -source store.go -destination store_generated.go -package writer . TX
type TX interface {
	LockAccounts(ctx context.Context, accounts ...string) error
	InsertTransaction(ctx context.Context, transaction ledger.TransactionData) (*ledger.Transaction, error)
	InsertMoves(ctx context.Context, move ...ledger.Move) error
	UpsertAccount(ctx context.Context, account ledger.Account) (bool, error)
	// RevertTransaction revert the transaction with identifier id
	// it returns :
	//  * the reverted transaction
	//  * a boolean indicating if the transaction has been reverted. false indicates an already reverted transaction (unless error != nil)
	//  * an error
	RevertTransaction(ctx context.Context, id int) (*ledger.Transaction, bool, error)
	UpdateTransactionMetadata(ctx context.Context, transactionID int, m metadata.Metadata) (*ledger.Transaction, error)
	DeleteTransactionMetadata(ctx context.Context, transactionID int, key string) (*ledger.Transaction, error)
	UpdateAccountMetadata(ctx context.Context, address string, m metadata.Metadata) error
	DeleteAccountMetadata(ctx context.Context, address, key string) error
	InsertLog(ctx context.Context, log ledger.Log) (*ledger.ChainedLog, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	AddToBalance(ctx context.Context, addr, asset string, amount *big.Int) (*big.Int, error)
}

//go:generate mockgen -source store.go -destination store_generated.go -package writer . Store
type Store interface {
	BeginTX(ctx context.Context) (TX, error)
	GetDB() bun.IDB
}
