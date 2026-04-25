-- name: ListDetailCandidates :many
SELECT
    w.href,
    w.name,
    w.image_url,
    w.release_date,
    w.category,
    w.created_at,
    w.updated_at,
    dfs.status AS fetch_status,
    dfs.last_attempted_at,
    dfs.last_fetched_at,
    dfs.last_run_id,
    dfs.last_error,
    dfs.last_error_type,
    dfs.last_error_stage,
    dfs.updated_at AS fetch_updated_at,
    wd.work_href IS NOT NULL AS has_detail
FROM works w
LEFT JOIN detail_fetch_state dfs ON dfs.work_href = w.href
LEFT JOIN work_details wd ON wd.work_href = w.href
WHERE (sqlc.arg(category) = '' OR w.category = sqlc.arg(category))
  AND (sqlc.arg(work_href) = '' OR RTRIM(w.href, '/') = RTRIM(sqlc.arg(work_href), '/'))
  AND (
      sqlc.arg(force_refresh) = 1
      OR dfs.status IS NULL
      OR dfs.status <> 'succeeded'
      OR wd.work_href IS NULL
  )
ORDER BY w.category ASC, w.href ASC
LIMIT CASE WHEN sqlc.arg(limit_rows) <= 0 THEN -1 ELSE sqlc.arg(limit_rows) END;

-- name: GetWorkDetail :one
SELECT
    work_href,
    category,
    title,
    summary,
    release_date,
    metascore,
    metascore_sentiment,
    metascore_review_count,
    user_score,
    user_score_sentiment,
    user_score_count,
    rating,
    duration,
    tagline,
    details_json,
    last_fetched_at,
    created_at,
    updated_at
FROM work_details
WHERE RTRIM(work_href, '/') = RTRIM(?, '/') LIMIT 1;

-- name: ListWorkDetails :many
SELECT
    work_href,
    category,
    title,
    summary,
    release_date,
    metascore,
    metascore_sentiment,
    metascore_review_count,
    user_score,
    user_score_sentiment,
    user_score_count,
    rating,
    duration,
    tagline,
    details_json,
    last_fetched_at,
    created_at,
    updated_at
FROM work_details
WHERE (sqlc.arg(category) = '' OR category = sqlc.arg(category))
  AND (sqlc.arg(work_href) = '' OR RTRIM(work_href, '/') = RTRIM(sqlc.arg(work_href), '/'))
ORDER BY category ASC, title ASC, work_href ASC
LIMIT CASE WHEN sqlc.arg(limit_rows) <= 0 THEN -1 ELSE sqlc.arg(limit_rows) END;

-- name: UpsertWorkDetail :exec
INSERT INTO work_details (
    work_href,
    category,
    title,
    summary,
    release_date,
    metascore,
    metascore_sentiment,
    metascore_review_count,
    user_score,
    user_score_sentiment,
    user_score_count,
    rating,
    duration,
    tagline,
    details_json,
    last_fetched_at,
    created_at,
    updated_at
) VALUES (
    sqlc.arg(work_href),
    sqlc.arg(category),
    sqlc.arg(title),
    sqlc.narg(summary),
    sqlc.narg(release_date),
    sqlc.narg(metascore),
    sqlc.narg(metascore_sentiment),
    sqlc.narg(metascore_review_count),
    sqlc.narg(user_score),
    sqlc.narg(user_score_sentiment),
    sqlc.narg(user_score_count),
    sqlc.narg(rating),
    sqlc.narg(duration),
    sqlc.narg(tagline),
    sqlc.arg(details_json),
    sqlc.arg(last_fetched_at),
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
)
ON CONFLICT(work_href) DO UPDATE SET
    category = excluded.category,
    title = excluded.title,
    summary = excluded.summary,
    release_date = excluded.release_date,
    metascore = excluded.metascore,
    metascore_sentiment = excluded.metascore_sentiment,
    metascore_review_count = excluded.metascore_review_count,
    user_score = excluded.user_score,
    user_score_sentiment = excluded.user_score_sentiment,
    user_score_count = excluded.user_score_count,
    rating = excluded.rating,
    duration = excluded.duration,
    tagline = excluded.tagline,
    details_json = excluded.details_json,
    last_fetched_at = excluded.last_fetched_at,
    updated_at = CURRENT_TIMESTAMP;

