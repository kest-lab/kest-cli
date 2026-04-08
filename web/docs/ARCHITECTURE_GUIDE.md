# 架构优化指南

本文档评估 Llamacto Web Scaffold 的整体架构，并提供改进建议。

## 🏗️ 当前架构评分: 8.5/10

### ✅ 优点

1. **清晰的分层架构** - 表现层、业务层、数据层分离
2. **模块化设计** - 组件、服务、状态管理独立
3. **类型安全** - 全面的 TypeScript 支持
4. **现代化技术栈** - Next.js 16 + React 19

### ⚠️ 改进空间

---

## 📁 目录结构优化

### 当前结构 (7/10)

```
src/
├── app/                     ✅ Next.js App Router
│   ├── (auth)/             ✅ 路由分组
│   ├── (site)/             ✅ 
│   └── console/            ✅ 
├── components/             ⚠️  需要优化
│   ├── ui/                ✅ 设计系统
│   ├── auth/              ✅ 功能分组
│   ├── login-form.tsx     ❌ 应该在 auth/ 下
│   ├── register-form.tsx  ❌ 应该在 auth/ 下
│   ├── chart-*.tsx        ❌ 应该在 charts/ 下
│   ├── data-table.tsx     ❌ 应该在 tables/ 下
│   └── ...                ⚠️  其他散落的组件
├── hooks/                  ✅ 
├── lib/                    ✅ 
│   ├── services/          ✅ 业务逻辑
│   └── types/             ✅ 类型定义
├── providers/              ✅ 
├── store/                  ✅ 状态管理
└── utils/                  ✅ 工具函数
```

### 建议结构 (9/10)

```
src/
├── app/                          # Next.js 路由
│   ├── (auth)/                  # 认证路由组
│   │   ├── login/
│   │   └── register/
│   ├── (marketing)/             # 新增：营销页面组
│   │   └── page.tsx             # 首页
│   └── (dashboard)/             # 重命名：控制台组
│       └── console/
│           ├── analytics/
│           └── settings/
│
├── components/
│   ├── ui/                      # 基础 UI 组件
│   │   ├── button.tsx
│   │   ├── card.tsx
│   │   └── ...
│   │
│   ├── features/                # 新增：功能组件
│   │   ├── auth/               # 认证相关
│   │   │   ├── login-form.tsx
│   │   │   ├── register-form.tsx
│   │   │   └── protected-route.tsx
│   │   │
│   │   ├── charts/             # 新增：图表组件
│   │   │   ├── area-chart.tsx
│   │   │   └── bar-chart.tsx
│   │   │
│   │   ├── tables/             # 新增：表格组件
│   │   │   └── data-table.tsx
│   │   │
│   │   └── navigation/         # 新增：导航组件
│   │       ├── app-sidebar.tsx
│   │       ├── nav-main.tsx
│   │       └── nav-user.tsx
│   │
│   └── layout/                  # 新增：布局组件
│       ├── site-header.tsx
│       └── site-footer.tsx
│
├── lib/
│   ├── api/                     # 新增：API 客户端
│   │   ├── client.ts           # 统一的 API 客户端
│   │   └── endpoints.ts        # API 端点定义
│   │
│   ├── services/                # 业务服务
│   │   ├── auth.service.ts
│   │   ├── user.service.ts
│   │   └── index.ts
│   │
│   ├── types/                   # 类型定义
│   │   ├── auth.types.ts
│   │   ├── user.types.ts
│   │   └── api.types.ts
│   │
│   ├── utils/                   # 工具函数
│   │   ├── cn.ts               # className 合并
│   │   ├── format.ts           # 格式化工具
│   │   └── validation.ts       # 验证工具
│   │
│   └── constants/               # 常量定义
│       ├── routes.ts
│       └── config.ts
│
├── hooks/                       # 自定义 Hooks
│   ├── use-auth.ts
│   ├── use-media-query.ts
│   └── index.ts
│
├── store/                       # 状态管理
│   ├── auth.store.ts
│   ├── ui.store.ts
│   └── index.ts
│
└── providers/                   # Context Providers
    ├── auth-provider.tsx
    ├── theme-provider.tsx
    └── query-provider.tsx
```

