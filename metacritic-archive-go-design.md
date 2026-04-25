# `metacritic-harvester` 设计文档

## 1. 项目定位

`metacritic-harvester` 是一个基于 Go 的 Metacritic 公开内容采集与本地归档工具。

当前设计聚焦在合法范围内抓取无需登录即可访问的公开页面，并以命令行工具为一阶段交付形态。整体目标不是一次性验证脚本，而是一套可持续演进的采集器架构。

当前一阶段目标：

- 先交付 CLI 版本
- 先支持列表页采集，不实现详情页和评论页抓取
- 先支持 `Games`、`Movies`、`TV Shows` 三个板块
- 先支持 `Metascore`、`User Score`、`Newest Releases` 三类列表
- 采集结果优先写入 SQLite
- 为后续 `detail`、`reviews`、`export`、`serve` 保留扩展接口

---

## 2. 当前基线

当前仓库中的 [main.go](D:/WorkSpace/Go/metacritic-harvester/main.go) 已完成最小可行验证，已经证明以下链路可工作：

- 基于 Colly 访问 Metacritic 列表页
- 设置基础请求头
- 代理轮转
- 简单失败重试
- 分页抓取
- 列表卡片解析
- 结果输出到 JSON

这份设计文档以该最小验证为现实基础，但不继续沿用其“单文件脚本 + 硬编码 URL + 内存聚合 + JSON 输出”的结构，而是重构为面向后续扩展的工程化方案。

---

## 3. 设计边界

### 3.1 当前纳入范围

- 子命令式 CLI
- 单任务列表采集：`crawl list`
- 三板块统一支持：`game`、`movie`、`tv`
- 三类列表统一支持：`metascore`、`userscore`、`newest`
- SQLite 主存储
- URL 构造与页面解析分层
- 代理轮转、失败重试、分页抓取

### 3.2 当前不纳入范围

- 登录态采集
- 浏览器自动化
- 分布式调度
- 高级反反爬策略
- Web 服务实现
- 详情页抓取实现
- 评论抓取实现
- 导出实现

说明：

- `crawl detail`、`crawl reviews`、`export`、`serve` 只在设计中预留，不在一期实现。
- 过滤条件本期不对外开放 CLI 参数，但领域模型和 URL builder 需要保留扩展位。

---

## 4. 总体架构

整体按四层拆分：

1. `cmd`
   - 命令入口与参数解析
2. `internal/app`
   - 任务编排与用例执行
3. `internal/crawler` 与 `internal/source/metacritic`
   - 前者负责通用抓取能力，后者负责 Metacritic 专属规则
4. `internal/storage`
   - SQLite 持久化

核心原则：

- 站点规则只放在 `internal/source/metacritic`
- 抓取器不拼站点 URL，不处理板块差异
- 存储层不理解页面结构
- 应用层不直接依赖 CSS 选择器
- 先闭环列表采集，再扩展详情、评论、导出、服务

---

## 5. CLI 设计

### 5.1 命令形态

对外命令统一设计为子命令：

```text
metacritic-harvester crawl list
metacritic-harvester crawl detail
metacritic-harvester crawl reviews
metacritic-harvester export
metacritic-harvester serve
```

当前一阶段只实现：

```text
metacritic-harvester crawl list
```

其余命令仅在目录和应用层接口上预留位置。

### 5.2 `crawl list` 参数

一期固定支持以下参数：

- `--category=game|movie|tv`
- `--metric=metascore|userscore|newest`
- `--pages=5`
- `--db=output/metacritic.db`
- `--debug=true|false`
- `--retries=3`
- `--proxies=http://127.0.0.1:7897,http://127.0.0.1:7898`

约束：

- `category` 使用 URL 友好内部枚举，不使用展示文案
- `metric` 使用统一内部枚举，`newest` 由 URL builder 映射为 `/new/`
- 本期不开放过滤条件 flags
- 非法枚举值、页数、代理格式必须在 CLI 层报错

---

## 6. 目录结构建议

```text
metacritic-harvester/
├── cmd/
│   └── metacritic-harvester/
│       └── main.go
├── internal/
│   ├── app/
│   │   ├── app.go
│   │   ├── list_service.go
│   │   ├── detail_service.go
│   │   └── review_service.go
│   ├── config/
│   │   ├── flags.go
│   │   └── config.go
│   ├── domain/
│   │   ├── category.go
│   │   ├── metric.go
│   │   ├── filter.go
│   │   ├── task.go
│   │   ├── work.go
│   │   └── list_entry.go
│   ├── crawler/
│   │   ├── collector.go
│   │   ├── retry.go
│   │   └── proxy.go
│   ├── source/
│   │   └── metacritic/
│   │       ├── urls.go
│   │       ├── selectors.go
│   │       ├── list_parser.go
│   │       ├── pagination.go
│   │       ├── detail_parser.go
│   │       └── review_parser.go
│   ├── storage/
│   │   ├── db.go
│   │   ├── schema.go
│   │   ├── list_repository.go
│   │   ├── detail_repository.go
│   │   └── review_repository.go
│   └── pipeline/
│       ├── task_runner.go
│       └── result.go
└── output/
```

