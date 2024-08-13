package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

func newAPIServer() *cobra.Command {
	return &cobra.Command{
		Use:          "api",
		Short:        "Launch api server",
		SilenceUsage: true,
		RunE:         runAPIServer(),
	}
}

func runAPIServer() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		return nil
	}
}

func apiServerOptions() fx.Option {
	ret := make([]fx.Option, 0)

	return fx.Options(ret...)
}