---

## 🔧 代码组织原则

### 1. 按功能分组 (Feature-Based)

**当前问题**: 组件散落在顶层

**建议**: 按功能模块组织

```
components/
├── features/
│   ├── auth/              # 所有认证相关组件
│   │   ├── index.ts      # 统一导出
│   │   ├── login-form/
│   │   │   ├── index.tsx
│   │   │   ├── login-form.test.tsx
│   │   │   └── login-form.styles.ts
│   │   └── register-form/
│   │       └── ...
│   │
│   └── dashboard/         # 所有仪表板组件
│       └── ...
```

**优点**:
- ✅ 更容易找到相关代码
- ✅ 更好的代码内聚性
- ✅ 便于团队协作
- ✅ 易于删除/移动整个功能

---

### 2. 服务层模式 (Service Layer Pattern)

**当前问题**: `auth.ts` 混合了多种职责

**建议**: 清晰的服务层架构

```typescript
// lib/services/base.service.ts
export abstract class BaseService {
  protected apiClient: ApiClient;
  
  constructor() {
    this.apiClient = createApiClient();
  }
  
  protected handleError(error: unknown): never {
    // 统一错误处理
  }
}

// lib/services/auth.service.ts
export class AuthService extends BaseService {
  async login(credentials: LoginCredentials): Promise<User> {
    return this.apiClient.post('/auth/login', credentials);
  }
  
  async logout(): Promise<void> {
    return this.apiClient.post('/auth/logout');
  }
}

// lib/services/index.ts
export const authService = new AuthService();
export const userService = new UserService();
```

**优点**:
- ✅ 单一职责原则
- ✅ 易于测试
- ✅ 代码复用
- ✅ 统一的错误处理

---

### 3. API 客户端抽象

**当前问题**: 直接使用 axios，缺少统一抽象

**建议**: 创建 API 客户端层

```typescript
// lib/api/client.ts
import type { ApiRequestConfig, ApiResponse } from './types';

class ApiClient {
  private baseURL: string;
  private defaultHeaders: Record<string, string>;
  
  constructor(config: ApiClientConfig) {
    this.baseURL = config.baseURL || '';
    this.defaultHeaders = config.headers || {};
  }
  
  async request<T>(config: ApiRequestConfig): Promise<ApiResponse<T>> {
    // 统一请求处理
    // - 添加认证头
    // - 错误处理
    // - 重试逻辑
    // - 日志记录
  }
  
  get<T>(url: string, config?: ApiRequestConfig) {
    return this.request<T>({ ...config, method: 'GET', url });
  }
  
  post<T>(url: string, data?: unknown, config?: ApiRequestConfig) {
    return this.request<T>({ ...config, method: 'POST', url, data });
  }
  
  // ... put, delete, patch
}

// lib/api/endpoints.ts
export const API_ENDPOINTS = {
  auth: {
    login: '/auth/login',
    logout: '/auth/logout',
    register: '/auth/register',
  },
  users: {
    list: '/users',
    detail: (id: string) => `/users/${id}`,
  },
} as const;

// 使用
import { apiClient } from '@/lib/api/client';
import { API_ENDPOINTS } from '@/lib/api/endpoints';

await apiClient.post(API_ENDPOINTS.auth.login, credentials);
```

---

### 4. 类型定义组织

**当前问题**: 类型定义散落各处

**建议**: 集中管理类型

