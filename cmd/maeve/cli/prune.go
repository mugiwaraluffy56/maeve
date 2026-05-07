package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func newPruneCommand() *cobra.Command {
	var sessionID string
	var threshold float64

	cmd := &cobra.Command{
		Use:   "prune",
		Short: "Remove low-importance context objects",
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
			count, err := app.store.PruneContextObjects(ctx, sessionID, threshold)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "pruned %d objects\n", count)
			return nil
		},
	}
	cmd.Flags().StringVar(&sessionID, "session", "", "session id")
	cmd.Flags().Float64Var(&threshold, "threshold", 0.3, "importance threshold")
	return cmd
}
