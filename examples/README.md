# Examples

This directory contains ready-to-run examples for the current CLI.

Files:

- `batch-tasks.yaml`: batch example writing into one SQLite database
- `batch-concurrent.yaml`: concurrent batch example writing into different SQLite files
- `batch-multi-db.yaml`: batch example writing into different SQLite files
- `schedule-jobs.yaml`: local cron schedule example
- `commands.md`: copy-paste CLI command examples

Typical workflow:

1. Run one or more `crawl list` tasks to seed `works`
2. Run `crawl detail` against the same SQLite database

```bash
go run ./cmd/metacritic-harvester crawl list --category=game --metric=metascore --pages=1 --db=output/metacritic.db
go run ./cmd/metacritic-harvester crawl detail --db=output/metacritic.db --category=game
```

Run a batch example:

```bash
go run ./cmd/metacritic-harvester crawl batch --file=examples/batch-tasks.yaml
go run ./cmd/metacritic-harvester crawl batch --file=examples/batch-concurrent.yaml --concurrency=2
```

Run the scheduler:

```bash
go run ./cmd/metacritic-harvester crawl schedule --file=examples/schedule-jobs.yaml
```

Run a single-task example:

```bash
go run ./cmd/metacritic-harvester crawl list --category=game --metric=metascore --year=2011:2014 --platform=pc,ps5 --genre=action,rpg --release-type=coming-soon
go run ./cmd/metacritic-harvester crawl detail --db=output/metacritic.db --category=game
```
