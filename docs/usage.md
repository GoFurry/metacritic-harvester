# Usage

## Requirements

- Go 1.26+
- network access to public Metacritic pages and APIs
- `sqlc` only if you need to regenerate database code

```bash
go mod tidy
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
sqlc generate
```

## Command overview

```bash
go run ./cmd/metacritic-harvester --help
go run ./cmd/metacritic-harvester crawl --help
go run ./cmd/metacritic-harvester latest --help
go run ./cmd/metacritic-harvester detail --help
go run ./cmd/metacritic-harvester review --help
go run ./cmd/metacritic-harvester serve --help
```

Implemented commands:

- `crawl list`
- `crawl detail`
- `crawl reviews`
- `crawl batch`
- `crawl schedule`
- `serve`
- `latest query / export / compare`
- `detail query / export / compare`
- `review query / export / compare`

## crawl list

```bash
go run ./cmd/metacritic-harvester crawl list --category=game --metric=metascore --pages=0 --db=output/metacritic.db
go run ./cmd/metacritic-harvester crawl list --category=movie --metric=userscore --year=2011:2014 --network=netflix,max --genre=drama,thriller
go run ./cmd/metacritic-harvester crawl list --category=tv --metric=newest --source=auto --pages=2
go run ./cmd/metacritic-harvester crawl list --category=game --metric=newest --timeout=4h --continue-on-error=false
go run ./cmd/metacritic-harvester crawl list --category=game --metric=metascore --rps=4 --burst=8
```

Common flags:

- `--category=game|movie|tv`
- `--metric=metascore|userscore|newest`
- `--source=api|html|auto`
- `--year=YYYY:YYYY`
- `--platform=...` for `game`
- `--network=...` for `movie|tv`
- `--genre=...`
- `--release-type=...` for `game|movie`
- `--pages`
- `--db`
- `--timeout`
- `--continue-on-error`
- `--rps`
- `--burst`
- `--retries`
- `--proxies`
- `--debug`

Notes:

- default source is `api`
- default timeout is `3h`
- default rate limit is `2 RPS` with `burst=2`
- `--continue-on-error=true` by default, so page-level failures are counted and logged but do not fail the whole run
- `--pages=0` means crawl all available list pages
- `--source=html` forces the legacy HTML path
- `--source=auto` means “API first, fallback to HTML on failure”
- list fallback happens at the run level
- `--rps` controls sustained request rate, while `--burst` controls short spikes above that rate
- set `--continue-on-error=false` to restore fail-fast behavior
- context cancellation or `--timeout` expiration still fails the run

## crawl detail

```bash
go run ./cmd/metacritic-harvester crawl detail --db=output/metacritic.db
go run ./cmd/metacritic-harvester crawl detail --db=output/metacritic.db --category=game --limit=20
go run ./cmd/metacritic-harvester crawl detail --db=output/metacritic.db --work-href=https://www.metacritic.com/game/baldurs-gate-3/ --source=auto
go run ./cmd/metacritic-harvester crawl detail --db=output/metacritic.db --category=tv --force
go run ./cmd/metacritic-harvester crawl detail --db=output/metacritic.db --category=movie --timeout=2h --continue-on-error=false
go run ./cmd/metacritic-harvester crawl detail --db=output/metacritic.db --category=game --concurrency=6 --rps=6 --burst=12
```

Common flags:

- `--category=game|movie|tv`
- `--work-href`
- `--limit`
- `--force`
- `--concurrency`
- `--source=api|html|auto`
- `--db`
- `--timeout`
- `--continue-on-error`
- `--rps`
- `--burst`
- `--retries`
- `--proxies`
- `--debug`

Notes:

- detail source defaults to `api`
- default timeout is `3h`
- default rate limit is `2 RPS` with `burst=2`
- `--continue-on-error=true` by default, so per-work failures are counted and logged but do not fail the whole run
- `--limit=0` means process all detail candidates
- `--work-href` accepts normalized absolute URLs and relative `/game/...` style paths
- detail fallback happens per work when `--source=auto`
- the API path can enrich from HTML/Nuxt for fields such as `where_to_buy` and `where_to_watch`
- enrich failure does not turn a successful detail fetch into a failed one
- `--concurrency` controls worker count; `--rps` and `--burst` control the shared HTTP rate limiter
- set `--continue-on-error=false` to stop the run on the first recorded work failure
- context cancellation or `--timeout` expiration still fails the run

## crawl reviews

```bash
go run ./cmd/metacritic-harvester crawl reviews --db=output/metacritic.db --category=game --review-type=critic --limit=10
go run ./cmd/metacritic-harvester crawl reviews --db=output/metacritic.db --category=movie --review-type=user --limit=10
go run ./cmd/metacritic-harvester crawl reviews --db=output/metacritic.db --work-href=https://www.metacritic.com/tv/shogun-2024/ --review-type=all --force
go run ./cmd/metacritic-harvester crawl reviews --db=output/metacritic.db --category=game --review-type=critic --timeout=90m --continue-on-error=false
go run ./cmd/metacritic-harvester crawl reviews --db=output/metacritic.db --category=game --review-type=all --concurrency=3 --rps=4 --burst=8
```

