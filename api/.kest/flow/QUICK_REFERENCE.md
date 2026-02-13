# 🚀 Kest 快速参考

## 📋 常用命令

```bash
# 运行测试
kest run test.flow.md                    # 运行单个 flow
kest run .kest/flow/*.flow.md            # 运行所有 flow

# 查看历史
kest history                             # 最近 20 条
kest history -n 50                       # 最近 50 条
kest show 345                            # 查看详情

# 环境管理
kest env list                            # 列出所有环境
kest env use staging                     # 切换环境（⚠️ 有 bug，会破坏配置）

# 重放请求
kest replay 345                          # 重放历史请求

# 变量管理
kest vars                                # 查看变量（需要实现）
```

---

## 📁 目录结构

```
.kest/
├── config.yaml              # 环境配置
├── flow/                    # Flow 测试文件
│   ├── 00-smoke-test.flow.md
│   ├── 01-auth-flow.flow.md
│   └── ...
└── logs/                    # 详细日志
    └── 2026-02-02_15-02-03_POST_api_v1_login.log
```

---

## 🌍 多环境配置

**配置文件**: `.kest/config.yaml`

```yaml
environments:
  local:
    base_url: http://127.0.0.1:8080
  staging:
    base_url: https://staging-api.kest.io
  production:
    base_url: https://api.kest.io
active_env: local
```

**切换环境**:
```bash
# ⚠️ 不要用这个命令（会破坏配置文件）
# kest env use staging

# 手动编辑配置文件
vim .kest/config.yaml
# 修改: active_env: staging
```

---

## 📝 Flow 文件语法

```markdown
# Test Name

## Step 1: Description

```kest
POST /api/v1/login
Content-Type: application/json

{
  "username": "test{{$timestamp}}",
  "password": "pass123"
}

[Captures]
token: data.access_token
user_id: data.user.id

[Asserts]
status == 200
body.code == 0
```
```

---

## ✅ 支持的断言

```kest
[Asserts]
status == 200                # HTTP 状态码
body.code == 0               # 响应字段值
body.data.name == "test"     # 嵌套字段
body.data.id != 0            # 不等于

# ❌ 不支持（需要实现）
# duration < 500ms           # 性能断言
# body.data.token exists     # 字段存在性
```

---

## 🔧 内置变量

```kest
{{$timestamp}}               # Unix 时间戳
{{$randomInt}}               # 随机整数
{{$uuid}}                    # UUID（需要实现）
{{captured_var}}             # 之前捕获的变量
```

---

## 📊 查看日志

```bash
# 查看所有日志
ls -lh .kest/logs/

# 查看特定日志
cat .kest/logs/2026-02-02_15-02-03_POST_api_v1_login.log

# 查找失败的请求
grep -l '"status": 4' .kest/logs/*.log
grep -l '"status": 5' .kest/logs/*.log

# 统计今天的测试
ls .kest/logs/$(date +%Y-%m-%d)*.log | wc -l
```

---

## 🎯 测试工作流

```bash
# 1. 确保服务器运行
curl http://127.0.0.1:8080/health

# 2. 运行冒烟测试
kest run .kest/flow/00-smoke-test.flow.md

# 3. 运行完整测试
kest run .kest/flow/01-auth-flow.flow.md

# 4. 查看历史
kest history -n 10

# 5. 查看失败的详情
kest show <id>

# 6. 查看日志文件
cat .kest/logs/<file>.log
```

---

## ⚠️ 已知问题

1. **`kest env use` 会破坏配置文件** - 手动编辑 `active_env`
2. **`kest env list` 不显示环境** - 直接查看配置文件
3. **不支持 `exists` 断言** - 使用 `!= ""` 代替
4. **不支持 `duration` 断言** - 等待实现

---

## 📚 更多文档

- [FLOW_GUIDE.md](../../kest-cli/FLOW_GUIDE.md) - Flow 完整指南
- [MULTI_ENV_GUIDE.md](./MULTI_ENV_GUIDE.md) - 多环境使用
- [LOGS_GUIDE.md](./LOGS_GUIDE.md) - 日志查看指南
- [TEST_REPORT.md](./TEST_REPORT.md) - 测试报告

---

**Keep Every Step Tested! 🦅**
