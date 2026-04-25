package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

const legacyUpgradeRunID = "legacy-upgrade-v1"

func InitSchema(ctx context.Context, db *sql.DB, schema string) error {
	if _, err := db.ExecContext(ctx, schema); err != nil {
		if !isLegacyMissingLineageColumnError(err) {
			return fmt.Errorf("initialize schema: %w", err)
		}
		if upgradeErr := ensureSchemaUpgrades(ctx, db); upgradeErr != nil {
			return upgradeErr
		}
		if _, retryErr := db.ExecContext(ctx, schema); retryErr != nil {
			return fmt.Errorf("initialize schema: %w", retryErr)
		}
	}
	if err := ensureSchemaUpgrades(ctx, db); err != nil {
		return err
	}
	return nil
}

func isLegacyMissingLineageColumnError(err error) bool {
	msg := err.Error()
	return strings.Contains(msg, "no such column: crawl_run_id") ||
		strings.Contains(msg, "no such column: source_crawl_run_id") ||
		strings.Contains(msg, "no such column: last_attempted_at") ||
		strings.Contains(msg, "no such column: last_run_id") ||
		strings.Contains(msg, "no such column: last_error_type") ||
		strings.Contains(msg, "no such column: last_error_stage") ||
		strings.Contains(msg, "no such column: updated_at") ||
		strings.Contains(msg, "no such column: review_type") ||
		strings.Contains(msg, "no such column: platform_key")
}

func ensureSchemaUpgrades(ctx context.Context, db *sql.DB) error {
	if err := upgradeReviewFetchStateTable(ctx, db); err != nil {
		return err
	}

	type columnDef struct {
		table      string
		columnName string
		columnSQL  string
	}

	columns := []columnDef{
		{
			table:      "list_entries",
			columnName: "crawl_run_id",
			columnSQL:  "ALTER TABLE list_entries ADD COLUMN crawl_run_id TEXT NOT NULL DEFAULT ''",
		},
		{
			table:      "latest_list_entries",
			columnName: "source_crawl_run_id",
			columnSQL:  "ALTER TABLE latest_list_entries ADD COLUMN source_crawl_run_id TEXT NOT NULL DEFAULT ''",
		},
		{
			table:      "detail_fetch_state",
			columnName: "last_attempted_at",
			columnSQL:  "ALTER TABLE detail_fetch_state ADD COLUMN last_attempted_at TEXT",
		},
		{
			table:      "detail_fetch_state",
			columnName: "last_run_id",
			columnSQL:  "ALTER TABLE detail_fetch_state ADD COLUMN last_run_id TEXT",
		},
		{
			table:      "detail_fetch_state",
			columnName: "last_error_type",
			columnSQL:  "ALTER TABLE detail_fetch_state ADD COLUMN last_error_type TEXT",
		},
		{
			table:      "detail_fetch_state",
			columnName: "last_error_stage",
			columnSQL:  "ALTER TABLE detail_fetch_state ADD COLUMN last_error_stage TEXT",
		},
		{
			table:      "detail_fetch_state",
			columnName: "updated_at",
			columnSQL:  "ALTER TABLE detail_fetch_state ADD COLUMN updated_at TEXT",
		},
		{
			table:      "latest_reviews",
			columnName: "source_crawl_run_id",
			columnSQL:  "ALTER TABLE latest_reviews ADD COLUMN source_crawl_run_id TEXT NOT NULL DEFAULT ''",
		},
		{
			table:      "latest_reviews",
			columnName: "last_crawled_at",
			columnSQL:  "ALTER TABLE latest_reviews ADD COLUMN last_crawled_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP",
		},
	}

	for _, column := range columns {
		tableExists, err := hasTable(ctx, db, column.table)
		if err != nil {
			return fmt.Errorf("check table %s: %w", column.table, err)
		}
		if !tableExists {
			continue
		}

		ok, err := hasColumn(ctx, db, column.table, column.columnName)
		if err != nil {
			return fmt.Errorf("check column %s.%s: %w", column.table, column.columnName, err)
		}
		if ok {
			continue
		}
		if _, err := db.ExecContext(ctx, column.columnSQL); err != nil {
			return fmt.Errorf("apply schema upgrade for %s.%s: %w", column.table, column.columnName, err)
		}
	}

	if err := backfillLegacyRunLineage(ctx, db); err != nil {
		return err
	}
	if err := backfillDetailFetchStateMetadata(ctx, db); err != nil {
		return err
	}
	if err := backfillReviewFetchStateMetadata(ctx, db); err != nil {
		return err
	}

	return nil
}

