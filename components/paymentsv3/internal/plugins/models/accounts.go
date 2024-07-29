package models

import (
	"encoding/json"
	"time"
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
