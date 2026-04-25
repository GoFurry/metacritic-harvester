CREATE TABLE IF NOT EXISTS crawl_runs (
    run_id TEXT PRIMARY KEY,
    source TEXT NOT NULL,
    task_name TEXT NOT NULL,
    category TEXT NOT NULL,
    metric TEXT NOT NULL,
    filter_key TEXT NOT NULL,
    started_at TEXT NOT NULL,
    finished_at TEXT,
    status TEXT NOT NULL,
    error_message TEXT
);

CREATE TABLE IF NOT EXISTS works (
    href TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    image_url TEXT,
    release_date TEXT,
    category TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_works_category_href ON works(category, href);

CREATE TABLE IF NOT EXISTS list_entries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    crawl_run_id TEXT NOT NULL DEFAULT '',
    work_href TEXT NOT NULL,
    category TEXT NOT NULL,
    metric TEXT NOT NULL,
    page_no INTEGER NOT NULL,
    rank_no INTEGER NOT NULL,
    metascore TEXT,
    user_score TEXT,
    filter_key TEXT NOT NULL,
    crawled_at TEXT NOT NULL,
    FOREIGN KEY(crawl_run_id) REFERENCES crawl_runs(run_id),
    FOREIGN KEY(work_href) REFERENCES works(href)
);

CREATE INDEX IF NOT EXISTS idx_list_entries_crawl_run_id ON list_entries(crawl_run_id);
CREATE INDEX IF NOT EXISTS idx_list_entries_work_href ON list_entries(work_href);
CREATE INDEX IF NOT EXISTS idx_list_entries_category_metric ON list_entries(category, metric);
CREATE INDEX IF NOT EXISTS idx_list_entries_filter_key ON list_entries(filter_key);
CREATE INDEX IF NOT EXISTS idx_list_entries_compare ON list_entries(crawl_run_id, category, metric, work_href, filter_key);

CREATE TABLE IF NOT EXISTS latest_list_entries (
    work_href TEXT NOT NULL,
    category TEXT NOT NULL,
    metric TEXT NOT NULL,
    filter_key TEXT NOT NULL,
    page_no INTEGER NOT NULL,
    rank_no INTEGER NOT NULL,
    metascore TEXT,
    user_score TEXT,
    source_crawl_run_id TEXT NOT NULL DEFAULT '',
    last_crawled_at TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (work_href, category, metric, filter_key),
    FOREIGN KEY(source_crawl_run_id) REFERENCES crawl_runs(run_id),
    FOREIGN KEY(work_href) REFERENCES works(href)
);

CREATE INDEX IF NOT EXISTS idx_latest_list_entries_run_id ON latest_list_entries(source_crawl_run_id);
CREATE INDEX IF NOT EXISTS idx_latest_list_entries_category_metric ON latest_list_entries(category, metric);
CREATE INDEX IF NOT EXISTS idx_latest_list_entries_filter_key ON latest_list_entries(filter_key);

CREATE TABLE IF NOT EXISTS work_details (
    work_href TEXT PRIMARY KEY,
    category TEXT NOT NULL,
    title TEXT NOT NULL,
    summary TEXT,
    release_date TEXT,
    metascore TEXT,
    metascore_sentiment TEXT,
    metascore_review_count INTEGER,
    user_score TEXT,
    user_score_sentiment TEXT,
    user_score_count INTEGER,
    rating TEXT,
    duration TEXT,
    tagline TEXT,
    details_json TEXT NOT NULL DEFAULT '{}',
    last_fetched_at TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(work_href) REFERENCES works(href)
);

CREATE INDEX IF NOT EXISTS idx_work_details_category ON work_details(category);
CREATE INDEX IF NOT EXISTS idx_work_details_last_fetched_at ON work_details(last_fetched_at);

CREATE TABLE IF NOT EXISTS work_detail_snapshots (
    work_href TEXT NOT NULL,
    crawl_run_id TEXT NOT NULL,
    category TEXT NOT NULL,
    title TEXT NOT NULL,
    summary TEXT,
    release_date TEXT,
    metascore TEXT,
    metascore_sentiment TEXT,
    metascore_review_count INTEGER,
    user_score TEXT,
    user_score_sentiment TEXT,
    user_score_count INTEGER,
    rating TEXT,
    duration TEXT,
    tagline TEXT,
    details_json TEXT NOT NULL DEFAULT '{}',
    fetched_at TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (work_href, crawl_run_id),
    FOREIGN KEY(work_href) REFERENCES works(href),
    FOREIGN KEY(crawl_run_id) REFERENCES crawl_runs(run_id)
);

