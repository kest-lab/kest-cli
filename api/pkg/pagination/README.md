# Pagination Package

统一的分页处理包，提供简洁易用的分页功能。

## 核心结构

### Request - 分页请求
```go
type Request struct {
    Page     int    `form:"page" json:"page"`         // 页码，默认 1
    PageSize int    `form:"page_size" json:"page_size"` // 每页大小，默认 10，最大 100
    Keyword  string `form:"keyword" json:"keyword"`    // 搜索关键词
}
```

### Result - 分页元数据
```go
type Result struct {
    Total    int64 `json:"total"`     // 总记录数
    Page     int   `json:"page"`      // 当前页码
    PageSize int   `json:"page_size"` // 每页大小
    LastPage int   `json:"last_page"` // 最后一页
    From     int   `json:"from"`      // 当前页起始位置
    To       int   `json:"to"`        // 当前页结束位置
}
```

## 使用方法

### 1. 在 DTO 中使用
```go
type ListRequest struct {
    pagination.Request  // 继承分页请求
    Status  string `form:"status"`   // 其他查询字段
    Plan    string `form:"plan"`
}
```

### 2. 在 Handler 中提取参数
```go
func (h *Handler) List(c *gin.Context) {
    var req ListRequest
    if err := c.ShouldBindQuery(&req); err != nil {
        response.Error(c, http.StatusBadRequest, err.Error())
        return
    }
    
    // 使用分页参数
    page := req.GetPage()        // 自动处理默认值
    pageSize := req.GetPageSize() // 自动处理限制
    offset := req.GetOffset()     // 计算 SQL 偏移量
}
```

### 3. 在 Repository 中查询
```go
func (r *repository) List(ctx context.Context, req *ListRequest) ([]*Model, int64, error) {
    var items []*Model
    var total int64
    
    // 构建查询
    query := r.db.WithContext(ctx).Table("models")
    if req.Keyword != "" {
        query = query.Where("name LIKE ?", "%"+req.Keyword+"%")
    }
    
    // 计算总数
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    // 分页查询
    offset := req.GetOffset()
    if err := query.Offset(offset).Limit(req.GetPageSize()).Find(&items).Error; err != nil {
        return nil, 0, err
    }
    
    return items, total, nil
}
```

### 4. 使用便捷方法
```go
// 自动分页查询
items, result, err := pagination.Paginate[Model](db, req)
if err != nil {
    return nil, err
}

// 从 Gin 上下文自动分页
items, result, err := pagination.PaginateFromContext[Model](c, db)
```

### 5. 响应数据
```go
// 使用统一的响应格式
response.SuccessWithPagination(c, items, total, req.GetPage(), req.GetPageSize())

// 或者使用 Result 结构
response.Success(c, map[string]interface{}{
    "items": items,
    "meta":  result,
})
```

## 工具方法

### Request 方法
- `GetPage() int` - 获取页码（默认 1）
- `GetPageSize() int` - 获取每页大小（默认 10，最大 100）
- `GetOffset() int` - 获取 SQL 偏移量

### 静态方法
- `FromQuery(query map[string][]string) *Request` - 从查询参数构建
- `FromContext(c *gin.Context) *Request` - 从 Gin 上下文提取
- `BuildResult(total, page, pageSize) *Result` - 构建分页元数据
- `Paginate[T](db, req) ([]T, *Result, error)` - 自动分页查询

## 响应格式

标准 API 响应格式：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [...],        // 数据列表
    "total": 100,          // 总记录数
    "page": 1,             // 当前页码
    "page_size": 20        // 每页大小
  }
}
```

## 注意事项

1. **默认值**: 页码默认 1，每页大小默认 10
2. **限制**: 每页大小最大 100，防止性能问题
3. **兼容性**: 与 `pkg/response.SuccessWithPagination()` 完美配合
4. **类型安全**: 使用泛型确保类型安全

## 迁移指南

如果你之前使用了 `paginator.go`，迁移到新的统一格式：

```go
// 之前
paginator, err := pagination.Paginate[Model](db, page, pageSize)
response.Success(c, paginator.ToMap())

// 现在
items, result, err := pagination.Paginate[Model](db, req)
response.SuccessWithPagination(c, items, result.Total, result.Page, result.PageSize)
```
