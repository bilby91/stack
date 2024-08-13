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

type Connector struct {
	// Unique ID of the connector
	ID ConnectorID `json:"id"`
	// Name given by the user to the connector
	Name string `json:"name"`
	// Creation date
	CreatedAt time.Time `json:"createdAt"`
	// Provider type
	Provider string `json:"provider"`

	// Config given by the user. It will be encrypted when stored
	Config json.RawMessage `json:"config"`
}

type ConnectorID struct {
	Reference string
	Provider  string
}

func (cid *ConnectorID) MarshalJSON() ([]byte, error) {
	return []byte(cid.String()), nil
}

func (cid *ConnectorID) UnmarshalJSON(data []byte) error {
	id, err := ConnectorIDFromString(string(data))
	if err != nil {
		return err
	}
	*cid = id
	return nil
}

func (cid *ConnectorID) String() string {
	if cid == nil || cid.Reference == "" {
		return ""
	}

	data, err := canonicaljson.Marshal(cid)
	if err != nil {
		panic(err)
	}

	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(data)
}

func ConnectorIDFromString(value string) (ConnectorID, error) {
	data, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(value)
	if err != nil {
		return ConnectorID{}, err
	}
	ret := ConnectorID{}
	err = canonicaljson.Unmarshal(data, &ret)
	if err != nil {
		return ConnectorID{}, err
	}

	return ret, nil
}

func MustConnectorIDFromString(value string) ConnectorID {
	id, err := ConnectorIDFromString(value)
	if err != nil {
		panic(err)
	}
	return id
}

func (cid ConnectorID) Value() (driver.Value, error) {
	return cid.String(), nil
}

func (cid *ConnectorID) Scan(value interface{}) error {
	if value == nil {
		return errors.New("connector id is nil")
	}

	if s, err := driver.String.ConvertValue(value); err == nil {

		if v, ok := s.(string); ok {

			id, err := ConnectorIDFromString(v)
			if err != nil {
				return fmt.Errorf("failed to parse connector id %s: %v", v, err)
			}

			*cid = id
			return nil
		}
	}

	return fmt.Errorf("failed to scan connector id: %v", value)
}
