package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

func newConnectorsServer() *cobra.Command {
	return &cobra.Command{
		Use:          "api",
		Short:        "Launch api server",
		SilenceUsage: true,
		RunE:         runConnectorsServer(),
	}
}

func runConnectorsServer() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		return nil
	}
}

func connectorsServerOptions() fx.Option {
	ret := make([]fx.Option, 0)

	return fx.Options(ret...)
}
