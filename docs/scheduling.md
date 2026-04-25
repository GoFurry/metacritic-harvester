# 调度说明

`crawl schedule` 会在前台启动一个本地调度器，按 cron 表达式触发批量任务文件。

## 命令

```bash
go run ./cmd/metacritic-harvester crawl schedule --file=examples/schedule-jobs.yaml
```

## 调度文件结构

```yaml
timezone: Asia/Shanghai

jobs:
  - name: morning-game
    cron: "0 9 * * *"
    batch_file: ../examples/batch-tasks.yaml
    enabled: true
    concurrency: 2

  - name: nightly-tv
    cron: "0 1 * * *"
    batch_file: ../examples/batch-concurrent.yaml
    enabled: true
```

字段说明：

- `timezone`
  - 可选，默认使用本地时区
- `jobs`
  - 至少一个

每个 job 支持：

- `name`
- `cron`
- `batch_file`
- `enabled`
- `concurrency`

## 运行行为

- 计划触发时会重新读取对应的 batch YAML
- 每个单任务抓取仍然会生成自己的 `run_id`
- 调度只负责触发，不做分布式协调
- 收到中断后停止接收新任务，并等待当前任务收尾

## cron 说明

当前支持：

- 标准 5 段 cron
- 可选秒字段

示例：

- `0 9 * * *`
- `0 */6 * * *`
- `*/30 * * * * *`

## 建议

- 如果多个任务写同一个 SQLite，调度时也会自动串行写入
- 如果你希望调度批次获得真正并发收益，优先让不同任务写不同 DB
- 调度文件和 batch 文件尽量分开管理
- 先用小规模任务验证，再挂长期调度
