package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func newSnapshotCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Manage compressed context snapshots",
	}
	cmd.AddCommand(newSnapshotSaveCommand())
	cmd.AddCommand(newSnapshotListCommand())
	return cmd
}

func newSnapshotSaveCommand() *cobra.Command {
	var sessionID string
	var budget int
	cmd := &cobra.Command{
		Use:   "save <name>",
		Short: "Save a named compressed snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			app, err := newApp(ctx, configPath, storePath)
			if err != nil {
				return err
			}
			defer app.close()

			sessionID, err := app.resolveSessionID(ctx, sessionID)
			if err != nil {
				return err
			}
			if budget == 0 {
				budget = app.cfg.Compression.DefaultBudgetTokens
			}
			payload, report, err := app.compress.Compress(ctx, sessionID, budget)
			if err != nil {
				return err
			}
			snap, err := app.snapshots.Save(ctx, sessionID, args[0], payload, report.OriginalTokens, report.CompressedTokens)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "snapshot %s\n", snap.ID)
			return nil
		},
	}
	cmd.Flags().StringVar(&sessionID, "session", "", "session id")
	cmd.Flags().IntVar(&budget, "budget", 0, "token budget")
	return cmd
}

func newSnapshotListCommand() *cobra.Command {
	var sessionID string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List snapshots",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			app, err := newApp(ctx, configPath, storePath)
			if err != nil {
				return err
			}
			defer app.close()

			sessionID, err := app.resolveSessionID(ctx, sessionID)
			if err != nil {
				return err
			}
			snapshots, err := app.snapshots.List(ctx, sessionID)
			if err != nil {
				return err
			}
			for _, snap := range snapshots {
				fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\t%d/%d tokens\n",
					snap.ID,
					snap.Name,
					snap.CompressedTokens,
					snap.OriginalTokens,
				)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&sessionID, "session", "", "session id")
	return cmd
}
