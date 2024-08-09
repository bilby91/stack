package models

import "time"

type Config interface {
	PageSize() int
	PollingIntervalDuration() time.Duration
}