说明：

- `detail` 和 `review` 相关文件可先只定义接口或占位结构，不要求一期实现逻辑。
- 当前不需要为 Web 服务单独建层，未来 `serve` 可在 `cmd` 下新增 HTTP 入口并复用 `internal/app`。

---

## 7. 领域模型

### 7.1 Category

```go
type Category string

const (
    CategoryGame  Category = "game"
    CategoryMovie Category = "movie"
    CategoryTV    Category = "tv"
)
```

### 7.2 Metric

```go
type Metric string

const (
    MetricMetascore Metric = "metascore"
    MetricUserScore Metric = "userscore"
    MetricNewest    Metric = "newest"
)
```

### 7.3 Filter

本期保留模型，但默认零值表示“不带过滤条件”。

```go
type Filter struct {
    Platform string
    Genre    string
    Year     string
    Sort     string
}
```

### 7.4 ListTask

`ListTask` 是一期最核心的任务对象。

```go
type ListTask struct {
    Category  Category
    Metric    Metric
    Filter    Filter
    MaxPages  int
    Debug     bool
}
```

### 7.5 Work

`Work` 表示作品主实体，按详情页 URL 唯一。

```go
type Work struct {
    Name        string
    Href        string
    ImageURL    string
    ReleaseDate string
    Category    Category
}
```

### 7.6 ListEntry

`ListEntry` 表示某次列表抓取得到的一条榜单记录。

```go
type ListEntry struct {
    WorkHref   string
    Category   Category
    Metric     Metric
    Page       int
    Rank       int
    Metascore  string
    UserScore  string
    FilterKey  string
    CrawledAt  time.Time
}
```

说明：

- `Work` 与 `ListEntry` 分离后，同一作品可出现在多个榜单和过滤组合里。
- `Metascore` 与 `UserScore` 同时保留，允许在不同列表下部分为空。
- `FilterKey` 用于后续从零值 `Filter` 平滑扩展到多过滤组合采集。

### 7.7 Future Fetch State

为后续阶段预留：

```go
type DetailFetchState struct {
    WorkHref      string
    LastFetchedAt *time.Time
    Status        string
}

type ReviewFetchState struct {
    WorkHref      string
    LastFetchedAt *time.Time
    Status        string
}
```

---

## 8. 核心接口

### 8.1 URL Builder

```go
func BuildListURL(category domain.Category, metric domain.Metric, filter domain.Filter, page int) string
```

职责：

- 根据板块选择 URL 基础路径
- 根据榜单选择 URL 模板
- 在后续阶段拼接过滤条件
- 根据页码统一追加分页参数

要求：

- 不允许在抓取流程中手工拼接 URL
- 分页逻辑不能写死 `game`

### 8.2 列表解析接口

```go
func ParseListItem(e *colly.HTMLElement, page int, rank int, task domain.ListTask) (domain.ListEntry, domain.Work, bool)
```

职责：

- 提取作品基础信息
- 提取当前列表项的榜单分值
- 填充 `page`、`rank`、`category`、`metric`
- 兼容 `game`、`movie`、`tv` 三个板块

### 8.3 分页解析接口

```go
func ParsePagination(e *colly.HTMLElement) int
```

职责：

- 解析站点实际最大页数
- 供任务调度层与 `MaxPages` 取最小值

### 8.4 仓储接口

```go
type ListRepository interface {
    UpsertWork(ctx context.Context, work domain.Work) error
    InsertListEntry(ctx context.Context, entry domain.ListEntry) error
}
```

说明：

- `Work` 与 `ListEntry` 的写入职责分离
- 一期可采用 `works UPSERT + list_entries INSERT`
- 未来再根据幂等性需求补充批量写入接口

### 8.5 任务执行接口

```go
type TaskRunner interface {
    RunListTask(ctx context.Context, task domain.ListTask) error
}
```

作用：

- 统一承接 CLI 命令
- 后续可复用到 `serve` 模式

---

## 9. Metacritic URL 规则

### 9.1 板块映射

固定映射为：

- `game`
- `movie`
- `tv`

### 9.2 列表路径规则

无过滤条件时使用以下规则：

- `metascore`: `/browse/{category}/`
- `userscore`: `/browse/{category}/all/all/all-time/userscore/`
- `newest`: `/browse/{category}/all/all/all-time/new/`

