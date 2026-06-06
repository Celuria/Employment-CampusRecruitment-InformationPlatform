# 就业中心校招信息平台 — 前端

基于 Vue 3 + Vite + TypeScript 的企业级前端工程，UI 设计对齐 `校招平台 Web 端 UI 设计.html`，当前阶段为**框架搭建完成**，业务功能待分模块迭代接入。

## 技术栈

| 类别 | 技术 |
| --- | --- |
| 框架 | Vue 3 (Composition API + `<script setup>`) |
| 构建 | Vite 8 |
| 语言 | TypeScript |
| 路由 | Vue Router 5 |
| 状态 | Pinia 3 |
| UI 组件 | Element Plus（表单/反馈）+ Tailwind CSS 4（布局/视觉） |
| HTTP | Axios |

## 快速开始

```bash
cd frontend
npm install
npm run dev      # 开发 http://localhost:5173
npm run build    # 生产构建
npm run preview  # 预览构建产物
```

### 环境变量

| 变量 | 说明 | 默认值 |
| --- | --- | --- |
| `VITE_APP_TITLE` | 应用标题 | 就业中心校招信息平台 |
| `VITE_API_BASE_URL` | 后端 API 地址（dev 代理目标） | http://localhost:8080 |

开发环境下，`/api` 请求会通过 Vite 代理转发至后端。

## 目录结构

```
src/
├── api/                 # API 层（按业务模块拆分）
│   ├── request.ts       # Axios 实例、拦截器
│   └── modules/         # auth / user / careerTalk / ...
├── components/
│   ├── business/        # 业务组件（卡片、筛选栏、侧边栏）
│   ├── common/          # 通用组件（分页、占位、空状态）
│   └── layout/          # 布局组件（Header、Footer）
├── composables/         # 组合式函数
├── constants/           # 常量与枚举
├── layouts/             # 页面布局（Default / Auth）
├── router/              # 路由定义与守卫
├── stores/              # Pinia 状态
├── styles/              # 全局样式与设计 Token
├── types/               # TypeScript 类型
├── utils/               # 工具函数
└── views/               # 页面视图（按模块分目录）
```

## 路由一览

| 路径 | 页面 | 登录要求 |
| --- | --- | --- |
| `/` | 首页 | 否 |
| `/login` | 登录 | 否（已登录重定向首页） |
| `/register` | 注册 | 否 |
| `/recommendations` | 个性化推荐 | 是 |
| `/career-talks` | 宣讲会列表 | 否 |
| `/career-talks/:id` | 宣讲会详情 | 否 |
| `/job-fairs` | 双选会列表 | 否 |
| `/job-fairs/:id` | 双选会详情 | 否 |
| `/calendar` | 我的日历 | 是 |
| `/profile/info` | 基本资料 | 是 |
| `/profile/preferences` | 偏好设置 | 是 |
| `/profile/reminders` | 提醒记录 | 是 |

---

## 模块功能与开发任务

以下按企业迭代方式划分：**模块职责 → 当前状态 → 待办任务（含优先级）**。

### 1. 工程基础层

**职责**：构建配置、环境变量、路径别名、设计 Token、全局样式。

| 状态 | 已完成 |
| --- | --- |
| 任务 | |
| P0 | Vite + Tailwind 4 + `@` 路径别名 |
| P0 | 设计 Token（brand / ink / warm 色板、card-shadow、btn-primary 等） |
| P1 | ESLint + Prettier 统一代码规范 |
| P2 | 按需引入 Element Plus，减小打包体积 |
| P2 | Husky + lint-staged 提交前检查 |

---

### 2. 布局与导航（`layouts/` + `components/layout/`）

**职责**：全局 Header/Footer、多布局切换、主导航与用户区域。

| 组件/文件 | 功能 |
| --- | --- |
| `DefaultLayout.vue` | 主布局：Header + 内容区 + Footer |
| `AuthLayout.vue` | 认证页居中布局 |
| `AppHeader.vue` | Logo、导航（首页/推荐/宣讲会/双选会/日历）、通知铃铛、用户头像/登录注册 |
| `AppFooter.vue` | 平台介绍、服务链接、帮助、版权 |

| 状态 | 框架 UI 已还原，交互为静态 |
| --- | --- |
| 任务 | |
| P0 | 导航 `active` 态与路由联动（已基本实现） |
| P0 | 未登录访问需鉴权路由 → 跳转登录并携带 `redirect` |
| P1 | 用户下拉菜单：个人中心、退出登录 |
| P1 | 通知铃铛 → 消息列表抽屉（活动变更、提醒状态） |
| P2 | 移动端导航折叠（响应式汉堡菜单） |

