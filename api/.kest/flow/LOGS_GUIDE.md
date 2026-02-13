# 📊 Kest 测试记录查看指南

## 🔍 查看测试历史

### 1. 查看最近的测试记录

```bash
# 查看最近 20 条记录（默认）
kest history

# 查看最近 50 条记录
kest history -n 50

# 查看所有项目的历史（跨项目）
kest history --global
```

**输出示例**：
```
ID    TIME                 METHOD URL                                      STATUS DURATION  
------------------------------------------------------------------------------------------
#345  15:02:03 today       GET    http://127.0.0.1:8080/api/v1/users/pr... 200    0ms       
#344  15:02:03 today       POST   http://127.0.0.1:8080/api/v1/login       200    46ms      
#343  15:02:03 today       POST   http://127.0.0.1:8080/api/v1/register    201    47ms      
```

---

### 2. 查看单个请求的详细信息

```bash
# 查看 ID 为 345 的请求详情
kest show 345
```

**输出示例**：
```
════ Record #345 ════
Time: 2026-02-02 15:02:03
               
─── Request ───
GET http://127.0.0.1:8080/api/v1/users/profile

Headers:
  Authorization: Bearer eyJhbGc...
  content-type: application/json
                
─── Response ───
Status: 200    Duration: 0ms

Headers:
  Content-Type: [application/json; charset=utf-8]
  Date: [Mon, 02 Feb 2026 15:02:03 GMT]

Body:
{
  "code": 0,
  "data": {
    "email": "test@kest.io",
    "username": "testuser"
  }
}
```

---

### 3. 重放历史请求

```bash
# 重放 ID 为 345 的请求
kest replay 345

# 重放并修改某些参数
kest replay 345 -H "Authorization: Bearer new_token"
```

---

## 📁 日志文件位置

### 日志目录结构

```
.kest/
├── config.yaml           # 配置文件
├── flow/                 # Flow 测试文件
└── logs/                 # 详细日志文件
    ├── 2026-02-02_15-02-03_POST_api_v1_login.log
    ├── 2026-02-02_15-02-03_GET_api_v1_users_profile.log
    └── ...
```

### 日志文件命名规则

```
格式: YYYY-MM-DD_HH-MM-SS_METHOD_path.log

示例:
2026-02-02_15-02-03_POST_api_v1_login.log
│          │         │    └─ API 路径
│          │         └─ HTTP 方法
│          └─ 时间戳 (时:分:秒)
└─ 日期 (年-月-日)
```

### 查看日志文件

```bash
# 查看所有日志文件
ls -lh .kest/logs/

# 查看特定日志文件
cat .kest/logs/2026-02-02_15-02-03_POST_api_v1_login.log

# 查看最新的 10 个日志
ls -lt .kest/logs/ | head -11

# 搜索包含特定内容的日志
grep -r "error" .kest/logs/

# 查看今天的所有日志
ls .kest/logs/$(date +%Y-%m-%d)*.log
```

---

## 📊 日志文件内容

每个日志文件包含完整的请求和响应信息：

```json
{
  "timestamp": "2026-02-02T15:02:03Z",
  "method": "POST",
  "url": "http://127.0.0.1:8080/api/v1/login",
  "request": {
    "headers": {
      "Content-Type": "application/json"
    },
    "body": "{\"username\":\"testuser\",\"password\":\"***\"}"
  },
  "response": {
    "status": 200,
    "duration_ms": 46,
    "headers": {
      "Content-Type": "application/json"
    },
    "body": "{\"code\":0,\"data\":{\"access_token\":\"...\"}}"
  }
}
```

---

## 🔧 实用技巧

### 1. 查找失败的请求

```bash
# 查看历史中的失败请求（状态码 >= 400）
kest history | grep -E "(4[0-9]{2}|5[0-9]{2})"

# 查看日志文件中的错误
grep -l "\"status\": 4" .kest/logs/*.log
grep -l "\"status\": 5" .kest/logs/*.log
```

### 2. 分析响应时间

```bash
# 查找响应时间超过 100ms 的请求
kest history | grep -E "[1-9][0-9]{2,}ms"

# 查看特定端点的所有请求
ls .kest/logs/*_POST_api_v1_login.log
```

### 3. 导出测试报告

```bash
# 导出历史记录到文件
kest history -n 100 > test-history.txt

# 统计请求数量
ls .kest/logs/*.log | wc -l

# 按日期统计
ls .kest/logs/2026-02-02*.log | wc -l
```

### 4. 清理旧日志

```bash
# 删除 7 天前的日志
find .kest/logs -name "*.log" -mtime +7 -delete

# 只保留最近 100 个日志文件
ls -t .kest/logs/*.log | tail -n +101 | xargs rm -f
```

---

## 🎯 常用场景

### 场景 1: 调试失败的测试

```bash
# 1. 运行测试
kest run .kest/flow/01-auth-flow.flow.md

# 2. 查看最近的历史
kest history -n 10

# 3. 查看失败请求的详情
kest show 345

# 4. 查看对应的日志文件
cat .kest/logs/2026-02-02_15-02-03_POST_api_v1_login.log
```

### 场景 2: 比较两次请求的差异

```bash
# 查看两个请求的详情
kest show 345 > req1.txt
kest show 346 > req2.txt

# 比较差异
diff req1.txt req2.txt
```

### 场景 3: 生成测试报告

```bash
# 统计今天的测试情况
echo "今天的测试统计:"
echo "总请求数: $(ls .kest/logs/$(date +%Y-%m-%d)*.log 2>/dev/null | wc -l)"
echo "成功请求: $(grep -l '\"status\": 2' .kest/logs/$(date +%Y-%m-%d)*.log 2>/dev/null | wc -l)"
echo "失败请求: $(grep -l '\"status\": [45]' .kest/logs/$(date +%Y-%m-%d)*.log 2>/dev/null | wc -l)"
```

---

## 📝 数据库位置

Kest 使用 SQLite 存储历史记录：

```bash
# 数据库位置（通常在用户目录）
~/.kest/kest.db

# 或者在项目目录
.kest/kest.db

# 直接查询数据库
sqlite3 ~/.kest/kest.db "SELECT id, method, url, response_status FROM records ORDER BY id DESC LIMIT 10;"
```

---

## 🎉 总结

Kest 提供了三种方式查看测试记录：

1. **命令行工具**:
   - `kest history` - 查看历史列表
   - `kest show <id>` - 查看详细信息
   - `kest replay <id>` - 重放请求

2. **日志文件**: `.kest/logs/*.log` - 完整的请求/响应日志

3. **SQLite 数据库**: `~/.kest/kest.db` - 结构化存储

**推荐工作流**:
1. 运行测试 → `kest run test.flow.md`
2. 查看历史 → `kest history`
3. 查看详情 → `kest show <id>`
4. 查看日志 → `cat .kest/logs/<file>.log`

---

**Keep Every Step Tested! 🦅**
