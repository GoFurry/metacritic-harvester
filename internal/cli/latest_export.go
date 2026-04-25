package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/GoFurry/metacritic-harvester/internal/storage"
)

func newLatestExportCommand() *cobra.Command {
	var (
		dbPath     string
		category   string
		metric     string
		filterKey  string
		format     string
		output     string
		checkpoint bool
	)

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export current latest_list_entries rows",
		RunE: func(cmd *cobra.Command, _ []string) (err error) {
			if err := validateOptionalCategoryMetric(category, metric); err != nil {
				return err
			}
			if format != "csv" && format != "json" {
				return fmt.Errorf("format must be csv or json")
			}
			if output == "" {
				return fmt.Errorf("output must not be empty")
			}

			repo, closeFn, err := openRepository(cmd.Context(), dbPath, checkpoint)
			if err != nil {
				return err
			}
			defer func() {
				err = finishReadRepository(err, closeFn)
			}()

			entries, err := repo.ListLatestEntries(cmd.Context(), storage.ListLatestEntriesFilter{
				Category:  category,
				Metric:    metric,
				FilterKey: filterKey,
				Limit:     -1,
			})
			if err != nil {
				return err
			}

			file, err := createOutputFile(output)
			if err != nil {
				return err
			}
			defer file.Close()

			mapped := mapLatestEntries(entries)
			if format == "json" {
				return writeJSON(file, mapped)
			}

			rows := make([][]string, 0, len(mapped))
			for _, entry := range mapped {
				rows = append(rows, []string{
					entry.WorkHref,
					entry.Category,
					entry.Metric,
					entry.FilterKey,
					fmt.Sprintf("%d", entry.PageNo),
					fmt.Sprintf("%d", entry.RankNo),
					entry.Metascore,
					entry.UserScore,
					entry.LastCrawledAt,
					entry.SourceCrawlRunID,
				})
			}

			return writeCSV(file, []string{
				"work_href",
				"category",
				"metric",
				"filter_key",
				"page_no",
				"rank_no",
				"metascore",
				"user_score",
				"last_crawled_at",
				"source_crawl_run_id",
			}, rows)
		},
	}

	cmd.Flags().StringVar(&dbPath, "db", "output/metacritic.db", "SQLite database path")
	cmd.Flags().StringVar(&category, "category", "", "Optional category filter: game|movie|tv")
	cmd.Flags().StringVar(&metric, "metric", "", "Optional metric filter: metascore|userscore|newest")
	cmd.Flags().StringVar(&filterKey, "filter-key", "", "Optional normalized filter key")
	cmd.Flags().StringVar(&format, "format", "csv", "Export format: csv|json")
	cmd.Flags().StringVar(&output, "output", "", "Output file path")
	addCheckpointFlag(cmd, &checkpoint)
	_ = cmd.MarkFlagRequired("output")

	return cmd
}