完整示例：

- `https://www.metacritic.com/browse/game/`
- `https://www.metacritic.com/browse/movie/`
- `https://www.metacritic.com/browse/tv/`
- `https://www.metacritic.com/browse/game/all/all/all-time/userscore/`
- `https://www.metacritic.com/browse/game/all/all/all-time/new/`

### 9.3 分页规则

分页统一通过 `?page=N` 追加。

例如：

- 第 1 页：不追加分页参数或统一视为 `page=1`
- 第 2 页：`...?page=2`

要求：

- URL builder 必须统一处理分页
- 不能在分页回调里单独硬编码路径模板

### 9.4 过滤条件扩展策略

本期只定义扩展位，不定义完整规则。

约束：

- 过滤条件只允许出现在 `internal/source/metacritic/urls.go`
- 其他层只传递 `Filter`
- 零值 `Filter` 必须生成当前无过滤 URL

---

## 10. 页面解析设计

### 10.1 选择器集中管理

所有 Metacritic 专属选择器集中在 `selectors.go`。

建议至少收敛以下常量：

```go
const (
    SelectorCard       = `div[data-testid="filter-results"]`
    SelectorTitle      = `h3[data-testid="product-title"] span:last-child`
    SelectorPagination = `nav[data-testid="navigation-pagination"]`
)
```

目的：

- 页面改版时缩小修改范围
- 避免选择器散落在 handler 中

### 10.2 列表项解析

列表解析至少抽取以下字段：

- `name`
- `href`
- `image_url`
- `release_date`
- `metascore`
- `user_score`
- `page`
- `rank`
- `category`
- `metric`
- `crawled_at`

解析策略：

- 尽量提取结构化节点
- 日期仍可先采用当前最小验证中的正则方案
- 分值提取应兼容不同榜单场景下字段缺失
- 解析器返回统一结构，而不是分别定义 `Game`、`Movie`、`TVShow`

### 10.3 分页解析

分页解析职责仅包括：

- 提取页面显示的最大页码
- 返回站点实际页数

任务层职责：

- 根据站点最大页数和 `MaxPages` 计算实际抓取页数
- 调度后续分页访问

---

## 11. 抓取器设计

`internal/crawler` 只负责通用抓取能力：

- 创建 Colly collector
- 配置请求头
- 配置限速
- 配置代理轮转
- 配置失败重试
- 提供统一事件回调挂载点

它不负责：

- 拼接 Metacritic URL
- 判断板块与榜单差异
- 解析列表字段
- 决定如何入库

### 11.1 请求头

沿用最小验证中已被证明有效的基础请求头策略。

### 11.2 代理策略

当前策略保持简单：

- 多代理 round-robin
- 单代理可退化为固定代理
- 重试时自然轮转到后续代理

### 11.3 重试策略

当前策略保持简单：

- 单请求失败后有限次重试
- 固定短暂等待
- 达到上限后记录错误并继续

不在一期引入：

- 错误分级策略
- 指数退避
- 代理健康检查

---

## 12. SQLite 设计

### 12.1 设计目标

SQLite 是当前主存储，而不是可选输出。

原因：

- 支持按条落盘
- 进程中断后数据不丢
- 便于后续增量抓取
- 便于后续详情页和评论页补抓
- 便于后续导出 JSON 和 CSV

### 12.2 表设计

#### 表：`works`

存储作品主实体。

建议字段：

- `href TEXT PRIMARY KEY`
- `name TEXT NOT NULL`
- `image_url TEXT`
- `release_date TEXT`
- `category TEXT NOT NULL`
- `created_at TEXT NOT NULL`
- `updated_at TEXT NOT NULL`

说明：

- 以详情页 `href` 为唯一键
- 一期只存列表页能拿到的基础作品信息

#### 表：`list_entries`

存储列表页快照记录。

建议字段：

- `id INTEGER PRIMARY KEY AUTOINCREMENT`
- `work_href TEXT NOT NULL`
- `category TEXT NOT NULL`
- `metric TEXT NOT NULL`
- `page_no INTEGER NOT NULL`
- `rank_no INTEGER NOT NULL`
- `metascore TEXT`
- `user_score TEXT`
- `filter_key TEXT NOT NULL`
- `crawled_at TEXT NOT NULL`

建议索引：

- `INDEX idx_list_entries_work_href(work_href)`
- `INDEX idx_list_entries_category_metric(category, metric)`
- `INDEX idx_list_entries_filter_key(filter_key)`

说明：

- 不再把作品主信息和榜单快照糊在一张表里
- 同一作品允许出现在不同榜单、不同页次、不同过滤组合中

#### 表：`detail_fetch_state`

