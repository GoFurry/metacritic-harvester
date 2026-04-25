-- name: CreateCrawlRun :exec
INSERT INTO crawl_runs (
    run_id,
    source,
    task_name,
    category,
    metric,
    filter_key,
    started_at,
    status,
    error_message
) VALUES (
    sqlc.arg(run_id),
    sqlc.arg(source),
    sqlc.arg(task_name),
    sqlc.arg(category),
    sqlc.arg(metric),
    sqlc.arg(filter_key),
    sqlc.arg(started_at),
    sqlc.arg(status),
    sqlc.narg(error_message)
);

-- name: CompleteCrawlRun :exec
UPDATE crawl_runs
SET finished_at = sqlc.arg(finished_at),
    status = sqlc.arg(status),
    error_message = NULL
WHERE run_id = sqlc.arg(run_id);

-- name: FailCrawlRun :exec
UPDATE crawl_runs
SET finished_at = sqlc.arg(finished_at),
    status = sqlc.arg(status),
    error_message = sqlc.arg(error_message)
WHERE run_id = sqlc.arg(run_id);

-- name: GetCrawlRun :one
SELECT run_id, source, task_name, category, metric, filter_key, started_at, finished_at, status, error_message
FROM crawl_runs
WHERE run_id = ? LIMIT 1;

-- name: ListCrawlRuns :many
SELECT run_id, source, task_name, category, metric, filter_key, started_at, finished_at, status, error_message
FROM crawl_runs
ORDER BY started_at DESC
LIMIT CASE WHEN sqlc.arg(limit_rows) <= 0 THEN -1 ELSE sqlc.arg(limit_rows) END;

-- name: UpsertWork :exec
INSERT INTO works (
    href,
    name,
    image_url,
    release_date,
    category,
    created_at,
    updated_at
) VALUES (
    sqlc.arg(href),
    sqlc.arg(name),
    sqlc.narg(image_url),
    sqlc.narg(release_date),
    sqlc.arg(category),
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
)
ON CONFLICT(href) DO UPDATE SET
    name = excluded.name,
    image_url = excluded.image_url,
    release_date = excluded.release_date,
    category = excluded.category,
    updated_at = CURRENT_TIMESTAMP;

-- name: InsertListEntry :exec
INSERT INTO list_entries (
    crawl_run_id,
    work_href,
    category,
    metric,
    page_no,
    rank_no,
    metascore,
    user_score,
    filter_key,
    crawled_at
) VALUES (
    sqlc.arg(crawl_run_id),
    sqlc.arg(work_href),
    sqlc.arg(category),
    sqlc.arg(metric),
    sqlc.arg(page_no),
    sqlc.arg(rank_no),
    sqlc.narg(metascore),
    sqlc.narg(user_score),
    sqlc.arg(filter_key),
    sqlc.arg(crawled_at)
);

-- name: UpsertLatestListEntry :exec
INSERT INTO latest_list_entries (
    work_href,
    category,
    metric,
    filter_key,
    page_no,
    rank_no,
    metascore,
    user_score,
    source_crawl_run_id,
    last_crawled_at,
    created_at,
    updated_at
) VALUES (
    sqlc.arg(work_href),
    sqlc.arg(category),
    sqlc.arg(metric),
    sqlc.arg(filter_key),
    sqlc.arg(page_no),
    sqlc.arg(rank_no),
    sqlc.narg(metascore),
    sqlc.narg(user_score),
    sqlc.arg(source_crawl_run_id),
    sqlc.arg(last_crawled_at),
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
)
ON CONFLICT(work_href, category, metric, filter_key) DO UPDATE SET
    page_no = excluded.page_no,
    rank_no = excluded.rank_no,
    metascore = excluded.metascore,
    user_score = excluded.user_score,
    source_crawl_run_id = excluded.source_crawl_run_id,
    last_crawled_at = excluded.last_crawled_at,
    updated_at = CURRENT_TIMESTAMP;

-- name: GetWorkByHref :one
SELECT href, name, image_url, release_date, category, created_at, updated_at
FROM works
WHERE href = ? LIMIT 1;

