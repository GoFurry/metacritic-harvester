# latest 命令组

`latest` 命令组只读 SQLite，不触发抓取。

说明：

- `latest` 只覆盖榜单当前态与榜单快照对比
- 详情读侧命令已经独立到顶层 `detail` 命令组

## latest query

用途：

- 直接查看 `latest_list_entries`
- 适合快速确认当前最新榜单状态

示例：

```bash
go run ./cmd/metacritic-harvester latest query --db=output/metacritic.db --category=game --metric=metascore
go run ./cmd/metacritic-harvester latest query --db=output/metacritic.db --work-href=https://www.metacritic.com/game/alpha --format=json
```

支持参数：

- `--db`
- `--category`
- `--metric`
- `--work-href`
- `--filter-key`
- `--limit`
- `--format=table|json`

## latest export

用途：

- 将当前最新视图导出为 `CSV` 或 `JSON`
- 适合给后续清洗和分析使用

示例：

```bash
go run ./cmd/metacritic-harvester latest export --db=output/metacritic.db --format=csv --output=output/latest.csv
go run ./cmd/metacritic-harvester latest export --db=output/metacritic.db --category=movie --metric=userscore --format=json --output=output/movie-userscore.json
go run ./cmd/metacritic-harvester latest export --db=output/metacritic.db --run-id=<run-id> --format=json --output=output/run-snapshot.json
go run ./cmd/metacritic-harvester latest export --db=output/metacritic.db --profile=summary --format=csv --output=output/latest-summary.csv
```

输出字段：

- `work_href`
- `category`
- `metric`
- `filter_key`
- `page_no`
- `rank_no`
- `metascore`
- `user_score`
- `last_crawled_at`
- `source_crawl_run_id`

璇存槑锛?

- 涓嶄紶 `--run-id` 鏃跺鍑哄綋鍓?`latest_list_entries`
- 浼犲叆 `--run-id` 鏃跺鍑哄崟鎵规 `list_entries` 蹇収
- `--profile=raw|flat|summary`锛岄粯璁?`raw`
- `raw` 涓?`flat` 绛変环
- `summary` 浼氳緭鍑?`run_id / category / metric / filter_key` 鐨勮仛鍚堟憳瑕?

## latest compare

用途：

- 用两个 `run_id` 对比两次抓取快照
- 判断作品新增、移除、分数变化、排名变化

示例：

```bash
go run ./cmd/metacritic-harvester latest compare --db=output/metacritic.db --from-run-id=<run-a> --to-run-id=<run-b>
go run ./cmd/metacritic-harvester latest compare --db=output/metacritic.db --from-run-id=<run-a> --to-run-id=<run-b> --format=json
go run ./cmd/metacritic-harvester latest compare --db=output/metacritic.db --from-run-id=<run-a> --to-run-id=<run-b> --format=csv --include-unchanged
```

支持参数：

- `--db`
- `--from-run-id`
- `--to-run-id`
- `--category`
- `--metric`
- `--format=table|json|csv`
- `--include-unchanged`

说明：

- `compare` 读取的是 `list_entries`
- `latest_list_entries` 只保存当前一份状态，不足以直接做两次批次比较
- 因此 `run_id` 是后续分析中非常重要的批次标识

## detail 命令组

详情当前态和详情历史对比不在 `latest` 下，而是在独立顶层命令组：

```bash
go run ./cmd/metacritic-harvester detail query --db=output/metacritic.db --category=game
go run ./cmd/metacritic-harvester detail export --db=output/metacritic.db --format=csv --output=output/details.csv
go run ./cmd/metacritic-harvester detail compare --db=output/metacritic.db --from-run-id=<detail-run-a> --to-run-id=<detail-run-b>
```

其中 `detail compare` 读取的是 `work_detail_snapshots`，而不是 `work_details` 当前态。
