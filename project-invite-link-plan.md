# Project Invite Link Plan

## 目标

为项目成员管理增加“邀请链接”能力，让管理员或 owner 可以生成一个可分享的邀请链接。收到链接的用户在登录后可以接受或拒绝邀请，接受后自动加入项目并获得指定角色。该能力用于补充并逐步替代当前“管理员手动添加成员”的流程。

本方案优先支持：

- 生成项目邀请链接
- 邀请页展示项目信息、角色和过期状态
- 登录后接受邀请加入项目
- 显式拒绝邀请
- 成员页中管理邀请链接和查看邀请状态

## 当前现状

- 项目成员只有直接 CRUD，没有邀请流程：
  - `GET /v1/projects/:id/members`
  - `POST /v1/projects/:id/members`
  - `PATCH /v1/projects/:id/members/:uid`
  - `DELETE /v1/projects/:id/members/:uid`
- 没有“邀请创建 / 接受 / 拒绝 / 失效”相关 API。
- 前端 Members 页面已经存在，可作为邀请入口的主要挂载位置。
- 项目已有两类可复用模式：
  - `api-spec share`：公开 slug 页面
  - `project cli token`：token 生成、哈希存储、过期和权限校验

## 产品设计

### 核心规则

- 邀请链接面向“已有系统账号或新注册账号”的用户。
- 被邀请人必须登录后才能接受或拒绝邀请。
- 邀请链接绑定单个项目和单个目标角色。
- 前端不开放通过邀请直接授予 `owner`，仅允许：
  - `admin`
  - `write`
  - `read`
- 同一邀请链接可以配置为：
  - 单次使用
  - 多次使用
- 邀请链接可以配置过期时间。
- 用户如果已经是该项目成员，则接受邀请时不再重复入库，应返回明确提示。

### 和现有“Add Member”的关系

建议分两阶段推进：

1. 第一阶段：并存
   - 保留当前 `Add Member`
   - 增加 `Invite Link`
2. 第二阶段：切换主入口
   - Members 页主按钮改为 `Generate Invite Link`
   - `Add Member` 收进二级菜单或保留给管理员内部操作

为了降低风险，本次后端和前端先按“并存”设计实现。

## 后端设计

### 数据模型

新增表：`project_invitations`

建议字段：

- `id`
- `project_id`
- `token_hash`
- `token_prefix`
- `slug`
- `role`
- `created_by`
- `status`
- `max_uses`
- `used_count`
- `expires_at`
- `last_used_at`
- `created_at`
- `updated_at`
- `deleted_at`

字段说明：

- `token_hash`
  - 用于服务端校验
  - 不存明文 token
- `token_prefix`
  - 用于管理页展示
- `slug`
  - 可作为公开邀请页路径的一部分
  - 建议直接使用明文 token，或用单独随机值映射
- `role`
  - 仅允许 `admin/write/read`
- `status`
  - `active`
  - `revoked`
  - `expired`
- `max_uses`
  - 默认 `1`
  - `0` 表示不限次数
- `used_count`
  - 每次成功 accept 后累加

不记录“邀请到某个邮箱”这一版能力，先做开放式项目邀请链接。

### 领域能力

新增 `project invitation` 模块，建议放在 `api/internal/modules/projectinvite` 或归入 `project` 模块下的 invitation 子能力。为了避免 project 模块继续膨胀，建议单独模块。

服务层职责：

- 创建邀请
- 获取邀请详情
- 接受邀请
- 拒绝邀请
- 撤销邀请
- 列出某项目下的邀请记录
- 校验邀请有效性

### API 设计

#### 管理侧 API

1. `POST /v1/projects/:id/invitations`

请求体：

```json
{
  "role": "read",
  "expires_at": "2026-05-01T00:00:00Z",
  "max_uses": 1
}
```

响应体：

