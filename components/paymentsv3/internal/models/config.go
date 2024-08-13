package models

import "time"

type Config struct {
	Name            string        `json:"name"`
	PollingDuration time.Duration `json:"pollingDuration"`
	PageSize        int           `json:"pageSize"`
}