-- name: GetLatestListEntry :one
SELECT work_href, category, metric, filter_key, page_no, rank_no, metascore, user_score, source_crawl_run_id, last_crawled_at, created_at, updated_at
FROM latest_list_entries
WHERE work_href = ? AND category = ? AND metric = ? AND filter_key = ?
LIMIT 1;

-- name: ListLatestEntries :many
SELECT work_href, category, metric, filter_key, page_no, rank_no, metascore, user_score, source_crawl_run_id, last_crawled_at, created_at, updated_at
FROM latest_list_entries
WHERE (sqlc.arg(category) = '' OR category = sqlc.arg(category))
  AND (sqlc.arg(metric) = '' OR metric = sqlc.arg(metric))
  AND (sqlc.arg(work_href) = '' OR RTRIM(work_href, '/') = RTRIM(sqlc.arg(work_href), '/'))
  AND (sqlc.arg(filter_key) = '' OR filter_key = sqlc.arg(filter_key))
ORDER BY category ASC, metric ASC, rank_no ASC, work_href ASC
LIMIT CASE WHEN sqlc.arg(limit_rows) <= 0 THEN -1 ELSE sqlc.arg(limit_rows) END;

-- name: ListListEntriesByRun :many
SELECT id, crawl_run_id, work_href, category, metric, page_no, rank_no, metascore, user_score, filter_key, crawled_at
FROM list_entries
WHERE crawl_run_id = sqlc.arg(crawl_run_id)
  AND (sqlc.arg(category) = '' OR category = sqlc.arg(category))
  AND (sqlc.arg(metric) = '' OR metric = sqlc.arg(metric))
  AND (sqlc.arg(work_href) = '' OR RTRIM(work_href, '/') = RTRIM(sqlc.arg(work_href), '/'))
  AND (sqlc.arg(filter_key) = '' OR filter_key = sqlc.arg(filter_key))
ORDER BY category ASC, metric ASC, filter_key ASC, rank_no ASC, work_href ASC
LIMIT CASE WHEN sqlc.arg(limit_rows) <= 0 THEN -1 ELSE sqlc.arg(limit_rows) END;

-- name: GetLatestEntryByWork :many
SELECT work_href, category, metric, filter_key, page_no, rank_no, metascore, user_score, source_crawl_run_id, last_crawled_at, created_at, updated_at
FROM latest_list_entries
WHERE work_href = sqlc.arg(work_href)
  AND (sqlc.arg(category) = '' OR category = sqlc.arg(category))
  AND (sqlc.arg(metric) = '' OR metric = sqlc.arg(metric))
ORDER BY category ASC, metric ASC, rank_no ASC;

