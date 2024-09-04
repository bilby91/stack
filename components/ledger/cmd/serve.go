package cmd

import (
	ledgercontroller "github.com/formancehq/ledger/internal/controller/ledger"
	"github.com/formancehq/ledger/internal/controller/ledger/writer"
	systemcontroller "github.com/formancehq/ledger/internal/controller/system"
	"github.com/formancehq/ledger/internal/storage"
	"github.com/formancehq/stack/libs/go-libs/auth"
	"github.com/formancehq/stack/libs/go-libs/aws/iam"
	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlpmetrics"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlptraces"
	"github.com/formancehq/stack/libs/go-libs/publish"

	"github.com/formancehq/ledger/internal/api"

	systemstore "github.com/formancehq/ledger/internal/storage/system"
	"github.com/formancehq/stack/libs/go-libs/ballast"
	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

const (
	BindFlag                   = "bind"
	BallastSizeInBytesFlag     = "ballast-size"
	NumscriptCacheMaxCountFlag = "numscript-cache-max-count"
	AutoUpgradeFlag            = "auto-upgrade"
)

func NewServe() *cobra.Command {
	cmd := &cobra.Command{
		Use: "serve",
		RunE: func(cmd *cobra.Command, args []string) error {
			serveConfiguration, err := discoverServeConfiguration(cmd)
			if err != nil {
				return err
			}

			connectionOptions, err := bunconnect.ConnectionOptionsFromFlags(cmd)
			if err != nil {
				return err
			}

			return service.New(cmd.OutOrStdout(),
				fx.NopLogger,
				publish.FXModuleFromFlags(cmd, service.IsDebug(cmd)),
				otlptraces.FXModuleFromFlags(cmd),
				otlpmetrics.FXModuleFromFlags(cmd),
				auth.FXModuleFromFlags(cmd),
				bunconnect.Module(*connectionOptions, service.IsDebug(cmd)),
				storage.NewFXModule(serveConfiguration.autoUpgrade),
				systemcontroller.NewFXModule(),
				ledgercontroller.NewFXModule(ledgercontroller.ModuleConfiguration{
					NSCacheConfiguration: writer.CacheConfiguration{
						MaxCount: serveConfiguration.numscriptCacheMaxCount,
					},
				}),
				ballast.Module(serveConfiguration.ballastSize),
				api.Module(api.Config{
					Version: Version,
					Debug:   service.IsDebug(cmd),
					Bind:    serveConfiguration.bind,
				}),
			).Run(cmd)
		},
	}
	cmd.Flags().Uint(BallastSizeInBytesFlag, 0, "Ballast size in bytes, default to 0")
	cmd.Flags().Int(NumscriptCacheMaxCountFlag, 1024, "Numscript cache max count")
	cmd.Flags().Bool(AutoUpgradeFlag, false, "Automatically upgrade all schemas")
	cmd.Flags().String(BindFlag, "0.0.0.0:3068", "API bind address")

	service.AddFlags(cmd.Flags())
	bunconnect.AddFlags(cmd.Flags())
	otlpmetrics.AddFlags(cmd.Flags())
	otlptraces.AddFlags(cmd.Flags())
	auth.AddFlags(cmd.Flags())
	publish.AddFlags(ServiceName, cmd.Flags(), func(cd *publish.ConfigDefault) {
		cd.PublisherCircuitBreakerSchema = systemstore.Schema
	})
	iam.AddFlags(cmd.Flags())

	return cmd
}

type serveConfiguration struct {
	ballastSize            uint
	numscriptCacheMaxCount uint
	autoUpgrade            bool
	bind                   string
}

func discoverServeConfiguration(cmd *cobra.Command) (*serveConfiguration, error) {
	ret := &serveConfiguration{}
	ret.ballastSize, _ = cmd.Flags().GetUint(BallastSizeInBytesFlag)
	ret.numscriptCacheMaxCount, _ = cmd.Flags().GetUint(NumscriptCacheMaxCountFlag)
	ret.autoUpgrade, _ = cmd.Flags().GetBool(AutoUpgradeFlag)
	ret.bind, _ = cmd.Flags().GetString(BindFlag)

	return ret, nil
}
