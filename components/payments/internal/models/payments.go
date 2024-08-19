package models

import (
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/gibson042/canonicaljson-go"
)

// Internal struct used by the plugins
type PSPPayment struct {
	// PSP payment/transaction reference. Should be unique.
	Reference string

	// Payment Creation date.
	CreatedAt time.Time

	// Type of payment: payin, payout, transfer etc...
	Type PaymentType

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

type Payment struct {
	// Unique Payment ID generated from payments information
	ID PaymentID `json:"id"`
	// Related Connector ID
	ConnectorID ConnectorID `json:"connectorID"`

	// PSP payment/transaction reference. Should be unique.
	Reference string `json:"reference"`

	// Payment Creation date.
	CreatedAt time.Time `json:"createdAt"`

	// Type of payment: payin, payout, transfer etc...
	Type PaymentType `json:"type"`

	// Payment Initial amount
	InitialAmount *big.Int `json:"initialAmount"`
	// Payment amount.
	Amount *big.Int `json:"amount"`

	// Currency. Should be in minor currencies unit.
	// For example: USD/2
	Asset string `json:"asset"`

	// Payment scheme if existing: visa, mastercard etc...
	Scheme PaymentScheme `json:"scheme"`

	// Payment status: pending, failed, succeeded etc...
	Status PaymentStatus `json:"status"`

	// Optional, can be filled for payouts and transfers for example.
	SourceAccountID *AccountID `json:"sourceAccountID"`
	// Optional, can be filled for payins and transfers for example.
	DestinationAccountID *AccountID `json:"destinationAccountID"`

	// Additional metadata
	Metadata map[string]string `json:"metadata"`

	// Related adjustment
	Adjustments []PaymentAdjustment `json:"adjustments"`
}

type PaymentAdjustment struct {
	// Unique ID of the payment adjustment
	ID PaymentAdjustmentID `json:"id"`
	// Related Payment ID
	PaymentID PaymentID `json:"paymentID"`

	// Creation date of the adjustment
	CreatedAt time.Time `json:"createdAt"`

	// Status of the payment adjustement
	Status PaymentStatus `json:"status"`

	// Optional
	// Amount moved
	Amount *big.Int `json:"amount"`
	// Optional
	// Asset related to amount
	Asset *string `json:"asset"`

	// Additional metadata
	Metadata map[string]string `json:"metadata"`
	// PSP response in raw
	Raw json.RawMessage `json:"raw"`
}

type PaymentReference struct {
	Reference string
	Type      PaymentType
}

type PaymentID struct {
	PaymentReference
	ConnectorID ConnectorID
}

func (pid *PaymentID) MarshalJSON() ([]byte, error) {
	return []byte(pid.String()), nil
}

func (pid PaymentID) String() string {
	data, err := canonicaljson.Marshal(pid)
	if err != nil {
		panic(err)
	}

	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(data)
}

func PaymentIDFromString(value string) (PaymentID, error) {
	ret := PaymentID{}
	data, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(value)
	if err != nil {
		return ret, err
	}
	err = canonicaljson.Unmarshal(data, &ret)
	if err != nil {
		return ret, err
	}

	return ret, nil
}

func MustPaymentIDFromString(value string) *PaymentID {
	data, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(value)
	if err != nil {
		panic(err)
	}
	ret := PaymentID{}
	err = canonicaljson.Unmarshal(data, &ret)
	if err != nil {
		panic(err)
	}

	return &ret
}

func (pid PaymentID) Value() (driver.Value, error) {
	return pid.String(), nil
}

func (pid *PaymentID) Scan(value interface{}) error {
	if value == nil {
		return errors.New("payment id is nil")
	}

	if s, err := driver.String.ConvertValue(value); err == nil {

		if v, ok := s.(string); ok {

			id, err := PaymentIDFromString(v)
			if err != nil {
				return fmt.Errorf("failed to parse paymentid %s: %v", v, err)
			}

			*pid = id
			return nil
		}
	}

	return fmt.Errorf("failed to scan paymentid: %v", value)
}

type PaymentAdjustmentID struct {
	PaymentID
	CreatedAt time.Time
	Status    PaymentStatus
}

func (pid *PaymentAdjustmentID) MarshalJSON() ([]byte, error) {
	return []byte(pid.String()), nil
}

func (pid PaymentAdjustmentID) String() string {
	data, err := canonicaljson.Marshal(pid)
	if err != nil {
		panic(err)
	}

	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(data)
}

func PaymentAdjustmentIDFromString(value string) (*PaymentAdjustmentID, error) {
	data, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(value)
	if err != nil {
		return nil, err
	}
	ret := PaymentAdjustmentID{}
	err = canonicaljson.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

func MustPaymentAdjustmentIDFromString(value string) *PaymentAdjustmentID {
	data, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(value)
	if err != nil {
		panic(err)
	}
	ret := PaymentAdjustmentID{}
	err = canonicaljson.Unmarshal(data, &ret)
	if err != nil {
		panic(err)
	}

	return &ret
}

func (pid PaymentAdjustmentID) Value() (driver.Value, error) {
	return pid.String(), nil
}

func (pid *PaymentAdjustmentID) Scan(value interface{}) error {
	if value == nil {
		return errors.New("payment adjustment id is nil")
	}

	if s, err := driver.String.ConvertValue(value); err == nil {

		if v, ok := s.(string); ok {

			id, err := PaymentAdjustmentIDFromString(v)
			if err != nil {
				return fmt.Errorf("failed to parse payment adjustment id %s: %v", v, err)
			}

			*pid = *id
			return nil
		}
	}

	return fmt.Errorf("failed to scan payment adjustement id: %v", value)
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

func (t PaymentType) String() string {
	switch t {
	case PAYMENT_TYPE_PAYIN:
		return "PAYIN"
	case PAYMENT_TYPE_PAYOUT:
		return "PAYOUT"
	case PAYMENT_TYPE_TRANSFER:
		return "TRANSFER"
	case PAYMENT_TYPE_OTHER:
		return "OTHER"
	default:
		return "UNKNOWN"
	}
}

func (t PaymentType) Value() (driver.Value, error) {
	switch t {
	case PAYMENT_TYPE_PAYIN:
		return "PAYIN", nil
	case PAYMENT_TYPE_PAYOUT:
		return "PAYOUT", nil
	case PAYMENT_TYPE_TRANSFER:
		return "TRANSFER", nil
	case PAYMENT_TYPE_OTHER:
		return "OTHER", nil
	default:
		return nil, fmt.Errorf("unknown payment type")
	}
}

func (t *PaymentType) Scan(value interface{}) error {
	if value == nil {
		return errors.New("payment type is nil")
	}

	s, err := driver.String.ConvertValue(value)
	if err != nil {
		return fmt.Errorf("failed to convert payment type")
	}

	v, ok := s.(string)
	if !ok {
		return fmt.Errorf("failed to cast payment type")
	}

	switch v {
	case "PAYIN":
		*t = PAYMENT_TYPE_PAYIN
	case "PAYOUT":
		*t = PAYMENT_TYPE_PAYOUT
	case "TRANSFER":
		*t = PAYMENT_TYPE_TRANSFER
	case "OTHER":
		*t = PAYMENT_TYPE_OTHER
	default:
		return fmt.Errorf("unknown payment type")
	}

	return nil
}

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

func (t PaymentStatus) String() string {
	switch t {
	case PAYMENT_STATUS_UNKNOWN:
		return "UNKNOWN"
	case PAYMENT_STATUS_PENDING:
		return "PENDING"
	case PAYMENT_STATUS_SUCCEEDED:
		return "SUCCEEDED"
	case PAYMENT_STATUS_CANCELLED:
		return "CANCELLED"
	case PAYMENT_STATUS_FAILED:
		return "FAILED"
	case PAYMENT_STATUS_EXPIRED:
		return "EXPIRED"
	case PAYMENT_STATUS_REFUNDED:
		return "REFUNDED"
	case PAYMENT_STATUS_REFUNDED_FAILURE:
		return "REFUNDED_FAILURE"
	case PAYMENT_STATUS_DISPUTE:
		return "DISPUTE"
	case PAYMENT_STATUS_DISPUTE_WON:
		return "DISPUTE_WON"
	case PAYMENT_STATUS_DISPUTE_LOST:
		return "DISPUTE_LOST"
	case PAYMENT_STATUS_OTHER:
		return "OTHER"
	default:
		return "UNKNOWN"
	}
}

func (t PaymentStatus) Value() (driver.Value, error) {
	switch t {
	case PAYMENT_STATUS_UNKNOWN:
		return "UNKNOWN", nil
	case PAYMENT_STATUS_PENDING:
		return "PENDING", nil
	case PAYMENT_STATUS_SUCCEEDED:
		return "SUCCEEDED", nil
	case PAYMENT_STATUS_CANCELLED:
		return "CANCELLED", nil
	case PAYMENT_STATUS_FAILED:
		return "FAILED", nil
	case PAYMENT_STATUS_EXPIRED:
		return "EXPIRED", nil
	case PAYMENT_STATUS_REFUNDED:
		return "REFUNDED", nil
	case PAYMENT_STATUS_REFUNDED_FAILURE:
		return "REFUNDED_FAILURE", nil
	case PAYMENT_STATUS_DISPUTE:
		return "DISPUTE", nil
	case PAYMENT_STATUS_DISPUTE_WON:
		return "DISPUTE_WON", nil
	case PAYMENT_STATUS_DISPUTE_LOST:
		return "DISPUTE_LOST", nil
	case PAYMENT_STATUS_OTHER:
		return "OTHER", nil
	default:
		return nil, fmt.Errorf("unknown payment status")
	}
}

func (t *PaymentStatus) Scan(value interface{}) error {
	if value == nil {
		return errors.New("payment status is nil")
	}

	s, err := driver.String.ConvertValue(value)
	if err != nil {
		return fmt.Errorf("failed to convert payment status")
	}

	v, ok := s.(string)
	if !ok {
		return fmt.Errorf("failed to cast payment status")
	}

	switch v {
	case "UNKNOWN":
		*t = PAYMENT_STATUS_UNKNOWN
	case "PENDING":
		*t = PAYMENT_STATUS_PENDING
	case "SUCCEEDED":
		*t = PAYMENT_STATUS_SUCCEEDED
	case "CANCELLED":
		*t = PAYMENT_STATUS_CANCELLED
	case "FAILED":
		*t = PAYMENT_STATUS_FAILED
	case "EXPIRED":
		*t = PAYMENT_STATUS_EXPIRED
	case "REFUNDED":
		*t = PAYMENT_STATUS_REFUNDED
	case "REFUNDED_FAILURE":
		*t = PAYMENT_STATUS_REFUNDED_FAILURE
	case "DISPUTE":
		*t = PAYMENT_STATUS_DISPUTE
	case "DISPUTE_WON":
		*t = PAYMENT_STATUS_DISPUTE_WON
	case "DISPUTE_LOST":
		*t = PAYMENT_STATUS_DISPUTE_LOST
	case "OTHER":
		*t = PAYMENT_STATUS_OTHER
	default:
		return fmt.Errorf("unknown payment status")
	}

	return nil
}

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

	PAYMENT_SCHEME_A2A
	PAYMENT_SCHEME_ACH_DEBIT
	PAYMENT_SCHEME_ACH
	PAYMENT_SCHEME_RTP

	PAYMENT_SCHEME_OTHER = 100 // match grpc tag
)

func (s PaymentScheme) String() string {
	switch s {
	case PAYMENT_SCHEME_UNKNOWN:
		return "UNKNOWN"
	case PAYMENT_SCHEME_CARD_VISA:
		return "CARD_VISA"
	case PAYMENT_SCHEME_CARD_MASTERCARD:
		return "CARD_MASTERCARD"
	case PAYMENT_SCHEME_CARD_AMEX:
		return "CARD_AMEX"
	case PAYMENT_SCHEME_CARD_DINERS:
		return "CARD_DINERS"
	case PAYMENT_SCHEME_CARD_DISCOVER:
		return "CARD_DISCOVER"
	case PAYMENT_SCHEME_CARD_JCB:
		return "CARD_JCB"
	case PAYMENT_SCHEME_CARD_UNION_PAY:
		return "CARD_UNION_PAY"
	case PAYMENT_SCHEME_CARD_ALIPAY:
		return "CARD_ALIPAY"
	case PAYMENT_SCHEME_CARD_CUP:
		return "CARD_CUP"
	case PAYMENT_SCHEME_SEPA_DEBIT:
		return "SEPA_DEBIT"
	case PAYMENT_SCHEME_SEPA_CREDIT:
		return "SEPA_CREDIT"
	case PAYMENT_SCHEME_SEPA:
		return "SEPA"
	case PAYMENT_SCHEME_GOOGLE_PAY:
		return "GOOGLE_PAY"
	case PAYMENT_SCHEME_APPLE_PAY:
		return "APPLE_PAY"
	case PAYMENT_SCHEME_DOKU:
		return "DOKU"
	case PAYMENT_SCHEME_DRAGON_PAY:
		return "DRAGON_PAY"
	case PAYMENT_SCHEME_MAESTRO:
		return "MAESTRO"
	case PAYMENT_SCHEME_MOL_PAY:
		return "MOL_PAY"
	case PAYMENT_SCHEME_A2A:
		return "A2A"
	case PAYMENT_SCHEME_ACH_DEBIT:
		return "ACH_DEBIT"
	case PAYMENT_SCHEME_ACH:
		return "ACH"
	case PAYMENT_SCHEME_RTP:
		return "RTP"
	case PAYMENT_SCHEME_OTHER:
		return "OTHER"
	default:
		return "UNKNOWN"
	}
}

func (s PaymentScheme) Value() (driver.Value, error) {
	switch s {
	case PAYMENT_SCHEME_UNKNOWN:
		return "UNKNOWN", nil
	case PAYMENT_SCHEME_CARD_VISA:
		return "CARD_VISA", nil
	case PAYMENT_SCHEME_CARD_MASTERCARD:
		return "CARD_MASTERCARD", nil
	case PAYMENT_SCHEME_CARD_AMEX:
		return "CARD_AMEX", nil
	case PAYMENT_SCHEME_CARD_DINERS:
		return "CARD_DINERS", nil
	case PAYMENT_SCHEME_CARD_DISCOVER:
		return "CARD_DISCOVER", nil
	case PAYMENT_SCHEME_CARD_JCB:
		return "CARD_JCB", nil
	case PAYMENT_SCHEME_CARD_UNION_PAY:
		return "CARD_UNION_PAY", nil
	case PAYMENT_SCHEME_CARD_ALIPAY:
		return "CARD_ALIPAY", nil
	case PAYMENT_SCHEME_CARD_CUP:
		return "CARD_CUP", nil
	case PAYMENT_SCHEME_SEPA_DEBIT:
		return "SEPA_DEBIT", nil
	case PAYMENT_SCHEME_SEPA_CREDIT:
		return "SEPA_CREDIT", nil
	case PAYMENT_SCHEME_SEPA:
		return "SEPA", nil
	case PAYMENT_SCHEME_GOOGLE_PAY:
		return "GOOGLE_PAY", nil
	case PAYMENT_SCHEME_APPLE_PAY:
		return "APPLE_PAY", nil
	case PAYMENT_SCHEME_DOKU:
		return "DOKU", nil
	case PAYMENT_SCHEME_DRAGON_PAY:
		return "DRAGON_PAY", nil
	case PAYMENT_SCHEME_MAESTRO:
		return "MAESTRO", nil
	case PAYMENT_SCHEME_MOL_PAY:
		return "MOL_PAY", nil
	case PAYMENT_SCHEME_A2A:
		return "A2A", nil
	case PAYMENT_SCHEME_ACH_DEBIT:
		return "ACH_DEBIT", nil
	case PAYMENT_SCHEME_ACH:
		return "ACH", nil
	case PAYMENT_SCHEME_RTP:
		return "RTP", nil
	case PAYMENT_SCHEME_OTHER:
		return "OTHER", nil
	default:
		return nil, fmt.Errorf("unknown payment type")
	}
}

func (s *PaymentScheme) Scan(value interface{}) error {
	if value == nil {
		return errors.New("payment type is nil")
	}

	st, err := driver.String.ConvertValue(value)
	if err != nil {
		return fmt.Errorf("failed to convert payment type")
	}

	v, ok := st.(string)
	if !ok {
		return fmt.Errorf("failed to cast payment type")
	}

	switch v {
	case "UNKNOWN":
		*s = PAYMENT_SCHEME_UNKNOWN
	case "CARD_VISA":
		*s = PAYMENT_SCHEME_CARD_VISA
	case "CARD_MASTERCARD":
		*s = PAYMENT_SCHEME_CARD_MASTERCARD
	case "CARD_AMEX":
		*s = PAYMENT_SCHEME_CARD_AMEX
	case "CARD_DINERS":
		*s = PAYMENT_SCHEME_CARD_DINERS
	case "CARD_DISCOVER":
		*s = PAYMENT_SCHEME_CARD_DISCOVER
	case "CARD_JCB":
		*s = PAYMENT_SCHEME_CARD_JCB
	case "CARD_UNION_PAY":
		*s = PAYMENT_SCHEME_CARD_UNION_PAY
	case "CARD_ALIPAY":
		*s = PAYMENT_SCHEME_CARD_ALIPAY
	case "CARD_CUP":
		*s = PAYMENT_SCHEME_CARD_CUP
	case "SEPA_DEBIT":
		*s = PAYMENT_SCHEME_SEPA_DEBIT
	case "SEPA_CREDIT":
		*s = PAYMENT_SCHEME_SEPA_CREDIT
	case "SEPA":
		*s = PAYMENT_SCHEME_SEPA
	case "GOOGLE_PAY":
		*s = PAYMENT_SCHEME_GOOGLE_PAY
	case "APPLE_PAY":
		*s = PAYMENT_SCHEME_APPLE_PAY
	case "DOKU":
		*s = PAYMENT_SCHEME_DOKU
	case "DRAGON_PAY":
		*s = PAYMENT_SCHEME_DRAGON_PAY
	case "MAESTRO":
		*s = PAYMENT_SCHEME_MAESTRO
	case "MOL_PAY":
		*s = PAYMENT_SCHEME_MOL_PAY
	case "A2A":
		*s = PAYMENT_SCHEME_A2A
	case "ACH_DEBIT":
		*s = PAYMENT_SCHEME_ACH_DEBIT
	case "ACH":
		*s = PAYMENT_SCHEME_ACH
	case "RTP":
		*s = PAYMENT_SCHEME_RTP
	case "OTHER":
		*s = PAYMENT_SCHEME_OTHER
	default:
		return fmt.Errorf("unknown payment type")
	}

	return nil
}

func FromPSPPaymentToPayment(from PSPPayment, connectorID ConnectorID) Payment {
	return Payment{
		ID: PaymentID{
			PaymentReference: PaymentReference{
				Reference: from.Reference,
				Type:      from.Type,
			},
			ConnectorID: connectorID,
		},
		ConnectorID:   connectorID,
		Reference:     from.Reference,
		CreatedAt:     from.CreatedAt,
		Type:          from.Type,
		InitialAmount: from.Amount,
		Amount:        from.Amount,
		Asset:         from.Asset,
		Scheme:        from.Scheme,
		Status:        from.Status,
		SourceAccountID: func() *AccountID {
			if from.SourceAccountReference == nil {
				return nil
			}
			return &AccountID{
				Reference:   *from.SourceAccountReference,
				ConnectorID: connectorID,
			}
		}(),
		DestinationAccountID: func() *AccountID {
			if from.DestinationAccountReference == nil {
				return nil
			}
			return &AccountID{
				Reference:   *from.DestinationAccountReference,
				ConnectorID: connectorID,
			}
		}(),
		Metadata: from.Metadata,
	}
}

func FromPSPPayments(from []PSPPayment, connectorID ConnectorID) []Payment {
	payments := make([]Payment, 0, len(from))
	for _, p := range from {
		payment := FromPSPPaymentToPayment(p, connectorID)
		payment.Adjustments = append(payment.Adjustments, FromPSPPaymentToPaymentAdjustement(p, connectorID))
		payments = append(payments, payment)
	}
	return payments
}

func FromPSPPaymentToPaymentAdjustement(from PSPPayment, connectorID ConnectorID) PaymentAdjustment {
	paymentID := PaymentID{
		PaymentReference: PaymentReference{
			Reference: from.Reference,
			Type:      from.Type,
		},
		ConnectorID: connectorID,
	}

	return PaymentAdjustment{
		ID: PaymentAdjustmentID{
			PaymentID: paymentID,
			CreatedAt: from.CreatedAt,
			Status:    from.Status,
		},
		PaymentID: paymentID,
		CreatedAt: from.CreatedAt,
		Status:    from.Status,
		Amount:    from.Amount,
		Asset:     &from.Asset,
		Metadata:  from.Metadata,
		Raw:       from.Raw,
	}
}
