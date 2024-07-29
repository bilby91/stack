package models

import (
	"context"
	"encoding/json"
)

type Plugin interface {
	Install(context.Context, InstallRequest) (InstallResponse, error)

	FetchAccounts(context.Context, FetchAccountsRequest) (FetchAccountsResponse, error)
	FetchPayments(context.Context, FetchPaymentsRequest) (FetchPaymentsResponse, error)
	FetchOthers(context.Context, FetchOthersRequest) (FetchOthersResponse, error)
}

type InstallRequest struct {
	Config json.RawMessage
}

type InstallResponse struct {
	Capabilities []Capability
	Workflow     Workflow
}

type FetchAccountsRequest struct {
	FromPayload json.RawMessage
	State       json.RawMessage
}

type FetchAccountsResponse struct {
	Accounts []Account
	NewState json.RawMessage
}

type FetchPaymentsRequest struct {
	FromPayload json.RawMessage
	State       json.RawMessage
}

type FetchPaymentsResponse struct {
	Payments []Payment
	NewState json.RawMessage
}

type FetchOthersRequest struct {
	Name  string
	State json.RawMessage
}

type FetchOthersResponse struct {
	Payload  json.RawMessage
	NewState json.RawMessage
}
