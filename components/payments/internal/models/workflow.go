package models

import (
	"time"
)

type Workflow struct {
	ID          string
	ConnectorID ConnectorID
	CreatedAt   time.Time
	Capability  Capability
	Metadata    map[string]string
}
