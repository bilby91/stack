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

type Account struct {
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

type ExpandedAccount struct {
	ID          AccountID
	ConnectorID ConnectorID

	Account
}

type AccountID struct {
	Reference   string
	ConnectorID ConnectorID
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

func AccountIDFromString(value string) (*AccountID, error) {
	data, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(value)
	if err != nil {
		return nil, err
	}
	ret := AccountID{}
	err = canonicaljson.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

func MustAccountIDFromString(value string) AccountID {
	id, err := AccountIDFromString(value)
	if err != nil {
		panic(err)
	}
	return *id
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

			*aid = *id
			return nil
		}
	}

	return fmt.Errorf("failed to scan account id: %v", value)
}
