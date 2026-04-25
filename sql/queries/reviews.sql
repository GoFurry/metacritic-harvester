-- name: ListReviewCandidateWorks :many
SELECT
    w.href,
    w.name,
    w.image_url,
    w.release_date,
    w.category,
    w.created_at,
    w.updated_at
FROM works w
WHERE (sqlc.arg(category) = '' OR w.category = sqlc.arg(category))
  AND (sqlc.arg(work_href) = '' OR RTRIM(w.href, '/') = RTRIM(sqlc.arg(work_href), '/'))
ORDER BY w.category ASC, w.href ASC
LIMIT CASE WHEN sqlc.arg(limit_rows) <= 0 THEN -1 ELSE sqlc.arg(limit_rows) END;

-- name: UpsertLatestReview :exec
INSERT INTO latest_reviews (
    review_key,
    external_review_id,
    work_href,
    category,
    review_type,
    platform_key,
    review_url,
    review_date,
    score,
    quote,
    publication_name,
    publication_slug,
    author_name,
    author_slug,
    season_label,
    username,
    user_slug,
    thumbs_up,
    thumbs_down,
    version_label,
    spoiler_flag,
    source_payload_json,
    source_crawl_run_id,
    last_crawled_at,
    created_at,
    updated_at
) VALUES (
    sqlc.arg(review_key),
    sqlc.narg(external_review_id),
    sqlc.arg(work_href),
    sqlc.arg(category),
    sqlc.arg(review_type),
    sqlc.arg(platform_key),
    sqlc.narg(review_url),
    sqlc.narg(review_date),
    sqlc.narg(score),
    sqlc.narg(quote),
    sqlc.narg(publication_name),
    sqlc.narg(publication_slug),
    sqlc.narg(author_name),
    sqlc.narg(author_slug),
    sqlc.narg(season_label),
    sqlc.narg(username),
    sqlc.narg(user_slug),
    sqlc.narg(thumbs_up),
    sqlc.narg(thumbs_down),
    sqlc.narg(version_label),
    sqlc.narg(spoiler_flag),
    sqlc.arg(source_payload_json),
    sqlc.arg(source_crawl_run_id),
    sqlc.arg(last_crawled_at),
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
)
ON CONFLICT(review_key) DO UPDATE SET
    external_review_id = excluded.external_review_id,
    work_href = excluded.work_href,
    category = excluded.category,
    review_type = excluded.review_type,
    platform_key = excluded.platform_key,
    review_url = excluded.review_url,
    review_date = excluded.review_date,
    score = excluded.score,
    quote = excluded.quote,
    publication_name = excluded.publication_name,
    publication_slug = excluded.publication_slug,
    author_name = excluded.author_name,
    author_slug = excluded.author_slug,
    season_label = excluded.season_label,
    username = excluded.username,
    user_slug = excluded.user_slug,
    thumbs_up = excluded.thumbs_up,
    thumbs_down = excluded.thumbs_down,
    version_label = excluded.version_label,
    spoiler_flag = excluded.spoiler_flag,
    source_payload_json = excluded.source_payload_json,
    source_crawl_run_id = excluded.source_crawl_run_id,
    last_crawled_at = excluded.last_crawled_at,
    updated_at = CURRENT_TIMESTAMP;

-- name: InsertReviewSnapshot :exec
INSERT INTO review_snapshots (
    review_key,
    crawl_run_id,
    external_review_id,
    work_href,
    category,
    review_type,
    platform_key,
    review_url,
    review_date,
    score,
    quote,
    publication_name,
    publication_slug,
    author_name,
    author_slug,
    season_label,
    username,
    user_slug,
    thumbs_up,
    thumbs_down,
    version_label,
    spoiler_flag,
    source_payload_json,
    crawled_at
) VALUES (
    sqlc.arg(review_key),
    sqlc.arg(crawl_run_id),
    sqlc.narg(external_review_id),
    sqlc.arg(work_href),
    sqlc.arg(category),
    sqlc.arg(review_type),
    sqlc.arg(platform_key),
    sqlc.narg(review_url),
    sqlc.narg(review_date),
    sqlc.narg(score),
    sqlc.narg(quote),
    sqlc.narg(publication_name),
    sqlc.narg(publication_slug),
    sqlc.narg(author_name),
    sqlc.narg(author_slug),
    sqlc.narg(season_label),
    sqlc.narg(username),
    sqlc.narg(user_slug),
    sqlc.narg(thumbs_up),
    sqlc.narg(thumbs_down),
    sqlc.narg(version_label),
    sqlc.narg(spoiler_flag),
    sqlc.arg(source_payload_json),
    sqlc.arg(crawled_at)
)
ON CONFLICT(review_key, crawl_run_id) DO NOTHING;

