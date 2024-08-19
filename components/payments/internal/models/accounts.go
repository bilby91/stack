package models

import (
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gibson042/canonicaljson-go"
)

// Internal struct used by the plugins
type PSPAccount struct {
	// PSP reference of the account. Should be unique.
	Reference string

	// Account's creation date
	CreatedAt time.Time

	// Optional, human readable name of the account (if existing)
	Name *string
	// Optional, if provided the default asset of the account
	// in minor currencies unit.
	DefaultAsset *string

	// Additional metadata
	Metadata map[string]string

	// PSP response in raw
	Raw json.RawMessage
}

type Account struct {
	// Unique Account ID generated from account information
	ID AccountID `json:"id"`
	// Related Connector ID
	ConnectorID ConnectorID `json:"connectorID"`

	// PSP reference of the account. Should be unique.
	Reference string `json:"reference"`

	// Account's creation date
	CreatedAt time.Time `json:"createdAt"`

	// Type of account: INTERNAL, EXTERNAL...
	Type AccountType `json:"type"`

	// Optional, human readable name of the account (if existing)
	Name *string `json:"name"`
	// Optional, if provided the default asset of the account
	// in minor currencies unit.
	DefaultAsset *string `json:"defaultAsset"`

	// Additional metadata
	Metadata map[string]string `json:"metadata"`

	// PSP response in raw
	Raw json.RawMessage `json:"raw"`
}

type AccountID struct {
	Reference   string
	ConnectorID ConnectorID
}

func (aid *AccountID) MarshalJSON() ([]byte, error) {
	return []byte(aid.String()), nil
}

func (aid *AccountID) String() string {
	if aid == nil || aid.Reference == "" {
		return ""
	}

	data, err := canonicaljson.Marshal(aid)
	if err != nil {
		panic(err)
	}

	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(data)
}

func AccountIDFromString(value string) (AccountID, error) {
	ret := AccountID{}

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

func MustAccountIDFromString(value string) AccountID {
	id, err := AccountIDFromString(value)
	if err != nil {
		panic(err)
	}
	return id
}

func (aid AccountID) Value() (driver.Value, error) {
	return aid.String(), nil
}

func (aid *AccountID) Scan(value interface{}) error {
	if value == nil {
		return errors.New("account id is nil")
	}

	if s, err := driver.String.ConvertValue(value); err == nil {

		if v, ok := s.(string); ok {

			id, err := AccountIDFromString(v)
			if err != nil {
				return fmt.Errorf("failed to parse account id %s: %v", v, err)
			}

			*aid = id
			return nil
		}
	}

	return fmt.Errorf("failed to scan account id: %v", value)
}

type AccountType string

const (
	ACCOUNT_TYPE_UNKNOWN AccountType = "UNKNOWN"
	// Internal accounts refers to user's digital e-wallets. It serves as a
	// secure storage for funds within the payments provider environment.
	ACCOUNT_TYPE_INTERNAL AccountType = "INTERNAL"
	// External accounts represents actual bank accounts of the user.
	ACCOUNT_TYPE_EXTERNAL AccountType = "EXTERNAL"
)

func FromPSPAccount(from PSPAccount, accountType AccountType, connectorID ConnectorID) Account {
	return Account{
		ID: AccountID{
			Reference:   from.Reference,
			ConnectorID: connectorID,
		},
		ConnectorID:  connectorID,
		Reference:    from.Reference,
		CreatedAt:    from.CreatedAt,
		Type:         accountType,
		Name:         from.Name,
		DefaultAsset: from.DefaultAsset,
		Metadata:     from.Metadata,
		Raw:          from.Raw,
	}
}

func FromPSPAccounts(from []PSPAccount, accountType AccountType, connectorID ConnectorID) []Account {
	accounts := make([]Account, 0, len(from))
	for _, a := range from {
		accounts = append(accounts, FromPSPAccount(a, accountType, connectorID))
	}
	return accounts
}
