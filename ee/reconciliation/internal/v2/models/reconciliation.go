package models

import (
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type ReconciliationStatus string

const (
	ReconciliationStatusFailed    = "FAILED"
	ReconciliationStatusSucceeded = "SUCCEEDED"
)

type Reconciliation struct {
	bun.BaseModel `bun:"reconciliationsv2.reconciliations" json:"-"`

	ID                   uuid.UUID            `bun:",pk,notnull" json:"id"`
	Name                 string               `bun:",notnull" json:"name"`
	PolicyID             uuid.UUID            `bun:",notnull" json:"policyID"`
	PolicyType           PolicyType           `bun:",notnull" json:"policyType"`
	CreatedAt            time.Time            `bun:",notnull" json:"createdAt"`
	ReconciliationStatus ReconciliationStatus `bun:",notnull" json:"status"`
}

// Reconciliations accounts based are a one time reconciliation between
// ledger accounts and payments accounts.
// In postgres, it will be represented as a single row for each reconciliation.
type ReconciliationAccountBased struct {
	bun.BaseModel `bun:"reconciliationsv2.reconciliations_account_based" json:"-"`

	ID                   uuid.UUID            `bun:",pk,notnull" json:"id"`
	PolicyID             uuid.UUID            `bun:",notnull" json:"policyID"`
	ReconciledAtLedger   time.Time            `bun:",notnull" json:"reconciledAtLedger"`
	ReconciledAtPayments time.Time            `bun:",notnull" json:"reconciledAtPayments"`
	ReconciliationStatus ReconciliationStatus `bun:",notnull" json:"status"`
	LedgerBalances       map[string]*big.Int  `bun:",jsonb" json:"ledgerBalances"`
	PaymentsBalances     map[string]*big.Int  `bun:",jsonb" json:"paymentsBalances"`
	DriftBalances        map[string]*big.Int  `bun:",jsonb" json:"driftBalances"`
	Error                string               `json:"error"`
}

// For reconciliations transactions based, we need to store every succeeded
// reconcililed payment ID and transactions inside a table, one per row.
// Example:
// - If one payment ID corresponds to one ledger transaction, we will have one
// row with both ID and the reconciliation rule ID and matched the two of them.
// - If one payment ID corresponds to two or more ledger transactions, we will
// have one row for each transaction with the the same payment ID for both and
// the rule ID that match both of them.
// - If one ledger transaction corresponds to two or more payments IDs, we will
// have one row for each payment ID with the same transaction ID for both and
// the rule ID that match both of them.
type ReconciliationTransactionSucceeded struct {
	bun.BaseModel `bun:"reconciliationsv2.reconciliations_transactions_succeeded" json:"-"`

	PaymentID     string   `bun:",pk,notnull" json:"paymentID"`
	TransactionID *big.Int `bun:",pk,notnull" json:"transactionID"`
	RuleID        string   `bun:",pk,notnull" json:"ruleID"`
}

// We have a configurable delay when the transaction or payment is pending, so
// they will be stored in a table until the delay is reached.
type ReconciliationTransactionPending struct {
	bun.BaseModel `bun:"reconciliationsv2.reconciliations_transactions_pending" json:"-"`

	PaymentID     string    `json:"paymentID"`
	TransactionID *big.Int  `json:"transactionID"`
	PolicyID      uuid.UUID `json:"policyID"`
	CreatedAt     time.Time `json:"createdAt"`
}

// Failed transactions/paymenbts are stored inside a separated table, they
// corresponds to failed to reconcile transactions or payments after the
// configurable delay.
type ReconciliationTransactionFailed struct {
	bun.BaseModel `bun:"reconciliationsv2.reconciliations_transactions_failed" json:"-"`

	PaymentID     string    `json:"paymentID"`
	TransactionID *big.Int  `json:"transactionID"`
	PolicyID      uuid.UUID `json:"policyID"`
}