-- name: UpsertReviewFetchStateRunning :exec
INSERT INTO review_fetch_state (
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
) VALUES (
    sqlc.arg(work_href),
    sqlc.arg(review_type),
    sqlc.arg(platform_key),
    'running',
    sqlc.arg(last_attempted_at),
    NULL,
    sqlc.arg(last_run_id),
    NULL,
    NULL,
    NULL,
    sqlc.arg(updated_at)
)
ON CONFLICT(work_href, review_type, platform_key) DO UPDATE SET
    status = 'running',
    last_attempted_at = excluded.last_attempted_at,
    last_run_id = excluded.last_run_id,
    last_error = NULL,
    last_error_type = NULL,
    last_error_stage = NULL,
    updated_at = excluded.updated_at;

-- name: UpsertReviewFetchStateSucceeded :exec
INSERT INTO review_fetch_state (
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
) VALUES (
    sqlc.arg(work_href),
    sqlc.arg(review_type),
    sqlc.arg(platform_key),
    'succeeded',
    sqlc.arg(last_attempted_at),
    sqlc.arg(last_fetched_at),
    sqlc.arg(last_run_id),
    NULL,
    NULL,
    NULL,
    sqlc.arg(updated_at)
)
ON CONFLICT(work_href, review_type, platform_key) DO UPDATE SET
    status = 'succeeded',
    last_attempted_at = excluded.last_attempted_at,
    last_fetched_at = excluded.last_fetched_at,
    last_run_id = excluded.last_run_id,
    last_error = NULL,
    last_error_type = NULL,
    last_error_stage = NULL,
    updated_at = excluded.updated_at;

-- name: UpsertReviewFetchStateFailed :exec
INSERT INTO review_fetch_state (
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
) VALUES (
    sqlc.arg(work_href),
    sqlc.arg(review_type),
    sqlc.arg(platform_key),
    'failed',
    sqlc.arg(last_attempted_at),
    NULL,
    sqlc.arg(last_run_id),
    sqlc.narg(last_error),
    sqlc.narg(last_error_type),
    sqlc.narg(last_error_stage),
    sqlc.arg(updated_at)
)
ON CONFLICT(work_href, review_type, platform_key) DO UPDATE SET
    status = 'failed',
    last_attempted_at = excluded.last_attempted_at,
    last_run_id = excluded.last_run_id,
    last_error = excluded.last_error,
    last_error_type = excluded.last_error_type,
    last_error_stage = excluded.last_error_stage,
    updated_at = excluded.updated_at;

-- name: GetReviewFetchState :one
SELECT
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
FROM review_fetch_state
WHERE RTRIM(work_href, '/') = RTRIM(sqlc.arg(work_href), '/')
  AND review_type = sqlc.arg(review_type)
  AND platform_key = sqlc.arg(platform_key)
LIMIT 1;

-- name: ListLatestReviews :many
SELECT
    review_key,
    external_review_id,
    work_href,
    category,
    review_type,
    platform_key,
    review_url,
    review_date,
    score,
    quote,
    publication_name,
    publication_slug,
    author_name,
    author_slug,
    season_label,
    username,
    user_slug,
    thumbs_up,
    thumbs_down,
    version_label,
    spoiler_flag,
    source_payload_json,
    source_crawl_run_id,
    last_crawled_at,
    created_at,
    updated_at
FROM latest_reviews
WHERE (sqlc.arg(category) = '' OR category = sqlc.arg(category))
  AND (sqlc.arg(review_type) = '' OR review_type = sqlc.arg(review_type))
  AND (sqlc.arg(platform_key) = '' OR platform_key = sqlc.arg(platform_key))
  AND (sqlc.arg(work_href) = '' OR RTRIM(work_href, '/') = RTRIM(sqlc.arg(work_href), '/'))
ORDER BY category ASC, review_type ASC, platform_key ASC, last_crawled_at DESC, review_key ASC
LIMIT CASE WHEN sqlc.arg(limit_rows) <= 0 THEN -1 ELSE sqlc.arg(limit_rows) END;

