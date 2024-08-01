package models

import (
	"encoding/json"
	"math/big"
	"time"
)

type Payment struct {
	// PSP payment/transaction reference. Should be unique.
	Reference string

	// Payment Creation date.
	CreatedAt time.Time

	// Type of payment: payin, payout, transfer etc...
	PaymentType PaymentType

	// Payment amount.
	Amount *big.Int

	// Currency. Should be in minor currencies unit.
	// For example: USD/2
	Asset string

	// Payment scheme if existing: visa, mastercard etc...
	Scheme PaymentScheme

	// Payment status: pending, failed, succeeded etc...
	Status PaymentStatus

	// Optional, can be filled for payouts and transfers for example.
	SourceAccountReference *string
	// Optional, can be filled for payins and transfers for example.
	DestinationAccountReference *string

	// Additional metadata
	Metadata map[string]string

	// PSP response in raw
	Raw json.RawMessage
}

// TODO(polo): match grpc et const
type (
	PaymentType   int
	PaymentStatus int
	PaymentScheme int
)

const (
	PAYMENT_TYPE_UNKNOWN PaymentType = iota
	PAYMENT_TYPE_PAYIN
	PAYMENT_TYPE_PAYOUT
	PAYMENT_TYPE_TRANSFER
	PAYMENT_TYPE_OTHER = 100 // match grpc tag
)

const (
	PAYMENT_STATUS_UNKNOWN PaymentStatus = iota
	PAYMENT_STATUS_PENDING
	PAYMENT_STATUS_SUCCEEDED
	PAYMENT_STATUS_CANCELLED
	PAYMENT_STATUS_FAILED
	PAYMENT_STATUS_EXPIRED
	PAYMENT_STATUS_REFUNDED
	PAYMENT_STATUS_REFUNDED_FAILURE
	PAYMENT_STATUS_DISPUTE
	PAYMENT_STATUS_DISPUTE_WON
	PAYMENT_STATUS_DISPUTE_LOST
	PAYMENT_STATUS_OTHER = 100 // match grpc tag
)

const (
	PAYMENT_SCHEME_UNKNOWN PaymentScheme = iota

	PAYMENT_SCHEME_CARD_VISA
	PAYMENT_SCHEME_CARD_MASTERCARD
	PAYMENT_SCHEME_CARD_AMEX
	PAYMENT_SCHEME_CARD_DINERS
	PAYMENT_SCHEME_CARD_DISCOVER
	PAYMENT_SCHEME_CARD_JCB
	PAYMENT_SCHEME_CARD_UNION_PAY
	PAYMENT_SCHEME_CARD_ALIPAY
	PAYMENT_SCHEME_CARD_CUP

	PAYMENT_SCHEME_SEPA_DEBIT
	PAYMENT_SCHEME_SEPA_CREDIT
	PAYMENT_SCHEME_SEPA

	PAYMENT_SCHEME_GOOGLE_PAY
	PAYMENT_SCHEME_APPLE_PAY

	PAYMENT_SCHEME_DOKU
	PAYMENT_SCHEME_DRAGON_PAY
	PAYMENT_SCHEME_MAESTRO
	PAYMENT_SCHEME_MOL_PAY

	PaymentSchePAYMENT_SCHEME_A2A
	PAYMENT_SCHEME_ACH_DEBIT
	PAYMENT_SCHEME_ACH
	PAYMENT_SCHEME_RTP

	PAYMENT_SCHEME_OTHER = 100 // match grpc tag
)