-- name: CompareCrawlRuns :many
WITH
from_entries AS (
    SELECT work_href, category, metric, filter_key, rank_no, metascore, user_score
    FROM list_entries le
    WHERE le.crawl_run_id = sqlc.arg(from_run_id)
      AND (sqlc.arg(category) = '' OR le.category = sqlc.arg(category))
      AND (sqlc.arg(metric) = '' OR le.metric = sqlc.arg(metric))
),
to_entries AS (
    SELECT work_href, category, metric, filter_key, rank_no, metascore, user_score
    FROM list_entries le
    WHERE le.crawl_run_id = sqlc.arg(to_run_id)
      AND (sqlc.arg(category) = '' OR le.category = sqlc.arg(category))
      AND (sqlc.arg(metric) = '' OR le.metric = sqlc.arg(metric))
),
compare_rows AS (
    SELECT
        COALESCE(f.work_href, t.work_href) AS work_href,
        COALESCE(f.category, t.category) AS category,
        COALESCE(f.metric, t.metric) AS metric,
        COALESCE(f.filter_key, t.filter_key) AS filter_key,
        IFNULL(f.rank_no, 0) AS from_rank,
        t.rank_no AS to_rank,
        CASE
            WHEN f.rank_no IS NOT NULL AND t.rank_no IS NOT NULL THEN t.rank_no - f.rank_no
            ELSE NULL
        END AS rank_diff,
        f.metascore AS from_metascore,
        t.metascore AS to_metascore,
        CASE
            WHEN f.metascore IS NOT NULL AND t.metascore IS NOT NULL THEN CAST(t.metascore AS REAL) - CAST(f.metascore AS REAL)
            ELSE NULL
        END AS metascore_diff,
        f.user_score AS from_user_score,
        t.user_score AS to_user_score,
        CASE
            WHEN f.user_score IS NOT NULL AND t.user_score IS NOT NULL THEN CAST(t.user_score AS REAL) - CAST(f.user_score AS REAL)
            ELSE NULL
        END AS user_score_diff,
        CASE
            WHEN f.work_href IS NULL THEN 'added'
            WHEN t.work_href IS NULL THEN 'removed'
            WHEN IFNULL(f.rank_no, -1) != IFNULL(t.rank_no, -1)
              OR IFNULL(f.metascore, '') != IFNULL(t.metascore, '')
              OR IFNULL(f.user_score, '') != IFNULL(t.user_score, '') THEN 'changed'
            ELSE 'unchanged'
        END AS change_type
    FROM from_entries f
    LEFT JOIN to_entries t
      ON f.work_href = t.work_href
     AND f.category = t.category
     AND f.metric = t.metric
     AND f.filter_key = t.filter_key

    UNION ALL

    SELECT
        COALESCE(f.work_href, t.work_href) AS work_href,
        COALESCE(f.category, t.category) AS category,
        COALESCE(f.metric, t.metric) AS metric,
        COALESCE(f.filter_key, t.filter_key) AS filter_key,
        IFNULL(f.rank_no, 0) AS from_rank,
        t.rank_no AS to_rank,
        CASE
            WHEN f.rank_no IS NOT NULL AND t.rank_no IS NOT NULL THEN t.rank_no - f.rank_no
            ELSE NULL
        END AS rank_diff,
        f.metascore AS from_metascore,
        t.metascore AS to_metascore,
        CASE
            WHEN f.metascore IS NOT NULL AND t.metascore IS NOT NULL THEN CAST(t.metascore AS REAL) - CAST(f.metascore AS REAL)
            ELSE NULL
        END AS metascore_diff,
        f.user_score AS from_user_score,
        t.user_score AS to_user_score,
        CASE
            WHEN f.user_score IS NOT NULL AND t.user_score IS NOT NULL THEN CAST(t.user_score AS REAL) - CAST(f.user_score AS REAL)
            ELSE NULL
        END AS user_score_diff,
        CASE
            WHEN f.work_href IS NULL THEN 'added'
            WHEN t.work_href IS NULL THEN 'removed'
            WHEN IFNULL(f.rank_no, -1) != IFNULL(t.rank_no, -1)
              OR IFNULL(f.metascore, '') != IFNULL(t.metascore, '')
              OR IFNULL(f.user_score, '') != IFNULL(t.user_score, '') THEN 'changed'
            ELSE 'unchanged'
        END AS change_type
    FROM to_entries t
    LEFT JOIN from_entries f
      ON f.work_href = t.work_href
     AND f.category = t.category
     AND f.metric = t.metric
     AND f.filter_key = t.filter_key
    WHERE f.work_href IS NULL
)
SELECT work_href, category, metric, filter_key, from_rank, to_rank, rank_diff, from_metascore, to_metascore, metascore_diff, from_user_score, to_user_score, user_score_diff, change_type
FROM compare_rows
WHERE (sqlc.arg(include_unchanged) = 1 OR change_type <> 'unchanged')
ORDER BY category ASC, metric ASC, filter_key ASC, change_type ASC, work_href ASC;

-- name: CountWorks :one
SELECT COUNT(*) FROM works;

-- name: CountListEntries :one
SELECT COUNT(*) FROM list_entries;

-- name: CountLatestListEntries :one
SELECT COUNT(*) FROM latest_list_entries;

-- name: CountCrawlRuns :one
SELECT COUNT(*) FROM crawl_runs;
