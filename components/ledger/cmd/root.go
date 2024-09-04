package cmd

import (
	"github.com/formancehq/stack/libs/go-libs/bun/bunmigrate"
	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/uptrace/bun"

	"github.com/spf13/cobra"
)

const (
	ServiceName = "ledger"
)

var (
	Version   = "develop"
	BuildDate = "-"
	Commit    = "-"
)

func NewRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:               "ledger",
		Short:             "ledger",
		DisableAutoGenTag: true,
		Version:           Version,
	}

	serve := NewServe()
	version := NewVersion()

	buckets := NewBucket()
	buckets.AddCommand(NewBucketUpgrade())

	root.AddCommand(serve)
	root.AddCommand(buckets)
	root.AddCommand(version)
	root.AddCommand(bunmigrate.NewDefaultCommand(func(cmd *cobra.Command, args []string, db *bun.DB) error {
		return upgradeAll(cmd, args)
	}))

	root.AddCommand(NewDocCommand())

	return root
}

func Execute() {
	service.Execute(NewRootCommand())
}