-- name: InsertWorkDetailSnapshot :exec
INSERT INTO work_detail_snapshots (
    work_href,
    crawl_run_id,
    category,
    title,
    summary,
    release_date,
    metascore,
    metascore_sentiment,
    metascore_review_count,
    user_score,
    user_score_sentiment,
    user_score_count,
    rating,
    duration,
    tagline,
    details_json,
    fetched_at
) VALUES (
    sqlc.arg(work_href),
    sqlc.arg(crawl_run_id),
    sqlc.arg(category),
    sqlc.arg(title),
    sqlc.narg(summary),
    sqlc.narg(release_date),
    sqlc.narg(metascore),
    sqlc.narg(metascore_sentiment),
    sqlc.narg(metascore_review_count),
    sqlc.narg(user_score),
    sqlc.narg(user_score_sentiment),
    sqlc.narg(user_score_count),
    sqlc.narg(rating),
    sqlc.narg(duration),
    sqlc.narg(tagline),
    sqlc.arg(details_json),
    sqlc.arg(fetched_at)
)
ON CONFLICT(work_href, crawl_run_id) DO NOTHING;

-- name: UpdateWorkFromDetail :exec
UPDATE works
SET name = CASE WHEN sqlc.arg(name) = '' THEN name ELSE sqlc.arg(name) END,
    release_date = CASE WHEN sqlc.arg(release_date) = '' THEN release_date ELSE sqlc.arg(release_date) END,
    category = sqlc.arg(category),
    updated_at = CURRENT_TIMESTAMP
WHERE href = sqlc.arg(href);

-- name: UpsertDetailFetchStateRunning :exec
INSERT INTO detail_fetch_state (
    work_href,
    status,
    last_attempted_at,
    last_fetched_at,
    last_run_id,
    last_error,
    last_error_type,
    last_error_stage,
    updated_at
) VALUES (
    sqlc.arg(work_href),
    'running',
    sqlc.arg(last_attempted_at),
    NULL,
    sqlc.arg(last_run_id),
    NULL,
    NULL,
    NULL,
    sqlc.arg(updated_at)
)
ON CONFLICT(work_href) DO UPDATE SET
    status = 'running',
    last_attempted_at = excluded.last_attempted_at,
    last_run_id = excluded.last_run_id,
    last_error = NULL,
    last_error_type = NULL,
    last_error_stage = NULL,
    updated_at = excluded.updated_at;

-- name: UpsertDetailFetchStateSucceeded :exec
INSERT INTO detail_fetch_state (
    work_href,
    status,
    last_attempted_at,
    last_fetched_at,
    last_run_id,
    last_error,
    last_error_type,
    last_error_stage,
    updated_at
) VALUES (
    sqlc.arg(work_href),
    'succeeded',
    sqlc.arg(last_attempted_at),
    sqlc.arg(last_fetched_at),
    sqlc.arg(last_run_id),
    NULL,
    NULL,
    NULL,
    sqlc.arg(updated_at)
)
ON CONFLICT(work_href) DO UPDATE SET
    status = 'succeeded',
    last_attempted_at = excluded.last_attempted_at,
    last_fetched_at = excluded.last_fetched_at,
    last_run_id = excluded.last_run_id,
    last_error = NULL,
    last_error_type = NULL,
    last_error_stage = NULL,
    updated_at = excluded.updated_at;

