package moneycorp

import (
	"encoding/json"
	"errors"
)

const (
	defaultPageSize = 100
)

type Config struct {
	ClientID string `json:"clientID"`
	APIKey   string `json:"apiKey"`
	Endpoint string `json:"endpoint"`
	PageSize int    `json:"pageSize"`
}

func (c Config) Validate() error {
	if c.ClientID == "" {
		return errors.New("missing clientID in config")
	}

	if c.APIKey == "" {
		return errors.New("missing api key in config")
	}

	if c.Endpoint == "" {
		return errors.New("missing endpoint in config")
	}

	if c.PageSize == 0 {
		return errors.New("invalid page size in config")
	}

	return nil
}

func (c *Config) FillDefault() {
	if c.PageSize == 0 {
		c.PageSize = defaultPageSize
	}
}

func unmarshalAndValidateConfig(payload json.RawMessage) (Config, error) {
	var config Config
	if err := json.Unmarshal(payload, &config); err != nil {
		return Config{}, err
	}

	config.FillDefault()

	return config, config.Validate()
}