---

### 3. 路由与权限（`router/`）

**职责**：路由注册、页面标题、登录守卫、滚动行为。

| 状态 | 路由表 + 守卫骨架已完成 |
| --- | --- |
| 任务 | |
| P0 | 完善 `requiresAuth` 嵌套路由匹配（已实现 `matched.some`） |
| P1 | 首次登录资料未完善 → 引导完善资料（Onboarding） |
| P1 | 路由级 Loading 与页面切换过渡动画 |
| P2 | 面包屑导航组件 |

---

### 4. 状态管理（`stores/`）

**职责**：认证态、用户信息、偏好设置的集中管理。

| Store | 功能 |
| --- | --- |
| `auth` | Token、登录/退出、`userInfo`、初始化 |
| `user` | 偏好设置、资料完善标记 |

| 状态 | Store 骨架已完成 |
| --- | --- |
| 任务 | |
| P0 | 登录成功后写入 Token 并拉取用户信息 |
| P0 | 退出登录清理 Token 与 Store |
| P1 | 资料/偏好变更后刷新推荐缓存 |
| P2 | 持久化「记住登录」策略优化 |

---

### 5. API 层（`api/`）

**职责**：统一请求封装、错误处理、按模块划分接口。

| 模块文件 | 对应后端接口 |
| --- | --- |
| `auth.ts` | POST `/auth/login`、`/auth/register`、`/auth/logout` |
| `user.ts` | GET/PUT `/users/me`、`/users/me/preferences` |
| `careerTalk.ts` | GET `/career-talks`、`/career-talks/{id}` |
| `jobFair.ts` | GET `/job-fairs`、`/job-fairs/{id}` |
| `recommendation.ts` | GET `/recommendations` |
| `calendar.ts` | CRUD `/calendar/events` |
| `reminder.ts` | GET `/reminders/logs` |

| 状态 | 接口函数已定义，待联调 |
| --- | --- |
| 任务 | |
| P0 | 与后端 Swagger 对齐请求/响应字段 |
| P0 | 401 自动跳转登录（已实现） |
| P1 | 请求 Loading 全局/局部策略 |
| P1 | 接口错误码映射与用户友好文案 |
| P2 | 请求取消（路由切换时 abort） |

---

### 6. 认证模块（`views/auth/`）

**职责**：注册、登录、记住登录、失败锁定提示。

| 状态 | 占位页 |
| --- | --- |
| 任务 | |
| P0 | 登录表单（账号、密码、记住登录）+ `loginApi` |
| P0 | 注册表单 + 密码强度校验（≥8 位，含字母数字） |
| P0 | 登录失败 Toast；连续 5 次锁定 15 分钟提示 |
| P1 | 验证码（可选） |
| P2 | 忘记密码 / SSO 对接（待学校确认） |

---

### 7. 宣讲会模块（`views/career-talks/` + `components/business/`）

**职责**：宣讲会列表、搜索筛选、卡片展示、侧边栏、详情、加入日历。

| 组件 | 功能 |
| --- | --- |
| `PageHeader` | 页面标题、描述、结果计数 |
| `CareerTalkFilter` | 关键词搜索、日期/地点/行业筛选、排序 |
| `CareerTalkCard` | 企业 Logo、标题、时间地点、标签、收藏、加入日历 |
| `CareerTalkSidebar` | 即将开始、热门公司、日历快捷入口、邮件提醒 |
| `AppPagination` | 分页 |

| 状态 | 列表页 UI 骨架 + Mock 数据 |
| --- | --- |
| 任务 | |
| P0 | 接入 `getCareerTalkListApi`，替换 Mock，绑定筛选参数 |
| P0 | 「加入日历」→ `addCalendarEventApi`，未登录跳转登录 |
| P0 | 重复添加 Toast「已在日历中」，按钮变「已加入日历」 |
| P1 | 详情页：完整字段展示、报名外链、来源 URL |
| P1 | 列表 Loading 骨架屏、空状态、错误重试 |
| P1 | 筛选 Pill 与 URL Query 同步（可分享链接） |
| P2 | 收藏功能 |
| P2 | 已结束活动置灰（UI 已有样式） |

---

### 8. 双选会模块（`views/job-fairs/`）

**职责**：双选会列表与详情，支持日期区间筛选。

| 状态 | 占位页 |
| --- | --- |
| 任务 | |
| P0 | 复用宣讲会列表模式，创建 `JobFairCard`、`JobFairFilter` |
| P0 | 接入 `getJobFairListApi` / `getJobFairDetailApi` |
| P1 | 日期区间选择器（双选会特有） |
| P1 | 报名截止时间突出展示 |
| P2 | 参与企业数量展示 |

