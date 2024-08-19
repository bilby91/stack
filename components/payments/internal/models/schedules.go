package models

import "time"

type Schedule struct {
	ID          string
	ConnectorID ConnectorID
	CreatedAt   time.Time
}
