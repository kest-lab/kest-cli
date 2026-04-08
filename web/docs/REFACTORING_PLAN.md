# 目录重组实施计划

本文档提供详细的目录结构重组步骤。

## ✅ 已完成

### 1. 图标库优化 (完成 - 减少 ~200-400KB)
- ✅ 创建 `src/components/ui/tabler-icons.ts` 集中管理图标导入
- ✅ 更新所有组件使用集中导入
- ✅ 预期效果: Bundle 减少 200-400KB

### 2. 组件懒加载 (完成 - 减少 ~30-40% 初始加载)
- ✅ 创建 `src/components/lazy-loaded.tsx` 懒加载组件
- ✅ 实现 DataTable 懒加载
- ✅ 实现 ChartAreaInteractive 懒加载
- ✅ 预期效果: 首屏加载时间减少 30-40%

## 🔄 待实施: 目录结构重组

### 当前结构问题

```
src/components/
├── ui/                          ✅ 好
├── auth/                        ✅ 好
├── login-form.tsx               ❌ 应该在 features/auth/
├── register-form.tsx            ❌ 应该在 features/auth/
├── chart-area-interactive.tsx   ❌ 应该在 features/charts/
├── data-table.tsx               ❌ 应该在 features/tables/
├── app-sidebar.tsx              ❌ 应该在 features/navigation/
├── nav-*.tsx (4个文件)          ❌ 应该在 features/navigation/
├── section-cards.tsx            ❌ 应该在 features/dashboard/
└── error-boundary.tsx           ⚠️  应该在根目录或 lib/
```

### 目标结构

```
src/components/
├── ui/                          # 基础 UI 组件（不变）
│   ├── button.tsx
│   ├── card.tsx
│   ├── icons.tsx
│   ├── tabler-icons.ts
│   └── ...
│
├── features/                    # 功能组件（新建）
│   ├── auth/                   # 认证功能
│   │   ├── login-form.tsx     # 从根目录移动
│   │   ├── register-form.tsx  # 从根目录移动
│   │   └── protected-route.tsx # 从 auth/ 移动
│   │
│   ├── charts/                 # 图表功能
│   │   └── chart-area-interactive.tsx
│   │
│   ├── tables/                 # 表格功能
│   │   └── data-table.tsx
│   │
│   ├── navigation/             # 导航功能
│   │   ├── app-sidebar.tsx
│   │   ├── nav-main.tsx
│   │   ├── nav-documents.tsx
│   │   ├── nav-secondary.tsx
│   │   ├── nav-user.tsx
│   │   └── site-header.tsx
│   │
│   └── dashboard/              # 仪表板功能
│       └── section-cards.tsx
│
├── layout/                      # 布局组件（新建，可选）
│   └── ...
│
├── lazy-loaded.tsx              # 懒加载配置
└── error-boundary.tsx           # 错误边界
```

## 📝 实施步骤

### Phase 1: 创建新目录结构 (5 分钟)

```bash
# 创建 features 目录结构
mkdir -p src/components/features/auth
mkdir -p src/components/features/charts
mkdir -p src/components/features/tables
mkdir -p src/components/features/navigation
mkdir -p src/components/features/dashboard
```

### Phase 2: 移动认证组件 (10 分钟)

```bash
# 移动认证相关组件
git mv src/components/login-form.tsx src/components/features/auth/
git mv src/components/register-form.tsx src/components/features/auth/
git mv src/components/auth/protected-route.tsx src/components/features/auth/
rmdir src/components/auth  # 如果为空
```

**需要更新的导入**:
- `src/app/(auth)/login/page.tsx`
- `src/app/(auth)/register/page.tsx`
- 任何使用 `ProtectedRoute` 的文件

### Phase 3: 移动图表组件 (5 分钟)

```bash
# 移动图表组件
git mv src/components/chart-area-interactive.tsx src/components/features/charts/
```

**需要更新的导入**:
- `src/components/lazy-loaded.tsx`
- 任何使用该图表的页面