---

### 9. 个性化推荐（`views/recommendations/`）

**职责**：基于学生画像的推荐列表、匹配原因、快捷操作。

| 状态 | 占位页 |
| --- | --- |
| 任务 | |
| P0 | 接入 `getRecommendationsApi` |
| P0 | 匹配原因标签（如「匹配您的意向岗位：Java 开发」） |
| P0 | 「加入日历」「查看详情」 |
| P1 | 未设置偏好 → 引导完善偏好 |
| P1 | 空推荐状态 |
| P2 | 「不感兴趣」反馈并刷新列表 |

---

### 10. 日历模块（`views/calendar/`）

**职责**：日程列表/月视图、编辑备注与提醒、删除、冲突提示。

| 状态 | 占位页 |
| --- | --- |
| 任务 | |
| P0 | 列表视图：活动名、类型、时间、地点、提醒状态 |
| P0 | 删除日程 + 二次确认，取消未发送提醒 |
| P1 | 月视图（FullCalendar 或自研） |
| P1 | 编辑弹窗：个人备注、提醒时间（1h/1d/3d） |
| P1 | 时间重叠非阻断提示 |
| P2 | 周视图 |
| P2 | 官网信息变更标记 |

---

### 11. 个人中心（`views/profile/`）

**职责**：基本资料、偏好设置、提醒记录。

| 子页面 | 功能 |
| --- | --- |
| `ProfileInfoView` | 姓名、学院、专业、年级、意向岗位、电话、邮箱 |
| `PreferencesView` | 意向岗位、偏好城市/公司、特别关注、提醒粒度 |
| `RemindersView` | 历史提醒发送记录与状态 |

| 状态 | Tab 导航 + 占位页 |
| --- | --- |
| 任务 | |
| P0 | 基本资料表单 + `getProfileApi` / `updateProfileApi` |
| P0 | 偏好设置表单 + `getPreferencesApi` / `updatePreferencesApi` |
| P1 | 学院/专业联动下拉（字典数据） |
| P1 | 邮箱脱敏展示 |
| P1 | 提醒记录列表 + 失败重试次数展示 |
| P2 | 头像上传 |

---

### 12. 首页（`views/home/`）

**职责**：平台入口、快捷导航至核心功能。

| 状态 | 四宫格快捷入口 |
| --- | --- |
| 任务 | |
| P1 | 登录后展示个性化欢迎语与待办日程摘要 |
| P2 | 数据概览（本周宣讲会数、已加入日历数） |

---

### 13. 通用组件（`components/common/`）

| 组件 | 任务 |
| --- | --- |
| `PagePlaceholder` | 已有，开发期占位 |
| `AppPagination` | P1 接入真实 total，省略号页码 |
| `EmptyState` | P1 场景化空状态（无搜索结果、无推荐等） |
| — | P1 全局 Toast/Message 封装（基于 Element Plus） |
| — | P1 `Loading` 骨架屏组件 |

---

### 14. 类型与工具（`types/` + `utils/` + `composables/`）

| 任务 | 优先级 |
| --- | --- |
| 与后端 DTO 对齐，补充 OpenAPI 代码生成（可选） | P1 |
| 接入 dayjs 统一日期格式化 | P1 |
| `usePagination` 与列表页联动 | P0 |
| `useAuth.requireAuth` 用于按钮级鉴权 | P0 |

---

## 建议迭代顺序（Sprint 对齐）

| Sprint | 前端重点 |
| --- | --- |
| Sprint 1 | 认证模块、个人资料页、Header 用户态 |
| Sprint 2 | 宣讲会/双选会列表详情、搜索筛选、分页 |
| Sprint 3 | 推荐页、日历 CRUD、偏好设置、提醒记录 |

## 与 UI 设计稿的对应关系

| 设计稿元素 | 代码位置 |
| --- | --- |
| 顶部导航栏 | `AppHeader.vue` |
| 搜索筛选栏 | `CareerTalkFilter.vue` |
| 宣讲会卡片（含收藏/加入日历/已加入） | `CareerTalkCard.vue` |
| 右侧边栏（即将开始/热门公司/日历/邮件） | `CareerTalkSidebar.vue` |
| 分页 | `AppPagination.vue` |
| 底部 Footer | `AppFooter.vue` |
| 品牌色与 card-shadow | `styles/index.css` |

---

**当前版本**：v0.1.0 — 框架搭建完成，可 `npm run dev` 浏览各页面骨架。
