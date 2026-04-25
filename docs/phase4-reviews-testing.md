# Phase 4 功能测试

本文档用于手工验证 Phase 4 评论抓取链路是否工作正常。

测试目标：

- 验证 `crawl reviews` 能从后端接口抓到评论
- 验证 `latest_reviews / review_snapshots / review_fetch_state / crawl_runs` 写入
- 验证 `review query / export / compare`
- 验证 reviews 任务可接入 `crawl batch`

建议使用单独数据库，避免和日常数据混在一起。

## 0. 前置检查

确认命令和测试通过：

```bash
go test ./...
go run ./cmd/metacritic-harvester crawl reviews --help
go run ./cmd/metacritic-harvester review --help
```

成功标准：

- `go test ./...` 通过
- `crawl reviews` 和 `review` 帮助正常输出

## 1. 准备测试数据库

建议使用新库：

```bash
go run ./cmd/metacritic-harvester crawl list --category=game --metric=metascore --pages=1 --db=output/phase4-test.db
go run ./cmd/metacritic-harvester crawl list --category=movie --metric=metascore --pages=1 --db=output/phase4-test.db
go run ./cmd/metacritic-harvester crawl list --category=tv --metric=metascore --pages=1 --db=output/phase4-test.db
```

成功标准：

- 三条命令都输出非空 `run_id`
- `output/phase4-test.db` 创建成功
- 后续 `crawl reviews` 有候选 `works` 可选

## 2. Game 评论抓取

抓取一个 game 候选作品的 critic + user 评论：

```bash
go run ./cmd/metacritic-harvester crawl reviews --db=output/phase4-test.db --category=game --limit=1 --review-type=all --page-size=20 --max-pages=1
```

成功标准：

- summary 中 `run_id` 非空
- `candidates >= 1`
- `scopes >= 1`
- `reviews >= 1`
- `snapshots` 和 `latest` 大于 0

额外验证 game 平台 scope：

```bash
go run ./cmd/metacritic-harvester review query --db=output/phase4-test.db --category=game --format=table --limit=20
```

观察点：

- `PLATFORM` 列应至少在部分记录中非空
- `TYPE` 同时可看到 `critic` 或 `user`

## 3. Movie 评论抓取

```bash
go run ./cmd/metacritic-harvester crawl reviews --db=output/phase4-test.db --category=movie --limit=1 --review-type=all --page-size=20 --max-pages=1
go run ./cmd/metacritic-harvester review query --db=output/phase4-test.db --category=movie --format=json --limit=10
```

成功标准：

- `crawl reviews` summary 中 `reviews >= 1`
- `review query --format=json` 能看到 `publication_name` 或 `username`
- movie 记录的 `platform_key` 通常为空

## 4. TV 评论抓取

```bash
go run ./cmd/metacritic-harvester crawl reviews --db=output/phase4-test.db --category=tv --limit=1 --review-type=all --page-size=20 --max-pages=1
go run ./cmd/metacritic-harvester review query --db=output/phase4-test.db --category=tv --format=json --limit=10
```

成功标准：

- `crawl reviews` summary 中 `reviews >= 1`
- JSON 中可看到 `review_type`
- 如果命中 critic 评论，可能出现 `season_label`

## 5. 导出测试

CSV 导出：

```bash
go run ./cmd/metacritic-harvester review export --db=output/phase4-test.db --category=game --format=csv --output=output/phase4-game-reviews.csv
```

JSON 导出：

```bash
go run ./cmd/metacritic-harvester review export --db=output/phase4-test.db --category=movie --format=json --output=output/phase4-movie-reviews.json
```

成功标准：

- 两个导出文件都生成成功
- CSV 首行包含 `review_key`
- JSON 中包含 `source_payload_json`

## 6. 快照对比测试

先对同一类作品再抓一轮，强制生成新的 `review_snapshots`：

```bash
go run ./cmd/metacritic-harvester crawl reviews --db=output/phase4-test.db --category=game --limit=1 --review-type=all --page-size=20 --max-pages=1 --force
```

记下两次 game 评论抓取的 `run_id`，然后执行：

```bash
go run ./cmd/metacritic-harvester review compare --db=output/phase4-test.db --from-run-id=<run-a> --to-run-id=<run-b>
go run ./cmd/metacritic-harvester review compare --db=output/phase4-test.db --from-run-id=<run-a> --to-run-id=<run-b> --include-unchanged --format=json
```

成功标准：

- compare 命令执行成功
- 默认输出可能为 0 行，这是允许的
- 加 `--include-unchanged` 后应至少能看到同一批评论的对照结果

## 7. 批量任务测试

创建一个 reviews batch 文件，例如：

```yaml
defaults:
  db: output/phase4-batch.db
  retries: 3
tasks:
  - name: game-reviews
    kind: reviews
    category: game
    review-type: all
    limit: 1
    page-size: 20
    max-pages: 1

  - name: movie-reviews
    kind: reviews
    category: movie
    review-type: critic
    limit: 1
    page-size: 20
    max-pages: 1
```

执行：

```bash
go run ./cmd/metacritic-harvester crawl batch --file=<your-reviews-batch.yaml> --concurrency=2
```

成功标准：

- 输出中出现 `kind=reviews`
- 每个 reviews task 都有独立 `run_id`
- batch summary 中能看到：
  - `review_scopes`
  - `review_fetched`
  - `reviews`
  - `review_snapshots`
  - `review_latest`

## 8. 数据库验证要点

如果你使用 SQLite 可视化工具或 `sqlite3`，建议确认这些表中确实有数据：

- `crawl_runs`
- `latest_reviews`
- `review_snapshots`
- `review_fetch_state`

建议检查：

- `latest_reviews.review_key` 非空
- `review_snapshots.crawl_run_id` 对应 `crawl_runs.run_id`
- `review_fetch_state.status` 会出现 `succeeded`

## 9. 问题排查

如果 `crawl reviews` 输出 `candidates=0`：

- 先确认已经跑过对应 category 的 `crawl list`
- 再确认 `works` 表里确实有该 category 的作品

如果 `game` 评论抓到了但 `platform` 全空：

- 先看当前命中的作品是否只有单平台上下文
- 再看该作品的评论接口是否真的返回平台字段

如果 `review compare` 默认没有结果：

- 这不一定是问题，说明两次快照没有变化
- 用 `--include-unchanged` 再确认对比链路是否正常

## 建议验收顺序

1. `go test ./...`
2. 三类 `crawl list`
3. `crawl reviews` for `game`
4. `crawl reviews` for `movie`
5. `crawl reviews` for `tv`
6. `review query`
7. `review export`
8. `review compare`
9. `crawl batch` with `kind: reviews`
