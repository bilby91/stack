package v2

import (
	"encoding/json"
	"net/http"

	"github.com/formancehq/paymentsv3/internal/api/backend"
	"github.com/formancehq/paymentsv3/internal/storage"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
)

// NOTE: in order to maintain previous version compatibility, we need to keep the
// same response structure as the previous version of the API

func listConnectors(backend backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO(polo): add opentalemetry span
		ctx := r.Context()

		connectors, err := backend.ListConnectors(
			ctx,
			storage.NewListConnectorsQuery(
				bunpaginate.NewPaginatedQueryOptions(storage.ConnectorQuery{}).
					// NOTE: previous version of payments did not have pagination, so
					// fetch everything and return it all
					WithPageSize(1000),
			),
		)
		if err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		data := make([]readConnectorsResponseElement, len(connectors.Data))
		for i := range connectors.Data {
			data[i] = readConnectorsResponseElement{
				Provider:    data[i].Provider,
				ConnectorID: data[i].ConnectorID,
				Name:        data[i].Name,
				Enabled:     true,
			}
		}

		err = json.NewEncoder(w).Encode(
			api.BaseResponse[[]readConnectorsResponseElement]{
				Data: &data,
			})
		if err != nil {
			api.InternalServerError(w, r, err)
			return
		}
	}
}

type ConnectorProvider string

const (
	ConnectorProviderBankingCircle ConnectorProvider = "BANKING-CIRCLE"
	ConnectorProviderCurrencyCloud ConnectorProvider = "CURRENCY-CLOUD"
	ConnectorProviderDummyPay      ConnectorProvider = "DUMMY-PAY"
	ConnectorProviderModulr        ConnectorProvider = "MODULR"
	ConnectorProviderStripe        ConnectorProvider = "STRIPE"
	ConnectorProviderWise          ConnectorProvider = "WISE"
	ConnectorProviderMangopay      ConnectorProvider = "MANGOPAY"
	ConnectorProviderMoneycorp     ConnectorProvider = "MONEYCORP"
	ConnectorProviderAtlar         ConnectorProvider = "ATLAR"
	ConnectorProviderAdyen         ConnectorProvider = "ADYEN"
	ConnectorProviderGeneric       ConnectorProvider = "GENERIC"
)

type readConnectorsResponseElement struct {
	Provider    ConnectorProvider `json:"provider"`
	ConnectorID string            `json:"connectorID"`
	Name        string            `json:"name"`
	Enabled     bool              `json:"enabled"`
}