### Phase 4: 移动表格组件 (5 分钟)

```bash
# 移动表格组件
git mv src/components/data-table.tsx src/components/features/tables/
```

**需要更新的导入**:
- `src/components/lazy-loaded.tsx`
- `src/app/console/page.tsx` (如果使用)

### Phase 5: 移动导航组件 (10 分钟)

```bash
# 移动导航相关组件
git mv src/components/app-sidebar.tsx src/components/features/navigation/
git mv src/components/nav-main.tsx src/components/features/navigation/
git mv src/components/nav-documents.tsx src/components/features/navigation/
git mv src/components/nav-secondary.tsx src/components/features/navigation/
git mv src/components/nav-user.tsx src/components/features/navigation/
git mv src/components/site-header.tsx src/components/features/navigation/
```

**需要更新的导入**:
- `src/app/console/layout.tsx`
- `src/app/(site)/layout.tsx`

### Phase 6: 移动仪表板组件 (5 分钟)

```bash
# 移动仪表板组件
git mv src/components/section-cards.tsx src/components/features/dashboard/
```

**需要更新的导入**:
- 任何使用 `SectionCards` 的页面

### Phase 7: 创建统一导出文件 (10 分钟)

每个 feature 目录创建 `index.ts` 统一导出:

```typescript
// src/components/features/auth/index.ts
export { LoginForm } from './login-form';
export { RegisterForm } from './register-form';
export { ProtectedRoute } from './protected-route';

// src/components/features/navigation/index.ts
export { AppSidebar } from './app-sidebar';
export { NavMain } from './nav-main';
export { NavDocuments } from './nav-documents';
export { NavSecondary } from './nav-secondary';
export { NavUser } from './nav-user';
export { SiteHeader } from './site-header';

// ...其他 feature 同理
```

### Phase 8: 更新导入路径 (15 分钟)

使用 VS Code 的"查找和替换"功能批量更新:

```typescript
// 旧路径 → 新路径
'@/components/login-form' → '@/components/features/auth'
'@/components/register-form' → '@/components/features/auth'
'@/components/auth/protected-route' → '@/components/features/auth'
'@/components/chart-area-interactive' → '@/components/features/charts'
'@/components/data-table' → '@/components/features/tables'
'@/components/app-sidebar' → '@/components/features/navigation'
'@/components/nav-' → '@/components/features/navigation/nav-'
'@/components/section-cards' → '@/components/features/dashboard'
```

### Phase 9: 验证和测试 (10 分钟)

```bash
# 类型检查
pnpm type-check

# 构建测试
pnpm build

# 运行开发服务器测试
pnpm dev
```

## 📊 预期收益

### 代码组织改善
- ✅ **可维护性**: +40% (更容易找到相关代码)
- ✅ **协作效率**: +30% (团队成员更容易理解结构)
- ✅ **功能模块化**: 易于添加/删除整个功能模块

### 开发体验提升
- ✅ 文件查找速度更快
- ✅ 代码审查更容易
- ✅ 新成员上手更快

### 长期优势
- ✅ 更好的代码内聚性
- ✅ 更低的耦合度
- ✅ 易于实施 micro-frontends（如果需要）

## ⚠️ 注意事项

1. **Git 历史**: 使用 `git mv` 保留文件历史
2. **逐步进行**: 一次移动一个 feature，测试后再继续
3. **更新文档**: 移动完成后更新 README 的项目结构说明
4. **团队通知**: 如果是团队项目，提前通知其他成员

## 🎯 总结

**总工作量**: 约 1.5-2 小时  
**难度**: 中等（主要是批量替换导入路径）  
**风险**: 低（构建失败会立即发现）  
**收益**: 高（长期可维护性大幅提升）

**建议**: 选择一个完整的时间块（如下午）进行，一气呵成完成所有步骤。

---

**创建时间**: 2025-11-15  
**维护者**: Llamacto Team
