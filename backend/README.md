# 就业中心校招信息平台 — 后端

基于 **Go 1.22 + Gin + GORM + MySQL + JWT** 的企业级分层架构，接口规范对齐 [`../frontend/接口文档.md`](../frontend/接口文档.md)。

## 技术栈

| 类别 | 技术 |
| --- | --- |
| 语言 | Go 1.22+ |
| Web 框架 | Gin |
| ORM | GORM |
| 数据库 | MySQL 8.0 |
| 配置 | Viper |
| 日志 | Zap |
| 鉴权 | JWT (golang-jwt) |
| 密码 | bcrypt |

## 快速开始

### 前置条件

- Go 1.22+
- MySQL 8.0

### 1. 安装依赖

```bash
cd backend
go mod tidy
```

### 2. 配置数据库

编辑 `config/config.yaml` 或设置环境变量（参考 `.env.example`）：

```yaml
database:
  host: 127.0.0.1
  port: 3306
  user: root
  password: root
  dbname: campus_recruitment
```

创建数据库：

```sql
CREATE DATABASE campus_recruitment DEFAULT CHARACTER SET utf8mb4;
```

### 3. 启动服务

```bash
make run
# 或
go run ./cmd/server/main.go
```

服务默认监听 **http://localhost:8080**

- 健康检查：`GET /health`
- API 前缀：`/api/v1`

## 项目结构

```
backend/
├── cmd/server/main.go          # 应用入口：加载配置、初始化依赖、启动 HTTP 服务
├── config/
│   ├── config.go               # 配置结构体与 Load 函数
│   └── config.yaml             # 默认配置文件
├── internal/                   # 内部业务代码（不可被外部项目 import）
│   ├── handler/                # HTTP 处理器层（Controller）
│   ├── service/                # 业务逻辑层
│   ├── repository/             # 数据访问层（DAO）
│   ├── model/                  # 数据库实体（Entity）
│   ├── dto/                    # 数据传输对象
│   │   ├── request/            # 请求 DTO
│   │   └── response/           # 响应 DTO
│   ├── middleware/             # Gin 中间件
│   └── router/                 # 路由注册
├── pkg/                        # 可复用公共包
│   ├── apperrors/              # 业务错误码
│   ├── database/               # 数据库连接
│   ├── jwt/                    # JWT 签发与解析
│   ├── logger/                 # 日志
│   ├── pagination/             # 分页工具
│   └── response/               # 统一响应封装
├── migrations/                 # SQL 迁移脚本
├── Makefile                    # 常用命令
├── go.mod
└── README.md
```

## 分层架构说明

```
HTTP Request
    ↓
Router（路由分发 + 中间件链）
    ↓
Handler（参数校验、调用 Service、返回响应）
    ↓
Service（业务逻辑、事务编排、权限校验）
    ↓
Repository（数据库 CRUD，不含业务逻辑）
    ↓
Model（表结构映射）
    ↓
MySQL
```

### 各层职责

| 层级 | 目录 | 职责 | 不应做 |
| --- | --- | --- | --- |
| **入口** | `cmd/server` | 组装依赖、启动服务 | 业务逻辑 |
| **路由** | `internal/router` | 注册路由、挂载中间件 | 业务逻辑 |
| **处理器** | `internal/handler` | 解析请求、调用 Service、写响应 | 直接操作 DB |
| **服务** | `internal/service` | 业务规则、事务、领域逻辑 | HTTP 相关操作 |
| **仓储** | `internal/repository` | 数据持久化、SQL 查询 | 业务判断 |
| **模型** | `internal/model` | 表结构定义 | 请求校验 |
| **DTO** | `internal/dto` | 接口入参/出参结构 | 数据库映射 |
| **中间件** | `internal/middleware` | 鉴权、CORS、日志、错误处理 | 业务逻辑 |
| **公共包** | `pkg/` | 跨层复用工具 | 依赖 internal |

## 文件功能详解

### cmd/server/main.go

应用唯一入口。负责：

1. 加载 `config/config.yaml`
2. 初始化 Logger、Database、JWT
3. GORM AutoMigrate 自动建表
4. 组装 Service → Handler → Router
5. 启动 Gin HTTP Server

### config/

| 文件 | 作用 |
| --- | --- |
| `config.go` | 定义 Server/Database/JWT/CORS/Auth 等配置结构体，`Load()` 读取 YAML |
| `config.yaml` | 默认配置：端口 8080、MySQL 连接、JWT 密钥、CORS 白名单、登录锁定策略 |

### internal/handler/

| 文件 | 作用 |
| --- | --- |
| `handler.go` | 聚合所有 Handler，`HealthHandler` 健康检查 |
| `helper.go` | `bindJSON`/`bindQuery`/`abortError` 等通用辅助 |
| `auth_handler.go` | 注册、登录、退出、验证码 |
| `user_handler.go` | 个人资料、偏好设置 |
| `business_handler.go` | 宣讲会、双选会、推荐、日历、提醒、管理端 |

