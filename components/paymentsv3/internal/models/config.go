package models

import (
	"errors"
	"time"
)

const (
	defaultPollingDuration = 2 * time.Minute
	defaultPageSize        = 100
)

type Config struct {
	Name            string        `json:"name"`
	PollingDuration time.Duration `json:"pollingDuration"`
	PageSize        int           `json:"pageSize"`
}

func (c Config) Validate() error {
	if c.Name == "" {
		return errors.New("name is required")
	}

	return nil
}

func DefaultConfig() Config {
	return Config{
		PollingDuration: defaultPollingDuration,
		PageSize:        defaultPageSize,
	}
}