-- name: UpsertDetailFetchStateFailed :exec
INSERT INTO detail_fetch_state (
    work_href,
    status,
    last_attempted_at,
    last_fetched_at,
    last_run_id,
    last_error,
    last_error_type,
    last_error_stage,
    updated_at
) VALUES (
    sqlc.arg(work_href),
    'failed',
    sqlc.arg(last_attempted_at),
    NULL,
    sqlc.arg(last_run_id),
    sqlc.narg(last_error),
    sqlc.narg(last_error_type),
    sqlc.narg(last_error_stage),
    sqlc.arg(updated_at)
)
ON CONFLICT(work_href) DO UPDATE SET
    status = 'failed',
    last_attempted_at = excluded.last_attempted_at,
    last_run_id = excluded.last_run_id,
    last_error = excluded.last_error,
    last_error_type = excluded.last_error_type,
    last_error_stage = excluded.last_error_stage,
    updated_at = excluded.updated_at;

-- name: RecoverStaleDetailFetchStates :execrows
UPDATE detail_fetch_state
SET status = 'failed',
    last_run_id = sqlc.arg(last_run_id),
    last_error = 'recovered stale running state',
    last_error_type = 'state_recovered',
    last_error_stage = 'recovery',
    updated_at = sqlc.arg(updated_at)
WHERE status = 'running'
  AND (
      last_attempted_at IS NULL
      OR last_attempted_at = ''
      OR last_attempted_at <= sqlc.arg(stale_before)
  )
  AND work_href IN (
      SELECT href
      FROM works
      WHERE (sqlc.arg(category) = '' OR category = sqlc.arg(category))
        AND (sqlc.arg(work_href) = '' OR RTRIM(href, '/') = RTRIM(sqlc.arg(work_href), '/'))
  );

-- name: GetDetailFetchState :one
SELECT
    work_href,
    status,
    last_attempted_at,
    last_fetched_at,
    last_run_id,
    last_error,
    last_error_type,
    last_error_stage,
    updated_at
FROM detail_fetch_state
WHERE work_href = ? LIMIT 1;

