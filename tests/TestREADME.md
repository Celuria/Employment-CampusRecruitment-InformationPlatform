# 校招信息平台 — 测试说明

本文档描述本项目的单元测试与集成测试策略、目录约定及落地顺序。**当前阶段仅建立目录与规范，尚未编写具体测试代码。**

---

## 1. 现状

| 项 | 状态 |
|----|------|
| 后端测试文件 | 0 个 `*_test.go` |
| 前端测试 | 未配置 Vitest |
| 运行入口 | `backend/Makefile` 中 `make test` → `go test ./...` |
| 架构 | Handler → Service（接口）→ Repository（接口）→ GORM |

项目分层清晰，Repository 与 Service 均定义了接口，适合采用 **testify 断言 + 手写 fake Repository + 表驱动测试** 的方式。

---

## 2. 测试目录

```
tests/
├── TestREADME.md              # 本文档
├── backend/
│   ├── integration/           # 后端集成测试（httptest、MySQL 等，需 -short 外运行）
│   └── fixtures/              # 测试夹具（SQL、JSON mock 数据等）
└── frontend/
    └── unit/                  # 前端 Vitest 用例（utils / constants）
```

### 2.1 与业务代码的位置关系

| 测试类型 | 存放位置 | 说明 |
|----------|----------|------|
| Go 单元测试 | 与被测代码**同目录** | Go 惯例：`pagination_test.go` 放在 `backend/pkg/pagination/` 旁；可访问未导出函数 |
| Go 集成测试 | `tests/backend/integration/` | 使用 `package xxx_test` 导入业务包；依赖 DB/HTTP 的用例加 `t.Skip` 或 `-short` |
| 测试夹具 | `tests/backend/fixtures/` | 如 `init_test_db.sql`、mock 响应 JSON |
| 前端单元测试 | `tests/frontend/unit/` | Vitest 用例；通过 `@/` 别名引用 `frontend/src` |

> **说明**：Go 单元测试文件不放入 `tests/` 根目录，而是与源码同包；`tests/backend/` 主要承载集成测试与共享夹具，避免与 `go test ./...` 的惯例冲突。

---

## 3. 测试金字塔

```
        ┌─────────────┐
        │ E2E / 手工   │  少量，后期可选 Playwright
        ├─────────────┤
        │  集成测试    │  Gin httptest + 测试库 MySQL
        ├─────────────┤
        │  单元测试    │  主力：testify + mock + 表驱动
        └─────────────┘
```

### 3.1 核心原则

1. **默认写单元测试**：无网络、无真实 DB，毫秒级执行。
2. **外部依赖隔离**：TCP / HTTP / DB / 邮件 → mock，或 `t.Skip("需要网络/DB，集成测试时启用")`。
3. **表驱动优先**：同一函数多组 `input → expected`，便于扩展用例。
4. **先测纯函数，再测 Service**：投入小、回归价值高。

### 3.2 参考风格（testify）

与常见 Go 探测/重试测试写法一致：

| 模式 | 用途 | 本项目示例 |
|------|------|------------|
| `assert.True` / `assert.Equal` | 确定性逻辑 | 校区拼接、分页规范化 |
| 测失败路径 | 边界与异常 | 密码错误、无匹配推荐 |
| `t.Skip(...)` | 依赖网络/DB | 集成测试、真实邮件发送 |

---

## 4. 后端测试策略

### 4.1 工具依赖（待引入）

```bash
cd backend
go get github.com/stretchr/testify/assert
go get github.com/stretchr/testify/require
# mock 可选，手写 fake 即可
go get github.com/stretchr/testify/mock
```

### 4.2 运行命令

```bash
# 全量（含集成，若未 Skip）
make test
go test ./... -v

# 单包
go test ./internal/service -v -run TestScoreCareerTalk

# 日常开发：跳过集成测试
go test ./... -short
```

### 4.3 三层策略

#### 第 1 层 — 纯函数（优先，无需 mock）

| 优先级 | 源文件 | 测试重点 |
|--------|--------|----------|
| P0 | `backend/pkg/pagination/pagination.go` | `Normalize` / `Offset` / `ParsePage` 边界 |
| P0 | `backend/internal/model/campus.go` | `CampusLabel` / `FormatEventLocation` |
| P0 | `backend/internal/service/recommendation_service.go` | `textMatch` / `scoreCareerTalk` / `scoreJobFair` |
| P1 | `backend/internal/service/reminder_service.go` | `calculateScheduledTime`（1h/1d/3d） |
| P1 | `backend/pkg/jwt/jwt.go` | Generate → Parse 往返 |
| P1 | `backend/internal/service/auth_service.go` | `validatePassword` / `validateUsername` |

