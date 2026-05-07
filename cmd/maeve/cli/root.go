package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "dev"
var configPath string
var storePath string

func Execute() {
	if err := NewRootCommand().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "maeve",
		Short:         "Memory and evolution engine for coding agents",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.PersistentFlags().StringVar(&configPath, "config", "", "config file path")
	cmd.PersistentFlags().StringVar(&storePath, "store", "", "SQLite store path")
	cmd.AddCommand(newCompressCommand())
	cmd.AddCommand(newInitCommand())
	cmd.AddCommand(newIngestCommand())
	cmd.AddCommand(newPruneCommand())
	cmd.AddCommand(newSnapshotCommand())
	cmd.AddCommand(newStatusCommand())
	cmd.AddCommand(newVersionCommand())
	return cmd
}

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print maeve version",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(cmd.OutOrStdout(), version)
			return nil
		},
	}
}
