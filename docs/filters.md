# Filter 参考

本文档整理当前 `crawl list` 已支持的过滤参数，以及 Metacritic 当前公开 URL 中可见的枚举值。

说明：

- CLI 多值参数统一使用逗号分隔
- 当前代码不对白名单做强校验，以下内容主要用于使用说明
- `movie` 和 `tv` 使用 `--network`

## 1. 通用年份参数

所有板块都支持：

```bash
--year=2011:2014
```

内部会映射为：

- `releaseYearMin=2011`
- `releaseYearMax=2014`

## 2. Game

### 2.1 参数

- `--platform`
- `--genre`
- `--release-type`
- `--year`

### 2.2 平台

可见平台值：

`3ds,dreamcast,game-boy-advance,gamecube,meta-quest,mobile,nintendo-64,nintendo-ds,nintendo-switch,nintendo-switch-2,pc,ps-vita,ps1,ps2,ps3,ps4,ps5,psp,wii,wii-u,xbox,xbox-360,xbox-one,xbox-series-x`

示例：

```bash
go run ./cmd/metacritic-harvester crawl list --category=game --metric=metascore --platform=pc,ps5
```

### 2.3 发售类型

当前已知值：

- `coming-soon`

示例：

```bash
go run ./cmd/metacritic-harvester crawl list --category=game --metric=newest --release-type=coming-soon
```

### 2.4 游戏类型

可见类型值：

`action,action-adventure,action-puzzle,action-rpg,adventure,application,arcade,beat---%27em---up,board-or-card-game,card-battle,compilation,edutainment,exercise-or-fitness,fighting,first---person-shooter,gambling,general,mmorpg,open---world,party-or-minigame,pinball,platformer,puzzle,racing,real---time-strategy,rhythm,roguelike,rpg,sandbox,shooter,simulation,sports,strategy,survival,tactics,third---person-shooter,trivia-or-game-show,turn---based-strategy,virtual,visual-novel`

示例：

```bash
go run ./cmd/metacritic-harvester crawl list --category=game --metric=metascore --genre=action,rpg
```

## 3. Movie

### 3.1 参数

- `--network`
- `--genre`
- `--release-type`
- `--year`

### 3.2 Network

可见值：

`amazon,apple-tv-plus,criterion-channel,discovery-plus,disney-plus,epix,fubotv,gamespot,hulu,itunes,iva,max,netflix,paramount-plus,peacock,plex,pluto,prime-video,starz,sundance-now,tubi,vudu,youtube-premium`

示例：

```bash
go run ./cmd/metacritic-harvester crawl list --category=movie --metric=userscore --network=netflix,max
```

### 3.3 发售类型

可见值：

- `coming-soon`
- `in-theaters`

### 3.4 类型

可见值：

`action,adventure,animation,biography,comedy,crime,documentary,drama,family,fantasy,film---noir,game---show,history,horror,music,musical,mystery,news,reality---tv,romance,sci---fi,short,sport,talk---show,thriller,war,western`

示例：

```bash
go run ./cmd/metacritic-harvester crawl list --category=movie --metric=userscore --genre=drama,thriller --release-type=coming-soon,in-theaters
```

## 4. TV

### 4.1 参数

- `--network`
- `--genre`
- `--year`

### 4.2 Network

当前可见值与 Movie 相同：

`amazon,apple-tv-plus,criterion-channel,discovery-plus,disney-plus,epix,fubotv,gamespot,hulu,itunes,iva,max,netflix,paramount-plus,peacock,plex,pluto,prime-video,starz,sundance-now,tubi,vudu,youtube-premium`

### 4.3 类型

当前可见值与 Movie 相同：

`action,adventure,animation,biography,comedy,crime,documentary,drama,family,fantasy,film---noir,game---show,history,horror,music,musical,mystery,news,reality---tv,romance,sci---fi,short,sport,talk---show,thriller,war,western`

示例：

```bash
go run ./cmd/metacritic-harvester crawl list --category=tv --metric=newest --network=hulu,netflix --genre=drama,thriller
```
