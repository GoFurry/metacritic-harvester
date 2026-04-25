# 批量任务

`crawl batch` 用于从 YAML 文件批量执行列表抓取任务。

## 命令

```bash
go run ./cmd/metacritic-harvester crawl batch --file=examples/batch-tasks.yaml
go run ./cmd/metacritic-harvester crawl batch --file=examples/batch-concurrent.yaml --concurrency=2
```

## 特点

- 支持 `defaults + tasks[]`
- 默认遇错继续执行
- 支持受控并发
- 输出每个任务的 `run_id`
- 最终输出总汇总

## YAML 结构

```yaml
defaults:
  db: output/metacritic.db
  pages: 2
  retries: 3
  debug: false
  concurrency: 2
  proxies:
    - http://127.0.0.1:7897

tasks:
  - name: game-metascore-pc
    category: game
    metric: metascore
    year: "2011:2014"
    platform: [pc, ps5]
    genre: [action, rpg]
    release-type: [coming-soon]

  - name: movie-userscore-streaming
    category: movie
    metric: userscore
    network: [netflix, max]
    genre: [drama, thriller]
    release-type: [coming-soon, in-theaters]

  - category: tv
    metric: newest
    network: [hulu, netflix]
    genre: [drama, thriller]
```

## 字段说明

`defaults` 支持：

- `db`
- `pages`
- `retries`
- `debug`
- `proxies`
- `concurrency`

`task` 支持：

- `name`
- `category`
- `metric`
- `year`
- `platform`
- `network`
- `genre`
- `release-type`
- `pages`
- `db`
- `retries`
- `debug`
- `proxies`

说明：

- `name` 可选，不写时自动生成 `category-metric-index`
- 任务字段优先于 `defaults`
- `concurrency` 当前只支持批次级，不支持单任务级
- 多值字段使用 YAML 数组，不用逗号字符串

## 适用范围

- `game`：`platform / genre / release-type / year`
- `movie`：`network / genre / release-type / year`
- `tv`：`network / genre / year`

## 并发建议

虽然批量任务已经支持并发，但 SQLite 的写入模型仍然有限。当前实现会对同一个 `db` 自动串行化，避免 `SQLITE_BUSY`。建议：

- 先从 `concurrency=2` 开始
- 同库任务可以安全执行，但不会获得真正的写入并发
- 如果任务很多，可考虑拆到不同 DB 文件

## 输出示例

每个任务会输出一行摘要：

```text
task=game-metascore-pc category=game metric=metascore run_id=... pages=2 works=40 list_entries=40 latest_entries=40 failures=0 status=success db=output/metacritic.db
```

最后会输出汇总：

```text
batch summary: total=3 succeeded=3 failed=0 works=120 list_entries=120 latest_entries=120 failures=0
```
