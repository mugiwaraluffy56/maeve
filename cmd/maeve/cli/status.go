package cli

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func newStatusCommand() *cobra.Command {
	var sessionID string
	var output string

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Print session stats",
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
			status, err := app.store.Status(ctx, sessionID)
			if err != nil {
				return err
			}
			if output == "json" {
				return json.NewEncoder(cmd.OutOrStdout()).Encode(status)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "session: %s\nobjects: %d\ntokens: %d\nscore: %.2f avg %.2f max\n",
				status.SessionID,
				status.ObjectCount,
				status.TokenCount,
				status.AvgScore,
				status.MaxScore,
			)
			return nil
		},
	}
	cmd.Flags().StringVar(&sessionID, "session", "", "session id")
	cmd.Flags().StringVarP(&output, "output", "o", "text", "output format: text or json")
	return cmd
}