Common flags:

- `--category=game|movie|tv`
- `--review-type=critic|user|all`
- `--work-href`
- `--platform`
- `--limit`
- `--page-size`
- `--max-pages`
- `--concurrency`
- `--force`
- `--db`
- `--timeout`
- `--continue-on-error`
- `--rps`
- `--burst`
- `--retries`
- `--proxies`
- `--debug`

Notes:

- reviews are `API-first`
- default timeout is `3h`
- default rate limit is `2 RPS` with `burst=2`
- `--continue-on-error=true` by default, so scope-level failures are counted and logged but do not fail the whole run
- `--limit=0` means process all review candidates
- review snapshots are written to `review_snapshots`
- current-state rows are written to `latest_reviews`
- recovery is scoped by `work_href + review_type + platform_key`
- `--concurrency` controls scope workers; `--rps` and `--burst` control the shared HTTP rate limiter
- set `--continue-on-error=false` to stop the run on the first recorded scope failure
- context cancellation or `--timeout` expiration still fails the run

## crawl batch

```bash
go run ./cmd/metacritic-harvester crawl batch --file=examples/batch-tasks.yaml
go run ./cmd/metacritic-harvester crawl batch --file=examples/batch-concurrent.yaml --concurrency=2
```

See [Batch tasks](./batch-tasks.md).

## crawl schedule

```bash
go run ./cmd/metacritic-harvester crawl schedule --file=examples/schedule-jobs.yaml
```

See [Scheduling](./scheduling.md).

## latest query / export / compare

```bash
go run ./cmd/metacritic-harvester latest query --db=output/metacritic.db --category=game --metric=metascore --limit=10
go run ./cmd/metacritic-harvester latest export --db=output/metacritic.db --format=csv --output=output/latest.csv
go run ./cmd/metacritic-harvester latest export --db=output/metacritic.db --run-id=<run-id> --profile=summary --format=json --output=output/latest-summary.json
go run ./cmd/metacritic-harvester latest compare --db=output/metacritic.db --from-run-id=<run-a> --to-run-id=<run-b>
```

See [Latest commands](./latest.md).

## detail query / export / compare

```bash
go run ./cmd/metacritic-harvester detail query --db=output/metacritic.db --category=game --format=json
go run ./cmd/metacritic-harvester detail export --db=output/metacritic.db --profile=flat --format=csv --output=output/detail-flat.csv
go run ./cmd/metacritic-harvester detail export --db=output/metacritic.db --run-id=<detail-run-id> --profile=summary --format=json --output=output/detail-summary.json
go run ./cmd/metacritic-harvester detail compare --db=output/metacritic.db --from-run-id=<detail-run-a> --to-run-id=<detail-run-b>
```

Notes:

- `detail query` reads current `work_details`
- `detail export --run-id` reads `work_detail_snapshots`
- `flat` expands common extras into CSV-friendly columns
- `summary` returns aggregated coverage rows

## review query / export / compare

```bash
go run ./cmd/metacritic-harvester review query --db=output/metacritic.db --category=game --review-type=critic --format=json
go run ./cmd/metacritic-harvester review export --db=output/metacritic.db --profile=flat --format=csv --output=output/review-flat.csv
go run ./cmd/metacritic-harvester review export --db=output/metacritic.db --run-id=<review-run-id> --profile=summary --format=json --output=output/review-summary.json
go run ./cmd/metacritic-harvester review compare --db=output/metacritic.db --from-run-id=<review-run-a> --to-run-id=<review-run-b>
```

Notes:

- `raw` keeps `source_payload_json`
- `flat` keeps normalized columns without payload noise
- `summary` groups by run, category, review type, and platform

## serve

Backend only:

```bash
go run ./cmd/metacritic-harvester serve --db=output/metacritic.db
```

Embedded full-stack console:

```bash
go run ./cmd/metacritic-harvester serve --db=output/metacritic.db --full-stack --enable-write
```

Notes:

- default address is `127.0.0.1:36666`
- default mode is read-only
- write operations require `--enable-write`
- live crawl logs stream from `/api/logs/stream`
- browser download exports are available for `latest / detail / review`
- batch and schedule remain CLI-driven

See [Serve](./serve.md).

## Data model summary

Current persistent tables used by the main workflow:

- `works`
- `crawl_runs`
- `list_entries`
- `latest_list_entries`
- `work_details`
- `work_detail_snapshots`
- `detail_fetch_state`
- `latest_reviews`
- `review_snapshots`
- `review_fetch_state`

## Test

```bash
go test ./...
```
