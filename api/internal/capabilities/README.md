# capabilities

> **Capabilities Layer – 技术能力层**

`capabilities` 提供**可复用、与业务无关的技术处理能力**。

---

## 架构位置

```
modules        →   capabilities   →   infra
(业务模块)          (技术能力)        (基础设施)
```

---

## 可用能力

| 能力 | 包路径 | 说明 |
|------|--------|------|
| **idgen** | `capabilities/idgen` | ID 生成（UUID、Snowflake、NanoID） |
| **crypto** | `capabilities/crypto` | 加密解密、哈希、密码 |

---

## 使用示例

```go
import (
    "github.com/zgiai/zgo/internal/capabilities/idgen"
    "github.com/zgiai/zgo/internal/capabilities/crypto"
)

// ID 生成
id := idgen.UUID()        // UUID v4
id := idgen.Snowflake()   // 雪花算法
id := idgen.NanoID()      // NanoID
id := idgen.ShortID()     // 短 ID

// 加密解密
enc := crypto.NewAESEncryptorFromString("secret-key")
ciphertext, _ := enc.EncryptString("敏感数据")
plaintext, _ := enc.DecryptString(ciphertext)

// 密码哈希
hash, _ := crypto.HashPassword("password")
ok := crypto.VerifyPassword("password", hash)

// HMAC 签名
sig := crypto.HMACSHA256Hex("data", "key")
```

---

## 设计原则

### ✅ 应该做的

- 提供**单一、明确的技术能力**
- 接口命名用**动词 + 对象**（如 `Encrypt()`, `Generate()`）
- 向上层隐藏实现细节
- 只依赖 `infra` 层

### ❌ 不允许做的

- 包含业务逻辑或业务判断
- 依赖 `modules` 层
- 成为 utils/helpers/common

---

## 添加新能力

1. 在 `capabilities/` 下创建新目录
2. 定义接口 + 实现
3. 添加测试
4. 更新本 README
