package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/GoFurry/metacritic-harvester/internal/config"
	"github.com/GoFurry/metacritic-harvester/internal/domain"
	"github.com/GoFurry/metacritic-harvester/internal/storage"
)

func newReviewExportCommand() *cobra.Command {
	var (
		dbPath     string
		category   string
		reviewType string
		platform   string
		workHref   string
		format     string
		output     string
		checkpoint bool
	)

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export current latest_reviews rows",
		RunE: func(cmd *cobra.Command, _ []string) (err error) {
			if err := validateOptionalCategoryMetric(category, ""); err != nil {
				return err
			}
			if strings.TrimSpace(reviewType) != "" {
				if _, err := domain.ParseReviewType(reviewType); err != nil {
					return err
				}
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

			rows, err := repo.ListLatestReviews(cmd.Context(), storage.ListLatestReviewsFilter{
				Category:   category,
				ReviewType: reviewType,
				Platform:   platform,
				WorkHref:   workHref,
				Limit:      -1,
			})
			if err != nil {
				return err
			}

			mapped := mapLatestReviews(rows)
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
					row.ReviewKey,
					row.ExternalReviewID,
					row.WorkHref,
					row.Category,
					row.ReviewType,
					row.PlatformKey,
					row.ReviewURL,
					row.ReviewDate,
					row.Score,
					row.Quote,
					row.PublicationName,
					row.PublicationSlug,
					row.AuthorName,
					row.AuthorSlug,
					row.SeasonLabel,
					row.Username,
					row.UserSlug,
					row.ThumbsUp,
					row.ThumbsDown,
					row.VersionLabel,
					row.SpoilerFlag,
					row.SourcePayloadJSON,
					row.SourceCrawlRunID,
					row.LastCrawledAt,
				})
			}
			return writeCSV(file, []string{
				"review_key",
				"external_review_id",
				"work_href",
				"category",
				"review_type",
				"platform_key",
				"review_url",
				"review_date",
				"score",
				"quote",
				"publication_name",
				"publication_slug",
				"author_name",
				"author_slug",
				"season_label",
				"username",
				"user_slug",
				"thumbs_up",
				"thumbs_down",
				"version_label",
				"spoiler_flag",
				"source_payload_json",
				"source_crawl_run_id",
				"last_crawled_at",
			}, csvRows)
		},
	}

	cmd.Flags().StringVar(&dbPath, "db", "output/metacritic.db", "SQLite database path")
	cmd.Flags().StringVar(&category, "category", "", "Optional category filter: game|movie|tv")
	cmd.Flags().StringVar(&reviewType, "review-type", "", "Optional review type filter: critic|user")
	cmd.Flags().StringVar(&platform, "platform", "", "Optional platform filter")
	cmd.Flags().StringVar(&workHref, "work-href", "", "Optional work href filter")
	cmd.Flags().StringVar(&format, "format", "csv", "Export format: csv|json")
	cmd.Flags().StringVar(&output, "output", "", "Output file path")
	addCheckpointFlag(cmd, &checkpoint)
	_ = cmd.MarkFlagRequired("output")

	return cmd
}
