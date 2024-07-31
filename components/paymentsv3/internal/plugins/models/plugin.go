package models

import (
	"context"
	"encoding/json"
)

type Plugin interface {
	Install(context.Context, InstallRequest) (InstallResponse, error)

	FetchNextAccounts(context.Context, FetchNextAccountsRequest) (FetchNextAccountsResponse, error)
	FetchNextPayments(context.Context, FetchNextPaymentsRequest) (FetchNextPaymentsResponse, error)
	FetchNextOthers(context.Context, FetchNextOthersRequest) (FetchNextOthersResponse, error)
}

type InstallRequest struct {
	Config json.RawMessage
}

type InstallResponse struct {
	Capabilities []Capability
	Workflow     Workflow
}

type FetchNextAccountsRequest struct {
	FromPayload json.RawMessage
	State       json.RawMessage
	PageSize    int
}

type FetchNextAccountsResponse struct {
	Accounts []Account
	NewState json.RawMessage
	HasMore  bool
}

type FetchNextPaymentsRequest struct {
	FromPayload json.RawMessage
	State       json.RawMessage
	PageSize    int
}

type FetchNextPaymentsResponse struct {
	Payments []Payment
	NewState json.RawMessage
	HasMore  bool
}

type FetchNextOthersRequest struct {
	Name        string
	FromPayload json.RawMessage
	State       json.RawMessage
	PageSize    int
}

type FetchNextOthersResponse struct {
	Others   []json.RawMessage
	NewState json.RawMessage
	HasMore  bool
}
