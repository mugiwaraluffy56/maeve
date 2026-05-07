package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func newInitCommand() *cobra.Command {
	var name string
	var branch string

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a maeve session for this repo",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			app, err := newApp(ctx, configPath, storePath)
			if err != nil {
				return err
			}
			defer app.close()

			sess, err := app.session.Init(ctx, name, app.repoPath, branch)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "session %s\n", sess.ID)
			return nil
		},
	}
	cmd.Flags().StringVar(&name, "name", "default", "session name")
	cmd.Flags().StringVar(&branch, "branch", "", "branch name")
	return cmd
}