预留详情抓取状态。

建议字段：

- `work_href TEXT PRIMARY KEY`
- `status TEXT NOT NULL`
- `last_fetched_at TEXT`
- `last_error TEXT`

#### 表：`review_fetch_state`

预留评论抓取状态。

建议字段：

- `work_href TEXT PRIMARY KEY`
- `status TEXT NOT NULL`
- `last_fetched_at TEXT`
- `last_error TEXT`

### 12.3 写入策略

一期建议：

- `works` 使用 UPSERT
- `list_entries` 使用 INSERT

原因：

- 作品主信息需要按 `href` 归并
- 榜单快照天然允许重复时间点采集

如需避免完全相同的重复快照，可后续在 `list_entries` 上追加唯一约束，例如：

- `UNIQUE(work_href, category, metric, page_no, rank_no, filter_key, crawled_at)`

本期先不强制。

---

## 13. 任务流程

### 13.1 一期 `crawl list` 流程

1. CLI 解析 `crawl list` 参数
2. 构造 `ListTask`
3. 调用 `BuildListURL` 生成起始 URL
4. 创建 collector，并挂载请求头、限速、代理、重试逻辑
5. 访问起始页
6. 解析分页，得到站点最大页数
7. 结合 `MaxPages` 生成实际抓取计划
8. 逐页解析列表项
9. 每条记录拆分为 `Work + ListEntry`
10. 立即写入 SQLite
11. 输出任务汇总信息

### 13.2 未来 `crawl detail` 流程

1. 查询待抓详情的 `work_href`
2. 访问详情页
3. 解析详情字段
4. 更新作品扩展信息
5. 更新 `detail_fetch_state`

### 13.3 未来 `crawl reviews` 流程

1. 查询待抓评论的 `work_href`
2. 访问评论列表页
3. 解析评论与分页
4. 写入评论相关表
5. 更新 `review_fetch_state`

---

## 14. 演进阶段

### Phase 1：列表采集闭环

目标：

- 落地 `crawl list`
- 覆盖 `game`、`movie`、`tv`
- 覆盖 `metascore`、`userscore`、`newest`
- 统一 URL builder
- 统一列表解析
- 结果写入 SQLite

### Phase 2：过滤条件与多任务运行

目标：

- 在 `Filter` 中逐步支持平台、类型、年份等条件
- 支持一批任务顺序执行
- 为定时采集或批量补采做准备

### Phase 3：详情页抓取

目标：

- 实现 `crawl detail`
- 按 `work_href` 做增量补抓
- 将列表页与详情页解析彻底解耦

### Phase 4：评论抓取

目标：

- 实现 `crawl reviews`
- 支持评论分页
- 建立评论去重和增量抓取策略

### Phase 5：导出与服务化

目标：

- 增加 `export`
- 增加 `serve`
- 复用 `app/service + repository + source adapter`
- 只新增 HTTP 入口，不重写底层抓取与存储逻辑

---

## 15. 验收与测试场景

至少覆盖以下测试与验收点：

- URL 构造测试：`3 categories x 3 metrics` 共 9 组基础组合
- 分页 URL 测试：确认分页逻辑不会写死 `game`
- 列表解析测试：同一解析器兼容 `game`、`movie`、`tv`
- 分数字段测试：`metascore`、`userscore`、`newest` 三类列表都能正确提取或允许为空
- SQLite 入库测试：同一 `href` 只保留一个 `work`
- 榜单快照测试：同一 `href` 允许产生多条 `list_entry`
- CLI 参数校验测试：非法 `category`、`metric`、页数、代理格式应直接失败
- 稳健性测试：代理轮转、失败重试、单页场景、无分页节点场景
- 兼容性测试：零值 `Filter` 生成当前无过滤 URL

---

## 16. 实现约束与默认假设

- 当前采集目标仅限 Metacritic 公开可访问内容
- 一期只实现列表采集，不实现详情和评论抓取逻辑
- 但设计中必须预留详情与评论的任务接口、仓储方向和状态表
- CLI 一期采用子命令形式，而不是单命令 flags 堆叠
- SQLite 是当前主链路，不把 JSON/CSV 作为一期主要输出
- `newest` 是对外语义名，URL builder 内部映射到 `/new/`
- `Filter` 当前默认零值，不影响现有无过滤 URL 规则

---

## 17. 设计原则总结

本项目后续实现遵循以下原则：

- 先列表页，后详情页，再评论页
- 先 CLI，后 Web 服务
- 先站点适配层隔离，后横向扩展功能
- 先 SQLite 主存储，后导出层
- 先统一任务模型，后增加批量任务与服务化入口
- 避免将站点规则散落到抓取器、应用层和存储层