```json
{
  "id": 1,
  "project_id": 12,
  "role": "read",
  "status": "active",
  "slug": "pji_xxxxx",
  "invite_url": "https://xxx/invite/project/pji_xxxxx",
  "max_uses": 1,
  "used_count": 0,
  "expires_at": "2026-05-01T00:00:00Z",
  "created_by": 7,
  "created_at": "2026-04-18T00:00:00Z",
  "updated_at": "2026-04-18T00:00:00Z"
}
```

权限：

- `admin` 或 `owner`

2. `GET /v1/projects/:id/invitations`

作用：

- 列出项目下邀请记录

权限：

- `admin` 或 `owner`

3. `DELETE /v1/projects/:id/invitations/:inviteId`

作用：

- 撤销邀请
- 撤销后链接不可再 accept

权限：

- `admin` 或 `owner`

#### 邀请页 API

4. `GET /v1/project-invitations/:slug`

作用：

- 公开读取邀请摘要
- 用于未登录和登录态都可打开邀请详情页

返回示例：

```json
{
  "project_id": 12,
  "project_name": "Catalog API",
  "project_slug": "catalog-api",
  "role": "read",
  "status": "active",
  "expires_at": "2026-05-01T00:00:00Z",
  "remaining_uses": 1,
  "requires_auth": true
}
```

说明：

- 不返回敏感管理信息
- 若邀请失效，也返回结构化状态，便于前端渲染“已过期/已撤销”

5. `POST /v1/project-invitations/:slug/accept`

作用：

- 当前登录用户接受邀请并加入项目

权限：

- 已登录

行为：

- 校验邀请状态
- 校验过期
- 校验使用次数
- 校验用户是否已是成员
- 成功后写入 `project_members`
- 更新 `used_count` 和 `last_used_at`

成功响应：

```json
{
  "project_id": 12,
  "member": {
    "user_id": 99,
    "role": "read"
  },
  "redirect_to": "/project/12"
}
```

6. `POST /v1/project-invitations/:slug/reject`

作用：

- 当前登录用户拒绝邀请

权限：

- 已登录

第一版建议行为：

- 不写 project_members
- 记录一条邀请响应日志或最小化 rejection 记录
- 如果不想新增日志表，也可先只返回成功，不做持久化

建议更完整做法：

- 新增 `project_invitation_responses`
  - `invitation_id`
  - `user_id`
  - `action`
  - `created_at`

### 路由与权限

管理侧路由挂到项目域下：

- `/projects/:id/invitations`
- `/projects/:id/invitations/:inviteId`

公开邀请路由不挂项目权限中间件：

- `/project-invitations/:slug`
- `/project-invitations/:slug/accept`
- `/project-invitations/:slug/reject`

其中：

- `GET invitation detail` 可以匿名
- `accept/reject` 需要 auth

### 校验细节

- `role`
  - 仅允许 `admin/write/read`
- `expires_at`
  - 可选
  - 若传入，必须晚于当前时间
- `max_uses`
  - 默认 `1`
  - 必须 `>= 0`
- 接受邀请时：
  - `status != active` 则拒绝
  - 已过期则拒绝
  - 达到使用上限则拒绝
  - 用户已是成员则返回 `409` 或结构化“already_member”

### 复用建议

- token 生成与 hash 存储逻辑复用 `project cli token` 的生成思路
- 公开页渲染和 slug 查询模式复用 `api-spec share`

## 前端设计

### 路由

新增邀请页：

- `web/src/app/invite/project/[slug]/page.tsx`

页面职责：

- 读取邀请详情
- 展示项目名、角色、有效期、状态
- 未登录时引导登录
- 登录后提供 Accept / Reject

### Members 管理页改造

当前 Members 页增加邀请管理区，不移除原有成员列表。

新增功能：

- `Generate Invite Link` 按钮
- 邀请创建弹窗
- 邀请记录列表
- 复制邀请链接
- 撤销邀请

#### 创建邀请弹窗

表单字段：

- `role`
- `expires_at`
- `max_uses`

默认值建议：

- `role = read`
- `max_uses = 1`
- `expires_at = now + 7 days`

提交成功后：

- 显示生成结果
- 支持一键复制 `invite_url`
- 刷新邀请列表