-- name: CompareReviewSnapshots :many
WITH
from_reviews AS (
    SELECT
        rs.review_key,
        rs.work_href,
        rs.category,
        rs.review_type,
        rs.platform_key,
        rs.score,
        rs.quote,
        rs.thumbs_up,
        rs.thumbs_down,
        rs.version_label,
        rs.spoiler_flag
    FROM review_snapshots rs
    WHERE rs.crawl_run_id = sqlc.arg(from_run_id)
      AND (sqlc.arg(category) = '' OR rs.category = sqlc.arg(category))
      AND (sqlc.arg(review_type) = '' OR rs.review_type = sqlc.arg(review_type))
      AND (sqlc.arg(platform_key) = '' OR rs.platform_key = sqlc.arg(platform_key))
      AND (sqlc.arg(work_href) = '' OR RTRIM(rs.work_href, '/') = RTRIM(sqlc.arg(work_href), '/'))
),
to_reviews AS (
    SELECT
        rs.review_key,
        rs.work_href,
        rs.category,
        rs.review_type,
        rs.platform_key,
        rs.score,
        rs.quote,
        rs.thumbs_up,
        rs.thumbs_down,
        rs.version_label,
        rs.spoiler_flag
    FROM review_snapshots rs
    WHERE rs.crawl_run_id = sqlc.arg(to_run_id)
      AND (sqlc.arg(category) = '' OR rs.category = sqlc.arg(category))
      AND (sqlc.arg(review_type) = '' OR rs.review_type = sqlc.arg(review_type))
      AND (sqlc.arg(platform_key) = '' OR rs.platform_key = sqlc.arg(platform_key))
      AND (sqlc.arg(work_href) = '' OR RTRIM(rs.work_href, '/') = RTRIM(sqlc.arg(work_href), '/'))
),
compare_rows AS (
    SELECT
        COALESCE(f.review_key, t.review_key) AS review_key,
        COALESCE(f.work_href, t.work_href) AS work_href,
        COALESCE(f.category, t.category) AS category,
        COALESCE(f.review_type, t.review_type) AS review_type,
        COALESCE(f.platform_key, t.platform_key) AS platform_key,
        f.score AS from_score,
        t.score AS to_score,
        CASE
            WHEN f.score IS NOT NULL AND t.score IS NOT NULL THEN t.score - f.score
            ELSE NULL
        END AS score_diff,
        f.quote AS from_quote,
        t.quote AS to_quote,
        f.thumbs_up AS from_thumbs_up,
        t.thumbs_up AS to_thumbs_up,
        f.thumbs_down AS from_thumbs_down,
        t.thumbs_down AS to_thumbs_down,
        f.version_label AS from_version_label,
        t.version_label AS to_version_label,
        f.spoiler_flag AS from_spoiler_flag,
        t.spoiler_flag AS to_spoiler_flag,
        CASE
            WHEN f.review_key IS NULL THEN 'added'
            WHEN t.review_key IS NULL THEN 'removed'
            WHEN IFNULL(f.score, -999999) <> IFNULL(t.score, -999999)
              OR IFNULL(f.quote, '') <> IFNULL(t.quote, '')
              OR IFNULL(f.thumbs_up, -1) <> IFNULL(t.thumbs_up, -1)
              OR IFNULL(f.thumbs_down, -1) <> IFNULL(t.thumbs_down, -1)
              OR IFNULL(f.version_label, '') <> IFNULL(t.version_label, '')
              OR IFNULL(f.spoiler_flag, -1) <> IFNULL(t.spoiler_flag, -1) THEN 'changed'
            ELSE 'unchanged'
        END AS change_type
    FROM from_reviews f
    LEFT JOIN to_reviews t ON f.review_key = t.review_key

    UNION ALL

    SELECT
        COALESCE(f.review_key, t.review_key) AS review_key,
        COALESCE(f.work_href, t.work_href) AS work_href,
        COALESCE(f.category, t.category) AS category,
        COALESCE(f.review_type, t.review_type) AS review_type,
        COALESCE(f.platform_key, t.platform_key) AS platform_key,
        f.score AS from_score,
        t.score AS to_score,
        CASE
            WHEN f.score IS NOT NULL AND t.score IS NOT NULL THEN t.score - f.score
            ELSE NULL
        END AS score_diff,
        f.quote AS from_quote,
        t.quote AS to_quote,
        f.thumbs_up AS from_thumbs_up,
        t.thumbs_up AS to_thumbs_up,
        f.thumbs_down AS from_thumbs_down,
        t.thumbs_down AS to_thumbs_down,
        f.version_label AS from_version_label,
        t.version_label AS to_version_label,
        f.spoiler_flag AS from_spoiler_flag,
        t.spoiler_flag AS to_spoiler_flag,
        CASE
            WHEN f.review_key IS NULL THEN 'added'
            WHEN t.review_key IS NULL THEN 'removed'
            WHEN IFNULL(f.score, -999999) <> IFNULL(t.score, -999999)
              OR IFNULL(f.quote, '') <> IFNULL(t.quote, '')
              OR IFNULL(f.thumbs_up, -1) <> IFNULL(t.thumbs_up, -1)
              OR IFNULL(f.thumbs_down, -1) <> IFNULL(t.thumbs_down, -1)
              OR IFNULL(f.version_label, '') <> IFNULL(t.version_label, '')
              OR IFNULL(f.spoiler_flag, -1) <> IFNULL(t.spoiler_flag, -1) THEN 'changed'
            ELSE 'unchanged'
        END AS change_type
    FROM to_reviews t
    LEFT JOIN from_reviews f ON f.review_key = t.review_key
    WHERE f.review_key IS NULL
)
SELECT
    review_key,
    work_href,
    category,
    review_type,
    platform_key,
    from_score,
    to_score,
    score_diff,
    from_quote,
    to_quote,
    from_thumbs_up,
    to_thumbs_up,
    from_thumbs_down,
    to_thumbs_down,
    from_version_label,
    to_version_label,
    from_spoiler_flag,
    to_spoiler_flag,
    change_type
FROM compare_rows
WHERE sqlc.arg(include_unchanged) = 1
   OR change_type <> 'unchanged'
ORDER BY category ASC, review_type ASC, platform_key ASC, change_type ASC, work_href ASC, review_key ASC;

-- name: CountLatestReviews :one
SELECT COUNT(*) FROM latest_reviews;

-- name: CountReviewSnapshots :one
SELECT COUNT(*) FROM review_snapshots;
