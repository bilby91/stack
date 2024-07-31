package moneycorp

import "github.com/formancehq/paymentsv3/internal/plugins/models"

var capabilities = []models.Capability{
	models.CAPABILITY_FETCH_ACCOUNTS,
	models.CAPABILITY_FETCH_PAYMENTS,
}
