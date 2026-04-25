package cli

import "github.com/spf13/cobra"

func newCrawlCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "crawl",
		Short: "Run crawler commands",
	}

	cmd.AddCommand(newCrawlListCommand())
	cmd.AddCommand(newCrawlDetailCommand())
	cmd.AddCommand(newCrawlReviewsCommand())
	cmd.AddCommand(newCrawlBatchCommand())
	cmd.AddCommand(newCrawlScheduleCommand())
	return cmd
}
