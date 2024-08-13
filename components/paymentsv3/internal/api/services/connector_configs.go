package services

import "github.com/formancehq/paymentsv3/internal/connectors/plugins"

func (s *Service) GetConnectorConfigs() plugins.Configs {
	return plugins.GetConfigs()
}
