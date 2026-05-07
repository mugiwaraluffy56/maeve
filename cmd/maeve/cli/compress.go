package cli

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func newCompressCommand() *cobra.Command {
	var sessionID string
	var budget int
	var report bool

	cmd := &cobra.Command{
		Use:   "compress",
		Short: "Compress current session context",
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
			if budget == 0 {
				budget = app.cfg.Compression.DefaultBudgetTokens
			}
			payload, summary, err := app.compress.Compress(ctx, sessionID, budget)
			if err != nil {
				return err
			}
			if report {
				return json.NewEncoder(cmd.OutOrStdout()).Encode(summary)
			}
			fmt.Fprint(cmd.OutOrStdout(), string(payload))
			return nil
		},
	}
	cmd.Flags().StringVar(&sessionID, "session", "", "session id")
	cmd.Flags().IntVar(&budget, "budget", 0, "token budget")
	cmd.Flags().BoolVar(&report, "report", false, "print compression report as JSON")
	return cmd
}
