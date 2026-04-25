# Command Examples

## Single-task list crawls

### Game / Metascore

```bash
go run ./cmd/metacritic-harvester crawl list --category=game --metric=metascore --pages=1 --db=output/metacritic.db
```

### Game with filters

```bash
go run ./cmd/metacritic-harvester crawl list --category=game --metric=metascore --year=2011:2014 --platform=pc,ps5 --genre=action,rpg --release-type=coming-soon
```

### Movie with filters

```bash
go run ./cmd/metacritic-harvester crawl list --category=movie --metric=userscore --year=2011:2014 --network=netflix,max --genre=drama,thriller --release-type=coming-soon,in-theaters
```

### TV with filters

```bash
go run ./cmd/metacritic-harvester crawl list --category=tv --metric=newest --year=2011:2014 --network=hulu,netflix --genre=drama,thriller
```

## Batch crawls

### Batch into one SQLite database

```bash
go run ./cmd/metacritic-harvester crawl batch --file=examples/batch-tasks.yaml
```

### Batch with concurrency

```bash
go run ./cmd/metacritic-harvester crawl batch --file=examples/batch-concurrent.yaml --concurrency=2
```

### Batch into separate SQLite databases

```bash
go run ./cmd/metacritic-harvester crawl batch --file=examples/batch-multi-db.yaml
```

## Scheduling

```bash
go run ./cmd/metacritic-harvester crawl schedule --file=examples/schedule-jobs.yaml
```

## Detail crawls

### Crawl all pending details in one database

```bash
go run ./cmd/metacritic-harvester crawl detail --db=output/metacritic.db
```

### Crawl only game details

```bash
go run ./cmd/metacritic-harvester crawl detail --db=output/metacritic.db --category=game
```

### Crawl one specific work

```bash
go run ./cmd/metacritic-harvester crawl detail --db=output/metacritic.db --work-href=https://www.metacritic.com/game/baldurs-gate-3
```

### Force refresh successful details

```bash
go run ./cmd/metacritic-harvester crawl detail --db=output/metacritic.db --category=tv --force
```

## Latest data

### Query latest rows

```bash
go run ./cmd/metacritic-harvester latest query --db=output/metacritic.db --category=game --metric=metascore
```

### Export latest rows

```bash
go run ./cmd/metacritic-harvester latest export --db=output/metacritic.db --format=csv --output=output/latest.csv
```

### Compare two runs

```bash
go run ./cmd/metacritic-harvester latest compare --db=output/metacritic.db --from-run-id=<run-a> --to-run-id=<run-b>
```

## Detail read-side

### Query current detail rows

```bash
go run ./cmd/metacritic-harvester detail query --db=output/metacritic.db --category=game
go run ./cmd/metacritic-harvester detail query --db=output/metacritic.db --work-href=https://www.metacritic.com/game/baldurs-gate-3 --format=json
```

### Export current detail rows

```bash
go run ./cmd/metacritic-harvester detail export --db=output/metacritic.db --format=csv --output=output/details.csv
go run ./cmd/metacritic-harvester detail export --db=output/metacritic.db --category=movie --format=json --output=output/movie-details.json
```

### Compare two detail runs

```bash
go run ./cmd/metacritic-harvester detail compare --db=output/metacritic.db --from-run-id=<detail-run-a> --to-run-id=<detail-run-b>
go run ./cmd/metacritic-harvester detail compare --db=output/metacritic.db --from-run-id=<detail-run-a> --to-run-id=<detail-run-b> --format=csv --include-unchanged
```