#### 邀请列表

字段：

- 角色
- 状态
- 可用次数
- 已使用次数
- 过期时间
- 创建时间
- 链接复制
- 撤销

状态文案：

- `Active`
- `Expired`
- `Revoked`
- `Used up`

### 邀请页 UX

#### 未登录

- 展示邀请摘要
- “登录后接受邀请”
- 登录成功后回跳原邀请页

建议做法：

- 登录页带 `redirect=/invite/project/:slug`

#### 已登录

展示：

- 项目名
- 目标角色
- 邀请状态
- 是否已过期
- 是否还有剩余次数

按钮：

- `Accept Invitation`
- `Reject Invitation`

#### 接受成功

- toast 提示成功
- 自动跳转 `/project/:projectId`

#### 拒绝成功

- toast 提示已拒绝
- 留在当前页并显示 `Invitation rejected`

### 前端数据层

新增：

- `web/src/types/project-invitation.ts`
- `web/src/services/project-invitation.ts`
- `web/src/hooks/use-project-invitations.ts`

建议 hooks：

- `useProjectInvitations(projectId)`
- `useCreateProjectInvitation(projectId)`
- `useDeleteProjectInvitation(projectId)`
- `useProjectInvitationDetail(slug)`
- `useAcceptProjectInvitation(slug)`
- `useRejectProjectInvitation(slug)`

### 和现有 Members 页的关系

第一版 UI 保持并存：

- 一级主按钮：
  - `Generate Invite Link`
- 二级按钮或 ActionMenu：
  - `Add Member`

如果后续确认完全替代，再把 `Add Member` 下沉。

## 数据迁移

新增 migration：

- 创建 `project_invitations`
- 可选：创建 `project_invitation_responses`

索引建议：

- `project_id`
- `slug` 唯一索引
- `token_hash` 唯一索引
- `status`

## 测试计划

### 后端

单元测试：

- 创建邀请成功
- 非 admin/owner 创建邀请失败
- 邀请过期校验
- 邀请次数耗尽校验
- accept 成功后写入 `project_members`
- 已是成员时 accept 返回冲突
- revoked invitation 无法 accept
- reject 成功

集成测试：

- 生成链接 -> 匿名打开 -> 登录 -> accept -> 成为项目成员
- 生成单次链接 -> 第一个用户 accept 成功 -> 第二个用户 accept 失败

### 前端

单测：

- 邀请状态 badge 显示
- 未登录态按钮和提示渲染
- 登录态 accept/reject 按钮渲染
- Members 页创建邀请弹窗提交

手工验证：

1. admin 生成 read 邀请链接
2. 新用户登录并打开链接
3. 接受邀请后跳入项目
4. Members 列表中能看到该用户
5. 单次链接再次使用时报失效
6. admin 撤销链接后邀请页显示 revoked

## 实施顺序

### Phase 1

- 后端 migration
- invitation model / repository / service / handler / routes
- invitation detail + create + revoke + accept

### Phase 2

- 前端 types / services / hooks
- Members 页增加邀请管理
- 邀请页路由与页面

### Phase 3

- reject 能力
- 登录回跳优化
- 邀请状态与错误态 polish

## 风险与取舍

- 如果完全替代“直接加成员”，管理员给现有内部用户赋权会变慢。
- 如果 accept 邀请需要登录，必须确保登录页支持稳定回跳。
- 如果邀请链接允许多次使用，要明确这是一种“开放式项目加入入口”，安全风险高于单次邀请。

## 建议默认策略

第一版默认配置：

- 仅 `admin/owner` 可生成邀请
- 仅允许 `admin/write/read`
- 默认 7 天过期
- 默认单次使用
- 邀请页允许匿名查看摘要
- accept/reject 必须登录
- Members 页先并存保留 `Add Member`

## 输出文件建议

如果后续要正式实施，建议把本方案拆为两份：

- 产品/交互方案
- 工程实施方案

当前这份文档可直接作为开发 kickoff 的主计划。
