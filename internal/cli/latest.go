package cli

import "github.com/spf13/cobra"

func newLatestCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "latest",
		Short: "Query and compare latest Metacritic list data",
	}

	cmd.AddCommand(newLatestQueryCommand())
	cmd.AddCommand(newLatestExportCommand())
	cmd.AddCommand(newLatestCompareCommand())
	return cmd
}