func upgradeReviewFetchStateTable(ctx context.Context, db *sql.DB) error {
	tableExists, err := hasTable(ctx, db, "review_fetch_state")
	if err != nil {
		return fmt.Errorf("check table review_fetch_state: %w", err)
	}
	if !tableExists {
		return nil
	}

	hasReviewType, err := hasColumn(ctx, db, "review_fetch_state", "review_type")
	if err != nil {
		return fmt.Errorf("check column review_fetch_state.review_type: %w", err)
	}
	hasPlatformKey, err := hasColumn(ctx, db, "review_fetch_state", "platform_key")
	if err != nil {
		return fmt.Errorf("check column review_fetch_state.platform_key: %w", err)
	}
	hasUpdatedAt, err := hasColumn(ctx, db, "review_fetch_state", "updated_at")
	if err != nil {
		return fmt.Errorf("check column review_fetch_state.updated_at: %w", err)
	}
	if hasReviewType && hasPlatformKey && hasUpdatedAt {
		return nil
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin review_fetch_state upgrade: %w", err)
	}
	defer tx.Rollback()

	statements := []string{
		"ALTER TABLE review_fetch_state RENAME TO review_fetch_state_legacy",
		`CREATE TABLE review_fetch_state (
    work_href TEXT NOT NULL,
    review_type TEXT NOT NULL,
    platform_key TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL,
    last_attempted_at TEXT,
    last_fetched_at TEXT,
    last_run_id TEXT,
    last_error TEXT,
    last_error_type TEXT,
    last_error_stage TEXT,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (work_href, review_type, platform_key),
    FOREIGN KEY(last_run_id) REFERENCES crawl_runs(run_id),
    FOREIGN KEY(work_href) REFERENCES works(href)
)`,
		`INSERT INTO review_fetch_state (
    work_href,
    review_type,
    platform_key,
    status,
    last_attempted_at,
    last_fetched_at,
    last_run_id,
    last_error,
    last_error_type,
    last_error_stage,
    updated_at
)
SELECT
    work_href,
    'all',
    '',
    status,
    NULL,
    last_fetched_at,
    NULL,
    last_error,
    NULL,
    NULL,
    COALESCE(last_fetched_at, CURRENT_TIMESTAMP)
FROM review_fetch_state_legacy`,
		"DROP TABLE review_fetch_state_legacy",
		"CREATE INDEX IF NOT EXISTS idx_review_fetch_state_status ON review_fetch_state(status)",
		"CREATE INDEX IF NOT EXISTS idx_review_fetch_state_last_attempted_at ON review_fetch_state(last_attempted_at)",
		"CREATE INDEX IF NOT EXISTS idx_review_fetch_state_last_run_id ON review_fetch_state(last_run_id)",
		"CREATE INDEX IF NOT EXISTS idx_review_fetch_state_scope ON review_fetch_state(work_href, review_type, platform_key)",
	}

	for _, stmt := range statements {
		if _, err := tx.ExecContext(ctx, stmt); err != nil {
			return fmt.Errorf("upgrade review_fetch_state table: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit review_fetch_state upgrade: %w", err)
	}
	return nil
}

func backfillDetailFetchStateMetadata(ctx context.Context, db *sql.DB) error {
	tableExists, err := hasTable(ctx, db, "detail_fetch_state")
	if err != nil {
		return fmt.Errorf("check table detail_fetch_state: %w", err)
	}
	if !tableExists {
		return nil
	}

	if _, err := db.ExecContext(ctx, `
UPDATE detail_fetch_state
SET updated_at = COALESCE(updated_at, last_attempted_at, last_fetched_at, CURRENT_TIMESTAMP)
WHERE updated_at IS NULL OR updated_at = ''
`); err != nil {
		return fmt.Errorf("backfill detail_fetch_state updated_at: %w", err)
	}

	if _, err := db.ExecContext(ctx, `
UPDATE detail_fetch_state
SET last_attempted_at = COALESCE(last_attempted_at, updated_at)
WHERE status = 'running' AND (last_attempted_at IS NULL OR last_attempted_at = '')
`); err != nil {
		return fmt.Errorf("backfill detail_fetch_state last_attempted_at: %w", err)
	}

	return nil
}

func backfillReviewFetchStateMetadata(ctx context.Context, db *sql.DB) error {
	tableExists, err := hasTable(ctx, db, "review_fetch_state")
	if err != nil {
		return fmt.Errorf("check table review_fetch_state: %w", err)
	}
	if !tableExists {
		return nil
	}

	if _, err := db.ExecContext(ctx, `
UPDATE review_fetch_state
SET updated_at = COALESCE(updated_at, last_attempted_at, last_fetched_at, CURRENT_TIMESTAMP)
WHERE updated_at IS NULL OR updated_at = ''
`); err != nil {
		return fmt.Errorf("backfill review_fetch_state updated_at: %w", err)
	}

	if _, err := db.ExecContext(ctx, `
UPDATE review_fetch_state
SET last_attempted_at = COALESCE(last_attempted_at, updated_at)
WHERE status = 'running' AND (last_attempted_at IS NULL OR last_attempted_at = '')
`); err != nil {
		return fmt.Errorf("backfill review_fetch_state last_attempted_at: %w", err)
	}

	if _, err := db.ExecContext(ctx, `
UPDATE review_fetch_state
SET review_type = 'all'
WHERE review_type IS NULL OR review_type = ''
`); err != nil {
		return fmt.Errorf("backfill review_fetch_state review_type: %w", err)
	}

	if _, err := db.ExecContext(ctx, `
UPDATE review_fetch_state
SET platform_key = ''
WHERE platform_key IS NULL
`); err != nil {
		return fmt.Errorf("backfill review_fetch_state platform_key: %w", err)
	}

	return nil
}

func backfillLegacyRunLineage(ctx context.Context, db *sql.DB) error {
	hasLegacyRows, err := hasLegacyLineageRows(ctx, db)
	if err != nil {
		return err
	}
	if !hasLegacyRows {
		return nil
	}

	now := time.Now().UTC().Format(time.RFC3339)
	if _, err := db.ExecContext(ctx, `
INSERT OR IGNORE INTO crawl_runs (
    run_id,
    source,
    task_name,
    category,
    metric,
    filter_key,
    started_at,
    finished_at,
    status,
    error_message
) VALUES (?, 'schema upgrade', 'legacy-upgrade', 'legacy', 'legacy', 'legacy', ?, ?, 'completed', NULL)
`, legacyUpgradeRunID, now, now); err != nil {
		return fmt.Errorf("create legacy crawl run: %w", err)
	}

	if _, err := db.ExecContext(ctx, `UPDATE list_entries SET crawl_run_id = ? WHERE crawl_run_id = ''`, legacyUpgradeRunID); err != nil {
		return fmt.Errorf("backfill list_entries crawl_run_id: %w", err)
	}
	if _, err := db.ExecContext(ctx, `UPDATE latest_list_entries SET source_crawl_run_id = ? WHERE source_crawl_run_id = ''`, legacyUpgradeRunID); err != nil {
		return fmt.Errorf("backfill latest_list_entries source_crawl_run_id: %w", err)
	}
	hasLatestReviews, err := hasTable(ctx, db, "latest_reviews")
	if err != nil {
		return fmt.Errorf("check table latest_reviews: %w", err)
	}
	if hasLatestReviews {
		if _, err := db.ExecContext(ctx, `UPDATE latest_reviews SET source_crawl_run_id = ? WHERE source_crawl_run_id = ''`, legacyUpgradeRunID); err != nil {
			return fmt.Errorf("backfill latest_reviews source_crawl_run_id: %w", err)
		}
	}

	return nil
}

func hasLegacyLineageRows(ctx context.Context, db *sql.DB) (bool, error) {
	var count int
	hasLatestReviews, err := hasTable(ctx, db, "latest_reviews")
	if err != nil {
		return false, fmt.Errorf("check table latest_reviews: %w", err)
	}
	latestReviewsExpr := "0"
	if hasLatestReviews {
		latestReviewsExpr = "(SELECT COUNT(*) FROM latest_reviews WHERE source_crawl_run_id = '')"
	}
	query := fmt.Sprintf(`
SELECT
    (SELECT COUNT(*) FROM list_entries WHERE crawl_run_id = '') +
    (SELECT COUNT(*) FROM latest_list_entries WHERE source_crawl_run_id = '') +
    %s
`, latestReviewsExpr)
	if err := db.QueryRowContext(ctx, query).Scan(&count); err != nil {
		return false, fmt.Errorf("check legacy run lineage rows: %w", err)
	}
	return count > 0, nil
}

func hasColumn(ctx context.Context, db *sql.DB, table string, column string) (bool, error) {
	rows, err := db.QueryContext(ctx, fmt.Sprintf("PRAGMA table_info(%s)", table))
	if err != nil {
		return false, err
	}
	defer rows.Close()

	var (
		cid      int
		name     string
		dataType string
		notNull  int
		defaultV sql.NullString
		pk       int
	)

	for rows.Next() {
		if err := rows.Scan(&cid, &name, &dataType, &notNull, &defaultV, &pk); err != nil {
			return false, err
		}
		if name == column {
			return true, nil
		}
	}

	if err := rows.Err(); err != nil {
		return false, err
	}
	return false, nil
}

func hasTable(ctx context.Context, db *sql.DB, table string) (bool, error) {
	var count int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name = ?`, table).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}
