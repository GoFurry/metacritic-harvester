# Reviews API Verification

验证日期：2026-04-24  
验证环境：本地 CLI / `curl.exe` 直连 `backend.metacritic.com`

## 结论

- `Phase 4` 里规划的评论接口模式在 2026-04-24 这次实测中全部可访问。
- `composer .../pages/...-reviews/...` 更适合拿页面上下文、summary、平台和附加信息。
- 真正稳定的评论分页列表应走 `/reviews/metacritic/...`。
- game 评论确实区分平台，movie / tv 当前实测不区分平台。
- user 评论列表当前实测有显式 `id`，critic 评论列表当前实测未见显式 `id`。

## 实测矩阵

| 名称 | 请求 | 结果 | 关键验证 |
| --- | --- | --- | --- |
| `game_critic_page` | `/composer/metacritic/pages/games-critic-reviews/left-4-dead-2/web?platform=xbox-360` | 成功 | 顶层 `components, meta`，`section = Xbox 360` |
| `game_critic_list` | `/reviews/metacritic/critic/games/left-4-dead-2/web?platform=xbox-360&offset=0&limit=2` | 成功 | `totalResults = 76`，返回 2 条评论 |
| `game_user_page` | `/composer/metacritic/pages/games-user-reviews/left-4-dead-2/web?platform=xbox-360` | 成功 | 顶层 `components, meta`，`section = Xbox 360` |
| `game_user_list` | `/reviews/metacritic/user/games/left-4-dead-2/web?platform=xbox-360&offset=0&limit=2` | 成功 | `totalResults = 128`，返回 2 条评论 |
| `movie_critic_page` | `/composer/metacritic/pages/movies-critic-reviews/boyhood/web` | 成功 | 顶层 `components, meta` |
| `movie_critic_list` | `/reviews/metacritic/critic/movies/boyhood/web?offset=0&limit=2` | 成功 | `totalResults = 50`，返回 2 条评论 |
| `movie_user_page` | `/composer/metacritic/pages/movies-user-reviews/boyhood/web` | 成功 | 顶层 `components, meta` |
| `movie_user_list` | `/reviews/metacritic/user/movies/boyhood/web?offset=0&limit=2` | 成功 | `totalResults = 393`，返回 2 条评论 |
| `tv_critic_page` | `/composer/metacritic/pages/shows-critic-reviews/oj-made-in-america/web` | 成功 | 顶层 `components, meta`，包含 `SeasonList` |
| `tv_critic_list` | `/reviews/metacritic/critic/shows/oj-made-in-america/web?offset=0&limit=2` | 成功 | `totalResults = 21`，返回 2 条评论 |
| `tv_user_page` | `/composer/metacritic/pages/shows-user-reviews/oj-made-in-america/web` | 成功 | 顶层 `components, meta`，包含 `SeasonList` |
| `tv_user_list` | `/reviews/metacritic/user/shows/oj-made-in-america/web?offset=0&limit=2` | 成功 | `totalResults = 3`，返回 2 条评论 |

## Composer 页验证

`composer` 页可以稳定拿到：

- `product`
- `critic-score-summary` 或 `user-score-summary`
- game 的平台上下文
- tv 的 `seasons`

但它不适合直接当评论列表主源。六组页面样本里，`ReviewList` 组件都返回了错误态：

| 页面 | `componentName` | `componentType` | `status` |
| --- | --- | --- | --- |
| game critic | `critic-reviews` | `ReviewList` | `400` |
| game user | `user-reviews` | `ReviewList` | `400` |
| movie critic | `critic-reviews` | `ReviewList` | `400` |
| movie user | `user-reviews` | `ReviewList` | `400` |
| tv critic | `critic-reviews` | `ReviewList` | `400` |
| tv user | `user-reviews` | `ReviewList` | `400` |

对应错误信息与现有样本一致：

```text
querystring/filterBySentiment must be equal to one of the allowed values
```

这说明 `composer` 页内嵌的 review-list fetch 在当前请求形态下不稳定，更适合作为：

- 页面上下文源
- summary 源
- game 平台枚举源

不适合作为评论明细抓取主链路。

## 列表接口验证

### Game critic

请求：

```bash
curl.exe -L "https://backend.metacritic.com/reviews/metacritic/critic/games/left-4-dead-2/web?platform=xbox-360&offset=0&limit=2"
```

验证结果：

- 成功返回 `data.totalResults = 76`
- `platform` 参数生效
- 首条评论字段包含：
  - `quote`
  - `score`
  - `url`
  - `date`
  - `author`
  - `authorSlug`
  - `publicationName`
  - `publicationSlug`
  - `reviewedProduct`
  - `platform`

### Game user

请求：

```bash
curl.exe -L "https://backend.metacritic.com/reviews/metacritic/user/games/left-4-dead-2/web?platform=xbox-360&offset=0&limit=2"
```

验证结果：

- 成功返回 `data.totalResults = 128`
- 首条评论字段包含：
  - `id`
  - `quote`
  - `score`
  - `thumbsUp`
  - `thumbsDown`
  - `date`
  - `author`
  - `version`
  - `spoiler`
  - `reviewedProduct`
  - `platform`

### Movie critic

请求：

```bash
curl.exe -L "https://backend.metacritic.com/reviews/metacritic/critic/movies/boyhood/web?offset=0&limit=2"
```

验证结果：

- 成功返回 `data.totalResults = 50`
- 不需要 `platform`
- 首条评论字段包含：
  - `quote`
  - `score`
  - `url`
  - `date`
  - `author`
  - `authorSlug`
  - `publicationName`
  - `publicationSlug`
  - `reviewedProduct`

### Movie user

请求：

```bash
curl.exe -L "https://backend.metacritic.com/reviews/metacritic/user/movies/boyhood/web?offset=0&limit=2"
```

验证结果：

- 成功返回 `data.totalResults = 393`
- 首条评论字段包含：
  - `id`
  - `quote`
  - `score`
  - `thumbsUp`
  - `thumbsDown`
  - `date`
  - `author`
  - `version`
  - `spoiler`
  - `reviewedProduct`

### TV critic

请求：

```bash
curl.exe -L "https://backend.metacritic.com/reviews/metacritic/critic/shows/oj-made-in-america/web?offset=0&limit=2"
```

验证结果：

- 成功返回 `data.totalResults = 21`
- 首条评论字段包含：
  - `quote`
  - `score`
  - `url`
  - `date`
  - `author`
  - `authorSlug`
  - `publicationName`
  - `publicationSlug`
  - `reviewedProduct`
  - `season`

### TV user

请求：

```bash
curl.exe -L "https://backend.metacritic.com/reviews/metacritic/user/shows/oj-made-in-america/web?offset=0&limit=2"
```

验证结果：

- 成功返回 `data.totalResults = 3`
- 首条评论字段包含：
  - `id`
  - `quote`
  - `score`
  - `thumbsUp`
  - `thumbsDown`
  - `date`
  - `author`
  - `version`
  - `spoiler`
  - `reviewedProduct`

## 对 Phase 4 的直接影响

- 默认评论抓取主链路应使用 `/reviews/metacritic/...`
- `composer` 页只做：
  - product 上下文
  - score summary
  - game 平台枚举
  - tv seasons 上下文
- `user review` 的主键应优先使用接口显式 `id`
- `critic review` 当前仍需准备组合业务键兜底
- game 的评论任务模型必须显式支持 `platform`
- tv critic 模型要预留 `season` 字段
