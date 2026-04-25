package cli

import "github.com/spf13/cobra"

func newDetailCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "detail",
		Short: "Query and compare current and historical detail data",
	}

	cmd.AddCommand(newDetailQueryCommand())
	cmd.AddCommand(newDetailExportCommand())
	cmd.AddCommand(newDetailCompareCommand())
	return cmd
}