CREATE INDEX IF NOT EXISTS idx_work_detail_snapshots_crawl_run_id_work_href ON work_detail_snapshots(crawl_run_id, work_href);
CREATE INDEX IF NOT EXISTS idx_work_detail_snapshots_work_href_fetched_at ON work_detail_snapshots(work_href, fetched_at DESC);
CREATE INDEX IF NOT EXISTS idx_work_detail_snapshots_category_work_href_fetched_at ON work_detail_snapshots(category, work_href, fetched_at DESC);

CREATE TABLE IF NOT EXISTS detail_fetch_state (
    work_href TEXT PRIMARY KEY,
    status TEXT NOT NULL,
    last_attempted_at TEXT,
    last_fetched_at TEXT,
    last_run_id TEXT,
    last_error TEXT,
    last_error_type TEXT,
    last_error_stage TEXT,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(last_run_id) REFERENCES crawl_runs(run_id),
    FOREIGN KEY(work_href) REFERENCES works(href)
);

CREATE INDEX IF NOT EXISTS idx_detail_fetch_state_status ON detail_fetch_state(status);
CREATE INDEX IF NOT EXISTS idx_detail_fetch_state_last_attempted_at ON detail_fetch_state(last_attempted_at);
CREATE INDEX IF NOT EXISTS idx_detail_fetch_state_last_run_id ON detail_fetch_state(last_run_id);
CREATE INDEX IF NOT EXISTS idx_detail_fetch_state_status_work_href ON detail_fetch_state(status, work_href);

CREATE TABLE IF NOT EXISTS review_fetch_state (
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
);

CREATE INDEX IF NOT EXISTS idx_review_fetch_state_status ON review_fetch_state(status);
CREATE INDEX IF NOT EXISTS idx_review_fetch_state_last_attempted_at ON review_fetch_state(last_attempted_at);
CREATE INDEX IF NOT EXISTS idx_review_fetch_state_last_run_id ON review_fetch_state(last_run_id);
CREATE INDEX IF NOT EXISTS idx_review_fetch_state_scope ON review_fetch_state(work_href, review_type, platform_key);

CREATE TABLE IF NOT EXISTS latest_reviews (
    review_key TEXT PRIMARY KEY,
    external_review_id TEXT,
    work_href TEXT NOT NULL,
    category TEXT NOT NULL,
    review_type TEXT NOT NULL,
    platform_key TEXT NOT NULL DEFAULT '',
    review_url TEXT,
    review_date TEXT,
    score REAL,
    quote TEXT,
    publication_name TEXT,
    publication_slug TEXT,
    author_name TEXT,
    author_slug TEXT,
    season_label TEXT,
    username TEXT,
    user_slug TEXT,
    thumbs_up INTEGER,
    thumbs_down INTEGER,
    version_label TEXT,
    spoiler_flag INTEGER,
    source_payload_json TEXT NOT NULL DEFAULT '{}',
    source_crawl_run_id TEXT NOT NULL DEFAULT '',
    last_crawled_at TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(source_crawl_run_id) REFERENCES crawl_runs(run_id),
    FOREIGN KEY(work_href) REFERENCES works(href)
);

CREATE INDEX IF NOT EXISTS idx_latest_reviews_work_href ON latest_reviews(work_href);
CREATE INDEX IF NOT EXISTS idx_latest_reviews_category_type_platform ON latest_reviews(category, review_type, platform_key);
CREATE INDEX IF NOT EXISTS idx_latest_reviews_run_id ON latest_reviews(source_crawl_run_id);

CREATE TABLE IF NOT EXISTS review_snapshots (
    review_key TEXT NOT NULL,
    crawl_run_id TEXT NOT NULL,
    external_review_id TEXT,
    work_href TEXT NOT NULL,
    category TEXT NOT NULL,
    review_type TEXT NOT NULL,
    platform_key TEXT NOT NULL DEFAULT '',
    review_url TEXT,
    review_date TEXT,
    score REAL,
    quote TEXT,
    publication_name TEXT,
    publication_slug TEXT,
    author_name TEXT,
    author_slug TEXT,
    season_label TEXT,
    username TEXT,
    user_slug TEXT,
    thumbs_up INTEGER,
    thumbs_down INTEGER,
    version_label TEXT,
    spoiler_flag INTEGER,
    source_payload_json TEXT NOT NULL DEFAULT '{}',
    crawled_at TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (review_key, crawl_run_id),
    FOREIGN KEY(crawl_run_id) REFERENCES crawl_runs(run_id),
    FOREIGN KEY(work_href) REFERENCES works(href)
);

CREATE INDEX IF NOT EXISTS idx_review_snapshots_run_id ON review_snapshots(crawl_run_id);
CREATE INDEX IF NOT EXISTS idx_review_snapshots_work_href ON review_snapshots(work_href);
CREATE INDEX IF NOT EXISTS idx_review_snapshots_scope ON review_snapshots(crawl_run_id, category, review_type, platform_key, work_href);
