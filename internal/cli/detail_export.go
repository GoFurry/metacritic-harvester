package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/GoFurry/metacritic-harvester/internal/config"
	"github.com/GoFurry/metacritic-harvester/internal/domain"
	"github.com/GoFurry/metacritic-harvester/internal/storage"
)

func newDetailExportCommand() *cobra.Command {
	var (
		dbPath     string
		category   string
		workHref   string
		format     string
		output     string
		checkpoint bool
	)

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export current work_details rows",
		RunE: func(cmd *cobra.Command, _ []string) (err error) {
			if err := validateOptionalCategoryMetric(category, ""); err != nil {
				return err
			}
			workHref = domain.NormalizeWorkHref(workHref, config.DefaultBaseURL)
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

			rows, err := repo.ListWorkDetails(cmd.Context(), storage.ListWorkDetailsFilter{
				Category: category,
				WorkHref: workHref,
				Limit:    -1,
			})
			if err != nil {
				return err
			}

			mapped, err := mapWorkDetails(rows)
			if err != nil {
				return err
			}

			file, err := createOutputFile(output)
			if err != nil {
				return err
			}
			defer file.Close()

			if format == "json" {
				return writeJSON(file, mapped)
			}

			csvRows := make([][]string, 0, len(mapped))
			for _, row := range mapped {
				csvRows = append(csvRows, []string{
					row.WorkHref,
					row.Category,
					row.Title,
					row.ReleaseDate,
					row.Metascore,
					row.UserScore,
					row.Rating,
					row.Duration,
					row.Tagline,
					row.LastFetchedAt,
					row.RawDetailsJSON,
				})
			}
			return writeCSV(file, []string{
				"work_href",
				"category",
				"title",
				"release_date",
				"metascore",
				"user_score",
				"rating",
				"duration",
				"tagline",
				"last_fetched_at",
				"details_json",
			}, csvRows)
		},
	}

	cmd.Flags().StringVar(&dbPath, "db", "output/metacritic.db", "SQLite database path")
	cmd.Flags().StringVar(&category, "category", "", "Optional category filter: game|movie|tv")
	cmd.Flags().StringVar(&workHref, "work-href", "", "Optional work href filter")
	cmd.Flags().StringVar(&format, "format", "csv", "Export format: csv|json")
	cmd.Flags().StringVar(&output, "output", "", "Output file path")
	addCheckpointFlag(cmd, &checkpoint)
	_ = cmd.MarkFlagRequired("output")

	return cmd
}
