package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/formancehq/payments/internal/api"
	v2 "github.com/formancehq/payments/internal/api/v2"
	v3 "github.com/formancehq/payments/internal/api/v3"
	"github.com/formancehq/payments/internal/connectors/engine"
	"github.com/formancehq/payments/internal/storage"
	sharedapi "github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/auth"
	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"
	"github.com/formancehq/stack/libs/go-libs/bun/bunmigrate"
	"github.com/formancehq/stack/libs/go-libs/health"
	"github.com/formancehq/stack/libs/go-libs/licence"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlptraces"
	"github.com/formancehq/stack/libs/go-libs/temporal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var (
	ServiceName = "payments"
	Version     = "develop"
	BuildDate   = "-"
	Commit      = "-"
)

const (
	pluginPathsFlag         = "plugin-paths"
	configEncryptionKeyFlag = "config-encryption-key"
	listenFlag              = "listen"
)

func NewRootCommand() *cobra.Command {
	viper.SetDefault("version", Version)

	root := &cobra.Command{
		Use:               "payments",
		Short:             "payments",
		DisableAutoGenTag: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return bindFlagsToViper(cmd)
		},
	}

	root.PersistentFlags().String(configEncryptionKeyFlag, "", "Config encryption key")

	version := newVersion()
	root.AddCommand(version)

	migrate := newMigrate()
	root.AddCommand(migrate)

	server := newServer()
	addAutoMigrateCommand(server)
	server.Flags().String(listenFlag, ":8080", "Listen address")
	server.Flags().StringSlice(pluginPathsFlag, nil, "Plugin paths")
	root.AddCommand(server)

	return root
}

func Execute() {
	if err := NewRootCommand().Execute(); err != nil {
		if _, err = fmt.Fprintln(os.Stderr, err); err != nil {
			panic(err)
		}

		os.Exit(1)
	}
}

func addAutoMigrateCommand(cmd *cobra.Command) {
	cmd.Flags().Bool(autoMigrateFlag, false, "Auto migrate database")
	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if viper.GetBool(autoMigrateFlag) {
			return bunmigrate.Run(cmd, args, Migrate)
		}
		return nil
	}
}

func commonOptions(cmd *cobra.Command) (fx.Option, error) {
	configEncryptionKey := viper.GetString(configEncryptionKeyFlag)
	if configEncryptionKey == "" {
		return nil, errors.New("missing config encryption key")
	}

	connectionOptions, err := bunconnect.ConnectionOptionsFromFlags(cmd.Context())
	if err != nil {
		return nil, err
	}

	pluginPaths, err := getPluginsMap(viper.GetStringSlice(pluginPathsFlag))
	if err != nil {
		return nil, err
	}

	return fx.Options(
		fx.Provide(func() *bunconnect.ConnectionOptions {
			return connectionOptions
		}),
		fx.Provide(func() sharedapi.ServiceInfo {
			return sharedapi.ServiceInfo{
				Version: Version,
			}
		}),
		otlptraces.CLITracesModule(),
		temporal.NewModule(
			engine.Tracer,
			temporal.SearchAttributes{
				SearchAttributes: engine.SearchAttributes,
			},
		),
		bunconnect.Module(*connectionOptions),
		auth.CLIAuthModule(),
		health.Module(),
		licence.CLIModule(ServiceName),
		storage.Module(*connectionOptions, configEncryptionKey),
		api.NewModule(viper.GetString(listenFlag)),
		engine.Module(pluginPaths),
		v2.NewModule(),
		v3.NewModule(),
	), nil
}

func getPluginsMap(paths []string) (map[string]string, error) {
	if len(paths) == 0 {
		return nil, errors.New("no plugin paths provided")
	}

	if len(paths)%2 != 0 {
		return nil, errors.New("invalid plugin paths provided")
	}

	plugins := make(map[string]string)
	for i := 0; i < len(paths); i += 2 {
		plugins[paths[i]] = paths[i+1]
	}

	return plugins, nil
}
