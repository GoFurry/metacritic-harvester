package cli

import "github.com/spf13/cobra"

func newReviewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "review",
		Short: "Query, export, and compare review data",
	}

	cmd.AddCommand(newReviewQueryCommand())
	cmd.AddCommand(newReviewExportCommand())
	cmd.AddCommand(newReviewCompareCommand())
	return cmd
}