```typescript
// lib/types/api.types.ts
export interface ApiResponse<T = unknown> {
  data: T;
  message?: string;
  code: number;
}

export interface PaginatedResponse<T> {
  items: T[];
  total: number;
  page: number;
  pageSize: number;
}

// lib/types/auth.types.ts
export interface User {
  id: string;
  email: string;
  name: string;
  role: UserRole;
}

export interface LoginCredentials {
  email: string;
  password: string;
}

// lib/types/index.ts
export * from './api.types';
export * from './auth.types';
export * from './user.types';
```

---

## 🎯 架构模式建议

### 1. Repository Pattern (仓储模式)

**适用场景**: 复杂的数据访问逻辑

```typescript
// lib/repositories/user.repository.ts
export class UserRepository {
  constructor(private apiClient: ApiClient) {}
  
  async findAll(filters?: UserFilters): Promise<User[]> {
    const response = await this.apiClient.get('/users', { params: filters });
    return response.data;
  }
  
  async findById(id: string): Promise<User | null> {
    try {
      const response = await this.apiClient.get(`/users/${id}`);
      return response.data;
    } catch (error) {
      if (error.status === 404) return null;
      throw error;
    }
  }
  
  async create(data: CreateUserDto): Promise<User> {
    const response = await this.apiClient.post('/users', data);
    return response.data;
  }
}

// 使用
const userRepo = new UserRepository(apiClient);
const user = await userRepo.findById('123');
```

---

### 2. Custom Hooks Pattern

**统一业务逻辑的 Hooks**

```typescript
// hooks/use-auth.ts
export function useAuth() {
  const user = useAuthStore((state) => state.user);
  const login = useAuthStore((state) => state.login);
  const logout = useAuthStore((state) => state.logout);
  
  return {
    user,
    isAuthenticated: !!user,
    isAdmin: user?.role === 'admin',
    login,
    logout,
  };
}

// hooks/use-api.ts
export function useApi<T>(
  fetcher: () => Promise<T>,
  options?: UseApiOptions
) {
  return useQuery({
    queryKey: options?.queryKey || [],
    queryFn: fetcher,
    ...options,
  });
}

// 使用
const { user, isAuthenticated, login } = useAuth();
const { data: users, isLoading } = useApi(() => userService.getAll());
```

---

### 3. Compound Components Pattern

**复杂组件的组合模式**

```typescript
// components/ui/data-table/index.tsx
interface DataTableProps<T> {
  data: T[];
  children: React.ReactNode;
}

export function DataTable<T>({ data, children }: DataTableProps<T>) {
  return (
    <DataTableContext.Provider value={{ data }}>
      <div className="data-table">{children}</div>
    </DataTableContext.Provider>
  );
}

DataTable.Header = DataTableHeader;
DataTable.Body = DataTableBody;
DataTable.Row = DataTableRow;
DataTable.Cell = DataTableCell;
DataTable.Pagination = DataTablePagination;

// 使用
<DataTable data={users}>
  <DataTable.Header>
    <DataTable.Cell>Name</DataTable.Cell>
    <DataTable.Cell>Email</DataTable.Cell>
  </DataTable.Header>
  <DataTable.Body />
  <DataTable.Pagination />
</DataTable>
```

---

## 🚦 状态管理优化

### 当前状态: Zustand (7/10)

**优点**:
- ✅ 简单易用
- ✅ 性能好
- ✅ TypeScript 支持

**建议增强**:

```typescript
// store/slices/auth.slice.ts
export interface AuthSlice {
  user: User | null;
  token: string | null;
  
  // Actions
  setUser: (user: User) => void;
  clearUser: () => void;
  
  // Async actions
  login: (credentials: LoginCredentials) => Promise<void>;
  logout: () => Promise<void>;
}

export const createAuthSlice: StateCreator<AuthSlice> = (set, get) => ({
  user: null,
  token: null,
  
  setUser: (user) => set({ user }),
  clearUser: () => set({ user: null, token: null }),
  
  login: async (credentials) => {
    const { user, token } = await authService.login(credentials);
    set({ user, token });
  },
  
  logout: async () => {
    await authService.logout();
    set({ user: null, token: null });
  },
});

// store/index.ts
export const useStore = create<AuthSlice & UiSlice>()(
  devtools(
    persist(
      (...args) => ({
        ...createAuthSlice(...args),
        ...createUiSlice(...args),
      }),
      { name: 'app-storage' }
    )
  )
);
```

