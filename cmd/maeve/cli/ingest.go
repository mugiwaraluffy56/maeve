package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newIngestCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ingest",
		Short: "Ingest context into the current session",
	}
	cmd.AddCommand(newIngestFileCommand())
	cmd.AddCommand(newIngestDiffCommand())
	return cmd
}

func newIngestFileCommand() *cobra.Command {
	var sessionID string
	cmd := &cobra.Command{
		Use:   "file <path>",
		Short: "Ingest a file snapshot",
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
			obj, err := app.ingest.File(ctx, sessionID, args[0])
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "object %s tokens=%d score=%.2f\n", obj.ID, obj.TokenCount, obj.Importance)
			return nil
		},
	}
	cmd.Flags().StringVar(&sessionID, "session", "", "session id")
	return cmd
}

func newIngestDiffCommand() *cobra.Command {
	var sessionID string
	cmd := &cobra.Command{
		Use:   "diff [path]",
		Short: "Ingest a git diff from a file or stdin",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			app, err := newApp(ctx, configPath, storePath)
			if err != nil {
				return err
			}
			defer app.close()

			var content []byte
			if len(args) == 0 {
				content, err = os.ReadFile("/dev/stdin")
			} else {
				content, err = os.ReadFile(args[0])
			}
			if err != nil {
				return err
			}

			sessionID, err := app.resolveSessionID(ctx, sessionID)
			if err != nil {
				return err
			}
			obj, err := app.ingest.Diff(ctx, sessionID, content)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "object %s tokens=%d score=%.2f\n", obj.ID, obj.TokenCount, obj.Importance)
			return nil
		},
	}
	cmd.Flags().StringVar(&sessionID, "session", "", "session id")
	return cmd
}