-- name: CompareWorkDetailSnapshots :many
WITH
from_snapshots AS (
    SELECT
        s.work_href,
        s.category,
        s.title,
        s.release_date,
        s.metascore,
        s.user_score,
        s.rating,
        s.duration,
        s.tagline,
        s.details_json
    FROM work_detail_snapshots s
    WHERE s.crawl_run_id = sqlc.arg(from_run_id)
      AND (sqlc.arg(category) = '' OR s.category = sqlc.arg(category))
      AND (sqlc.arg(work_href) = '' OR RTRIM(s.work_href, '/') = RTRIM(sqlc.arg(work_href), '/'))
),
to_snapshots AS (
    SELECT
        s.work_href,
        s.category,
        s.title,
        s.release_date,
        s.metascore,
        s.user_score,
        s.rating,
        s.duration,
        s.tagline,
        s.details_json
    FROM work_detail_snapshots s
    WHERE s.crawl_run_id = sqlc.arg(to_run_id)
      AND (sqlc.arg(category) = '' OR s.category = sqlc.arg(category))
      AND (sqlc.arg(work_href) = '' OR RTRIM(s.work_href, '/') = RTRIM(sqlc.arg(work_href), '/'))
),
compare_rows AS (
    SELECT
        COALESCE(f.work_href, t.work_href) AS work_href,
        COALESCE(f.category, t.category) AS category,
        CASE WHEN f.work_href IS NULL THEN NULL ELSE f.title END AS from_title,
        t.title AS to_title,
        f.release_date AS from_release_date,
        t.release_date AS to_release_date,
        f.metascore AS from_metascore,
        t.metascore AS to_metascore,
        f.user_score AS from_user_score,
        t.user_score AS to_user_score,
        f.rating AS from_rating,
        t.rating AS to_rating,
        f.duration AS from_duration,
        t.duration AS to_duration,
        f.tagline AS from_tagline,
        t.tagline AS to_tagline,
        CASE WHEN f.work_href IS NULL THEN NULL ELSE f.details_json END AS from_details_json,
        t.details_json AS to_details_json,
        CASE
            WHEN f.work_href IS NOT NULL
             AND t.work_href IS NOT NULL
             AND IFNULL(f.details_json, '{}') <> IFNULL(t.details_json, '{}') THEN 1
            ELSE 0
        END AS details_json_changed,
        CASE
            WHEN f.work_href IS NULL THEN 'added'
            WHEN t.work_href IS NULL THEN 'removed'
            WHEN IFNULL(f.title, '') <> IFNULL(t.title, '')
              OR IFNULL(f.release_date, '') <> IFNULL(t.release_date, '')
              OR IFNULL(f.metascore, '') <> IFNULL(t.metascore, '')
              OR IFNULL(f.user_score, '') <> IFNULL(t.user_score, '')
              OR IFNULL(f.rating, '') <> IFNULL(t.rating, '')
              OR IFNULL(f.duration, '') <> IFNULL(t.duration, '')
              OR IFNULL(f.tagline, '') <> IFNULL(t.tagline, '')
              OR IFNULL(f.details_json, '{}') <> IFNULL(t.details_json, '{}') THEN 'changed'
            ELSE 'unchanged'
        END AS change_type
    FROM from_snapshots f
    LEFT JOIN to_snapshots t
      ON f.work_href = t.work_href

    UNION ALL

    SELECT
        COALESCE(f.work_href, t.work_href) AS work_href,
        COALESCE(f.category, t.category) AS category,
        CASE WHEN f.work_href IS NULL THEN NULL ELSE f.title END AS from_title,
        t.title AS to_title,
        f.release_date AS from_release_date,
        t.release_date AS to_release_date,
        f.metascore AS from_metascore,
        t.metascore AS to_metascore,
        f.user_score AS from_user_score,
        t.user_score AS to_user_score,
        f.rating AS from_rating,
        t.rating AS to_rating,
        f.duration AS from_duration,
        t.duration AS to_duration,
        f.tagline AS from_tagline,
        t.tagline AS to_tagline,
        CASE WHEN f.work_href IS NULL THEN NULL ELSE f.details_json END AS from_details_json,
        t.details_json AS to_details_json,
        CASE
            WHEN f.work_href IS NOT NULL
             AND t.work_href IS NOT NULL
             AND IFNULL(f.details_json, '{}') <> IFNULL(t.details_json, '{}') THEN 1
            ELSE 0
        END AS details_json_changed,
        CASE
            WHEN f.work_href IS NULL THEN 'added'
            WHEN t.work_href IS NULL THEN 'removed'
            WHEN IFNULL(f.title, '') <> IFNULL(t.title, '')
              OR IFNULL(f.release_date, '') <> IFNULL(t.release_date, '')
              OR IFNULL(f.metascore, '') <> IFNULL(t.metascore, '')
              OR IFNULL(f.user_score, '') <> IFNULL(t.user_score, '')
              OR IFNULL(f.rating, '') <> IFNULL(t.rating, '')
              OR IFNULL(f.duration, '') <> IFNULL(t.duration, '')
              OR IFNULL(f.tagline, '') <> IFNULL(t.tagline, '')
              OR IFNULL(f.details_json, '{}') <> IFNULL(t.details_json, '{}') THEN 'changed'
            ELSE 'unchanged'
        END AS change_type
    FROM to_snapshots t
    LEFT JOIN from_snapshots f
      ON f.work_href = t.work_href
    WHERE f.work_href IS NULL
)
SELECT
    work_href,
    category,
    from_title,
    to_title,
    from_release_date,
    to_release_date,
    from_metascore,
    to_metascore,
    from_user_score,
    to_user_score,
    from_rating,
    to_rating,
    from_duration,
    to_duration,
    from_tagline,
    to_tagline,
    from_details_json,
    to_details_json,
    details_json_changed,
    change_type
FROM compare_rows
WHERE sqlc.arg(include_unchanged) = 1
   OR change_type <> 'unchanged'
ORDER BY category ASC, work_href ASC;

-- name: CountWorkDetails :one
SELECT COUNT(*) FROM work_details;

-- name: CountWorkDetailSnapshots :one
SELECT COUNT(*) FROM work_detail_snapshots;

-- name: CountWorkDetailSnapshotsByWorkHref :one
SELECT COUNT(*) FROM work_detail_snapshots WHERE work_href = ?;
