package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

func newAllInOneServer() *cobra.Command {
	return &cobra.Command{
		Use:          "serve",
		Aliases:      []string{"server"},
		Short:        "Launch api server",
		SilenceUsage: true,
		RunE:         runServer(),
	}
}

func runServer() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		return nil
	}
}

func serverOptions() fx.Option {
	ret := make([]fx.Option, 0)

	return fx.Options(ret...)
}
