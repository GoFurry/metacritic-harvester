# metacritic-harvester 中文说明

[English README](../README.md) | [路线图](./roadmap.md) | [使用方式](./usage.md)

`metacritic-harvester` 是一个使用 Go 编写的 Metacritic 公开内容采集工具。当前以 CLI 为主，聚焦公开列表页与详情页的合法采集、SQLite 落库、批量任务和最新结果查询。

## 当前完成度

当前仓库已经完成 Phase 1、Phase 2、Phase 2.5 和 Phase 3 的主体能力：

- `cobra` CLI 入口
- `crawl list`
- `crawl detail`
- `crawl batch`
- `crawl schedule`
- `detail query`
- `detail export`
- `detail compare`
- `latest query`
- `latest export`
- `latest compare`
- `game` / `movie` / `tv`
- `metascore` / `userscore` / `newest`
- 过滤参数
- SQLite + `sqlc`
- 快照表 + 最新表
- `crawl_runs` 批次记录
- `work_details` 当前详情视图
- `work_detail_snapshots` 详情历史快照
- `detail_fetch_state` 详情抓取状态
- 顺序和受控并发批量执行
- 本地 cron 调度

## 当前数据模型

- `works`
  - 作品主表，按 `href` 去重并更新
- `crawl_runs`
  - 每次单任务抓取一条记录，保存 `run_id`、来源、状态、时间和错误
- `list_entries`
  - 历史快照表，保留时序变化
- `latest_list_entries`
  - 当前最新值表，便于查询和导出
- `work_details`
  - 当前详情表，保存详情核心字段和扩展 JSON
- `work_detail_snapshots`
  - 详情历史快照表，按 `work_href + crawl_run_id` 保留每次成功抓取结果
- `detail_fetch_state`
  - 详情抓取状态表，保存成功时间和最近错误

这意味着当前系统同时支持：

- 查看历史变化
- 获取当前最新状态
- 按两个 `run_id` 做快照对比
- 基于已有榜单结果继续抓取详情页数据

## 当前目录结构

```text
cmd/
  metacritic-harvester/
internal/
  app/
  cli/
  config/
  crawler/
  domain/
  source/metacritic/
  storage/
sql/
  schema.sql
  queries/
docs/
examples/
```

## 推荐阅读顺序

- [使用方式](./usage.md)
- [批量任务](./batch-tasks.md)
- [最新记录与对比](./latest.md)
- [调度说明](./scheduling.md)
- [过滤参数](./filters.md)
- [路线图](./roadmap.md)
- [设计文档](../metacritic-archive-go-design.md)

## 下一阶段

Phase 3.5 已经完成，后续更自然的方向是：

1. `crawl reviews`
2. 更多导出与分析视图
3. `serve` 服务化入口
