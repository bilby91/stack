package models

type (
	PaymentType   string
	PaymentStatus string
	PaymentScheme string
)

const (
	PaymentTypePayIn    PaymentType = "PAYIN"
	PaymentTypePayOut   PaymentType = "PAYOUT"
	PaymentTypeTransfer PaymentType = "TRANSFER"
	PaymentTypeOther    PaymentType = "OTHER"
)

const (
	PaymentStatusPending         PaymentStatus = "PENDING"
	PaymentStatusSucceeded       PaymentStatus = "SUCCEEDED"
	PaymentStatusCancelled       PaymentStatus = "CANCELLED"
	PaymentStatusFailed          PaymentStatus = "FAILED"
	PaymentStatusExpired         PaymentStatus = "EXPIRED"
	PaymentStatusRefunded        PaymentStatus = "REFUNDED"
	PaymentStatusRefundedFailure PaymentStatus = "REFUNDED_FAILURE"
	PaymentStatusDispute         PaymentStatus = "DISPUTE"
	PaymentStatusDisputeWon      PaymentStatus = "DISPUTE_WON"
	PaymentStatusDisputeLost     PaymentStatus = "DISPUTE_LOST"
	PaymentStatusOther           PaymentStatus = "OTHER"
)

const (
	PaymentSchemeUnknown PaymentScheme = "unknown"
	PaymentSchemeOther   PaymentScheme = "other"

	PaymentSchemeCardVisa       PaymentScheme = "visa"
	PaymentSchemeCardMasterCard PaymentScheme = "mastercard"
	PaymentSchemeCardAmex       PaymentScheme = "amex"
	PaymentSchemeCardDiners     PaymentScheme = "diners"
	PaymentSchemeCardDiscover   PaymentScheme = "discover"
	PaymentSchemeCardJCB        PaymentScheme = "jcb"
	PaymentSchemeCardUnionPay   PaymentScheme = "unionpay"
	PaymentSchemeCardAlipay     PaymentScheme = "alipay"
	PaymentSchemeCardCUP        PaymentScheme = "cup"

	PaymentSchemeSepaDebit  PaymentScheme = "sepa debit"
	PaymentSchemeSepaCredit PaymentScheme = "sepa credit"
	PaymentSchemeSepa       PaymentScheme = "sepa"

	PaymentSchemeApplePay  PaymentScheme = "apple pay"
	PaymentSchemeGooglePay PaymentScheme = "google pay"

	PaymentSchemeDOKU      PaymentScheme = "doku"
	PaymentSchemeDragonPay PaymentScheme = "dragonpay"
	PaymentSchemeMaestro   PaymentScheme = "maestro"
	PaymentSchemeMolPay    PaymentScheme = "molpay"

	PaymentSchemeA2A      PaymentScheme = "a2a"
	PaymentSchemeACHDebit PaymentScheme = "ach debit"
	PaymentSchemeACH      PaymentScheme = "ach"
	PaymentSchemeRTP      PaymentScheme = "rtp"
)