### internal/service/

| 文件 | 作用 |
| --- | --- |
| `service.go` | 定义 Service 接口、`NewServices()` 依赖注入 |
| `auth_service.go` | 注册（bcrypt 加密）、登录（JWT 签发、失败锁定） |
| `user_service.go` | 资料 CRUD、偏好 Upsert |
| `stub_service.go` | 宣讲会/双选会列表查询；推荐/日历/管理端占位实现 |

### internal/repository/

| 文件 | 作用 |
| --- | --- |
| `repository.go` | 所有 Repository 接口与 GORM 实现：User、Preference、CareerTalk、JobFair |

### internal/model/

| 文件 | 实体表 |
| --- | --- |
| `user.go` | `users` — 用户账号、角色、状态、登录锁定 |
| `user_preference.go` | `user_preferences` — 偏好设置 |
| `career_talk.go` | `career_talks` — 宣讲会 |
| `job_fair.go` | `job_fairs` — 双选会 |
| `calendar_event.go` | `calendar_events` — 日历事件 |
| `reminder_log.go` | `reminder_logs` — 邮件提醒记录 |
| `types.go` | `JSONStrings` 等自定义 GORM 类型 |

### internal/middleware/

| 文件 | 作用 |
| --- | --- |
| `auth.go` | `Auth` 强制鉴权、`OptionalAuth` 可选鉴权、`RequireAdmin` 管理员校验 |
| `error.go` | 统一错误响应、Panic Recovery |
| `logger.go` | 请求访问日志 |

### internal/router/router.go

按权限分组注册路由：

| 路由组 | 中间件 | 说明 |
| --- | --- | --- |
| 公开 | 无 | `/auth/register`、`/auth/login` |
| 可选鉴权 | `OptionalAuth` | 宣讲会/双选会（登录后返回个性化字段） |
| 学生 | `Auth` + `RequireAuth` | 资料、推荐、日历、提醒 |
| 管理端 | `Auth` + `RequireAuth` + `RequireAdmin` | `/admin/*` |

### pkg/

| 包 | 作用 |
| --- | --- |
| `response` | 统一 `{code, message, data}` 响应，分页 `PageResult` |
| `apperrors` | 业务错误码，与接口文档第 14 章对齐 |
| `jwt` | Token 签发/解析，Claims 含 userID/username/role |
| `logger` | Zap JSON 日志 |
| `database` | GORM MySQL 连接池配置 |
| `pagination` | 分页参数规范化（page 从 1 开始，max 100） |

## 已实现 vs 待实现

| 模块 | 状态 | 说明 |
| --- | --- | --- |
| 项目框架 | ✅ | 分层架构、路由、中间件、统一响应 |
| 认证 | ✅ | 注册、登录、JWT、密码校验、登录锁定 |
| 用户资料/偏好 | ✅ | CRUD 骨架 |
| 宣讲会/双选会查询 | ✅ | 分页列表、详情（仅 published） |
| 推荐 | 🔲 | 返回空列表，待实现规则引擎 |
| 日历 | 🔲 | 占位，待实现 CRUD + 去重 |
| 提醒 | 🔲 | 占位，待实现 |
| 管理端 CRUD | 🔲 | 路由已注册，返回 501 占位 |
| 数据同步 | 🔲 | 触发接口占位 |
| 定时任务 | 🔲 | 邮件提醒、信息同步 |

## 权限控制

JWT Payload 包含 `role` 字段：

- `student`：学生，可访问学生端 API
- `admin`：管理员，可访问 `/admin/*`

后端强制校验：

```go
admin := v1.Group("/admin",
    middleware.Auth(jwtManager),
    middleware.RequireAuth(),
    middleware.RequireAdmin(),
)
```

## 统一响应示例

```bash
# 健康检查
curl http://localhost:8080/health

# 注册
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"2022001001","password":"Pass1234","email":"test@edu.cn"}'

# 登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"2022001001","password":"Pass1234"}'

# 宣讲会列表
curl "http://localhost:8080/api/v1/career-talks?page=1&pageSize=10"
```

## 开发规范

1. **Handler** 只做参数绑定和响应，不写 SQL
2. **Service** 承载业务逻辑，返回 `*apperrors.AppError`
3. **Repository** 只做数据访问，不含 HTTP 或业务规则
4. 新增接口先更新 `frontend/接口文档.md`，再实现代码
5. 管理端写操作需记录审计日志（待实现）

## 后续迭代建议

| Sprint | 后端任务 |
| --- | --- |
| Sprint 1 | 完善认证、用户资料、管理员种子数据 |
| Sprint 2 | 宣讲会/双选会完整 CRUD（管理端）、筛选排序 |
| Sprint 3 | 推荐引擎、日历 CRUD、邮件提醒定时任务 |

---

**当前版本**：v0.1.0 — 框架搭建完成，认证与基础查询可用。