**计划测试文件（与源码同目录）：**

- `backend/pkg/pagination/pagination_test.go`
- `backend/internal/model/campus_test.go`
- `backend/internal/service/recommendation_service_test.go`
- `backend/internal/service/reminder_service_test.go`
- `backend/internal/service/auth_service_test.go`
- `backend/pkg/jwt/jwt_test.go`

#### 第 2 层 — Service 单元测试（mock Repository）

在 `service` 包内手写 `fakeUserRepo` 等结构体，实现 `internal/repository/repository.go` 中的接口。

| Service | 测试重点 | Mock 对象 |
|---------|----------|-----------|
| `AuthService` | 登录成功 / 密码错误 / 账号锁定 | `UserRepository` |
| `RecommendationService` | fallback、分页切片、InCalendar | 各 Repository |
| `CalendarService` | 重复加入日历 | `CalendarRepository` |
| `AdminService` | 重复邮箱创建用户 | `UserRepository` |

不测 GORM SQL 正确性，留给集成层。

#### 第 3 层 — 集成测试（`tests/backend/integration/`）

| 类型 | 做法 | 何时运行 |
|------|------|----------|
| HTTP API | `httptest.NewRecorder` + Gin 路由 | CI nightly / 本地手动 |
| MySQL | 库名 `campus_recruitment_test` + `backend/scripts/init_database.sql` | 发布前 / CI |
| 外部网络 | 始终 `t.Skip` | 不纳入默认流水线 |

集成测试模板约定：

```go
func TestXxx_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("集成测试，使用 -short 时跳过")
    }
    // ...
}
```

**计划测试文件：**

- `tests/backend/integration/auth_api_test.go`
- `tests/backend/integration/recommendation_api_test.go`

**计划夹具：**

- `tests/backend/fixtures/init_test_db.sql`
- `tests/backend/fixtures/users.json`

---

## 5. 前端测试策略（辅助）

后端为主；前端仅覆盖 **utils / constants** 中的纯逻辑，不为每个 `.vue` 写组件测试。

| 源文件 | 测试重点 |
|--------|----------|
| `frontend/src/utils/location.ts` | 与后端 `FormatEventLocation` 行为一致 |
| `frontend/src/utils/format.ts` | 日期 picker ↔ API 格式互转 |
| `frontend/src/utils/calendar.ts` | `buildCalendarDays` 月历网格 |
| `frontend/src/constants/profile.ts` | `getPositionsByCollege` 学院联动 |

### 5.1 工具（待配置）

```bash
cd frontend
npm i -D vitest @vue/test-utils happy-dom
# package.json 增加: "test": "vitest", "test:run": "vitest run"
```

**计划测试文件：**

- `tests/frontend/unit/location.spec.ts`
- `tests/frontend/unit/format.spec.ts`
- `tests/frontend/unit/calendar.spec.ts`
- `tests/frontend/unit/profile.spec.ts`

组件、Pinia Store、API 层暂用手工回归或后期 E2E。

---

## 6. 落地顺序

| 步骤 | 内容 | 产出 |
|------|------|------|
| 1 | 引入 testify；pagination、campus 首批用例 | 2 个 `*_test.go`，`make test` 绿灯 |
| 2 | 推荐打分、提醒时间、auth 校验 | 表驱动纯函数测试 |
| 3 | AuthService / RecommendationService | fake Repository + 核心分支 |
| 4 | 集成测试目录 + `-short` 约定 | `tests/backend/integration/` |
| 5 | 可选：Vitest + 前端 utils | `tests/frontend/unit/` |
| 6 | 可选：GitHub Actions `go test ./... -short` | CI 流水线 |

---

## 7. 不建议测试的范围

- **Gin Handler 逐行**：与 Service 重复；仅对鉴权中间件、错误码保留少量 httptest。
- **Element Plus 组件渲染**：UI 靠手工或 E2E。
- **默认跑真实 MySQL / 发真实邮件**：慢且不稳定，必须 `t.Skip` 或 `-short`。

---

## 8. 预期成果

- 日常 `make test` / `go test ./... -short` 在 3 秒内完成，无外部依赖。
- 推荐、鉴权、分页等核心逻辑具备回归保护。
- 集成测试可选启用，不阻塞本地开发。
- 前端 location / profile 与后端 campus 行为可交叉校验一致性。

---

## 9. 快速索引

| 资源 | 路径 |
|------|------|
| 后端 Makefile | `backend/Makefile` |
| 数据库初始化 | `backend/scripts/init_database.sql` |
| 演示数据 | `backend/scripts/seed_demo_data.sql` |
| Repository 接口 | `backend/internal/repository/repository.go` |
| Service 接口 | `backend/internal/service/service.go` |