---

## 📝 代码规范建议

### 1. 导入顺序

```typescript
// ✅ 推荐的导入顺序
// 1. React 相关
import { useState, useEffect } from 'react';
import type { FC } from 'react';

// 2. 第三方库
import { motion } from 'framer-motion';
import { useQuery } from '@tanstack/react-query';

// 3. 内部模块 - 按字母顺序
import { Button } from '@/components/ui/button';
import { useAuth } from '@/hooks/use-auth';
import { cn } from '@/lib/utils';
import type { User } from '@/lib/types';

// 4. 相对导入
import { LoginFormSchema } from './schema';
import styles from './login-form.module.css';
```

### 2. 文件命名规范

```
# 组件文件
login-form.tsx          ✅ kebab-case
LoginForm.tsx           ❌ PascalCase (不推荐)

# 工具函数
format-date.ts          ✅ kebab-case
formatDate.ts           ❌ camelCase (不推荐)

# 类型文件
auth.types.ts           ✅ 描述性命名
types.ts                ❌ 太通用

# 测试文件
login-form.test.tsx     ✅ 与源文件同名 + .test
login-form.spec.tsx     ✅ 也可以 .spec
```

### 3. 组件结构规范

```typescript
// ✅ 推荐的组件结构
import type { FC } from 'react';

// 1. 类型定义
interface LoginFormProps {
  onSuccess?: () => void;
  className?: string;
}

// 2. 常量定义
const DEFAULT_VALUES = {
  email: '',
  password: '',
};

// 3. 主组件
export const LoginForm: FC<LoginFormProps> = ({ 
  onSuccess,
  className 
}) => {
  // Hooks
  const { login } = useAuth();
  const form = useForm();
  
  // Handlers
  const handleSubmit = async (data: LoginFormData) => {
    await login(data);
    onSuccess?.();
  };
  
  // Render
  return (
    <form onSubmit={form.handleSubmit(handleSubmit)}>
      {/* JSX */}
    </form>
  );
};

// 4. 子组件（如果需要）
const LoginFormFooter = () => {
  return <div>Footer</div>;
};
```

---

## 📊 质量指标

### 代码质量目标

| 指标 | 当前 | 目标 |
|------|------|------|
| TypeScript 覆盖率 | 95% | 100% |
| 组件复用率 | 70% | 85% |
| 单一职责组件 | 75% | 90% |
| 文档覆盖率 | 40% | 80% |
| 测试覆盖率 | 0% | 70% |

---

## 🔄 迁移计划

### Phase 1: 目录重构 (2-3天)
1. 创建新的目录结构
2. 移动文件到对应位置
3. 更新所有导入路径
4. 运行测试确保无误

### Phase 2: 服务层重构 (3-5天)
1. 创建 BaseService
2. 重构现有服务
3. 添加统一错误处理
4. 编写单元测试

### Phase 3: 组件优化 (5-7天)
1. 按功能重组组件
2. 抽取公共逻辑到 hooks
3. 优化组件复用性
4. 添加组件文档

### Phase 4: 类型系统增强 (2-3天)
1. 集中管理类型定义
2. 添加工具类型
3. 完善 API 类型
4. 类型安全检查

---

## 📚 最佳实践参考

- [Next.js Project Structure](https://nextjs.org/docs/getting-started/project-structure)
- [React Clean Architecture](https://github.com/eduardomoroni/react-clean-architecture)
- [TypeScript Best Practices](https://typescript-eslint.io/rules/)
- [Feature-Sliced Design](https://feature-sliced.design/)

---

**最后更新**: 2025-11-15
**维护者**: Llamacto Team
