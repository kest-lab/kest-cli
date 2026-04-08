# 性能优化指南

本文档提供 Llamacto Web Scaffold 的性能优化最佳实践和建议。

## 📊 当前状态评估

### 构建大小
- **.next 文件夹**: ~62MB
- **文件数量**: 75+ TypeScript/TSX 文件
- **依赖包**: 40+ 生产依赖

---

## 🎯 高优先级优化 (立即实施)

### 1. 组件懒加载

**问题**: 所有组件都在初始加载时打包，增加首屏加载时间。

**解决方案**: 使用 `dynamic()` 实现代码分割

```tsx
// 示例：懒加载重组件
import dynamic from 'next/dynamic';

// 图表组件（通常较大）
const ChartAreaInteractive = dynamic(
  () => import('@/components/chart-area-interactive'),
  {
    loading: () => <div className="h-96 animate-pulse bg-muted" />,
    ssr: false, // 图表组件通常不需要 SSR
  }
);

// 数据表格（包含大量依赖）
const DataTable = dynamic(() => import('@/components/data-table'), {
  loading: () => <div>Loading table...</div>,
});

// React Query Devtools（仅开发环境）
const ReactQueryDevtools = dynamic(
  () =>
    import('@tanstack/react-query-devtools').then(
      (mod) => mod.ReactQueryDevtools
    ),
  { ssr: false }
);
```

**实施位置**:
- ✅ `src/components/chart-area-interactive.tsx`
- ✅ `src/components/data-table.tsx`
- ✅ `src/providers/query-provider.tsx` (devtools)
- ✅ 所有模态框和抽屉组件

**预期收益**: 减少 30-40% 的初始 bundle 大小

---

### 2. 图标库优化

**问题**: 
- `lucide-react`: ~600KB
- `@tabler/icons-react`: ~1MB
- 只使用了其中很少一部分图标

**解决方案**: 按需导入图标

```tsx
// ❌ 错误：导入整个图标库
import * as Icons from 'lucide-react';

// ✅ 正确：只导入需要的图标
import { Home, Settings, User } from 'lucide-react';
```

**实施清单**:
- [ ] 审查所有图标导入
- [ ] 创建自定义图标导出文件
- [ ] 使用 tree-shaking

```tsx
// src/components/ui/icons.ts
export {
  Home,
  Settings,
  User,
  LogOut,
  Menu,
  // 只导出实际使用的图标
} from 'lucide-react';
```

**预期收益**: 减少 200-400KB bundle 大小

---

### 3. 依赖包优化

**当前问题**:
```json
{
  "axios": "^1.13.2",           // 28KB (gzipped)
  "@tanstack/react-query": "^5.90.9",  // 46KB (gzipped)
}
```

**优化建议**:

#### A. 考虑轻量级替代方案
```bash
# axios (28KB) → native fetch (0KB)
# 如果不需要拦截器和复杂配置，使用原生 fetch

# date-fns → day.js
# 如果需要日期处理，优先使用 day.js (2KB vs 33KB)
```

#### B. 移除未使用的依赖
```bash
# 检查未使用的依赖
pnpm dlx depcheck

# 分析 bundle 大小
pnpm dlx @next/bundle-analyzer
```

**预期收益**: 减少 10-20% bundle 大小

---

### 4. 字体优化

**当前状态**: 使用 Google Fonts (Inter)

**优化方案**:
```tsx
// src/app/layout.tsx
import { Inter } from "next/font/google";

const inter = Inter({ 
  subsets: ["latin"],
  display: 'swap', // ✅ 已优化
  preload: true,   // ✅ 已优化
  // 添加以下优化
  variable: '--font-inter',
  fallback: ['system-ui', 'arial'],
  adjustFontFallback: true, // 减少布局偏移
});
```

**预期收益**: 改善 CLS (Cumulative Layout Shift)

---

## 🚀 中优先级优化 (短期实施)

### 5. 数据获取优化

**实施 ISR (Incremental Static Regeneration)**

```tsx
// src/app/page.tsx
export const revalidate = 3600; // 1小时重新验证

// src/app/console/page.tsx
export const revalidate = 60; // 1分钟重新验证（动态数据）
```

**实施 Streaming SSR**

```tsx
import { Suspense } from 'react';

export default function Page() {
  return (
    <>
      <Header /> {/* 立即显示 */}
      <Suspense fallback={<LoadingSkeleton />}>
        <DataComponent /> {/* 流式传输 */}
      </Suspense>
    </>
  );
}
```

---

### 6. 图片优化

**创建图片组件规范**:

```tsx
// src/components/ui/optimized-image.tsx
import Image from 'next/image';

interface OptimizedImageProps {
  src: string;
  alt: string;
  width: number;
  height: number;
  priority?: boolean;
}

export function OptimizedImage({
  src,
  alt,
  width,
  height,
  priority = false,
}: OptimizedImageProps) {
  return (
    <Image
      src={src}
      alt={alt}
      width={width}
      height={height}
      priority={priority}
      sizes="(max-width: 768px) 100vw, (max-width: 1200px) 50vw, 33vw"
      quality={85}
      placeholder="blur"
      blurDataURL="data:image/svg+xml;base64,..." // 添加 blur placeholder
    />
  );
}
```

---

### 7. CSS 优化

**当前状态**: Tailwind CSS v4

**优化建议**:

```css
/* src/app/globals.css */

/* ✅ 使用 CSS containment */
.card {
  contain: layout style paint;
}

/* ✅ 使用 content-visibility */
.lazy-section {
  content-visibility: auto;
  contain-intrinsic-size: 0 500px;
}

/* ✅ 优化动画性能 */
@media (prefers-reduced-motion: reduce) {
  * {
    animation-duration: 0.01ms !important;
    animation-iteration-count: 1 !important;
    transition-duration: 0.01ms !important;
  }
}
```

---

## 🎨 低优先级优化 (长期优化)

### 8. 构建时优化

**添加构建分析**:

```json
// package.json
{
  "scripts": {
    "analyze": "ANALYZE=true pnpm build",
    "build:profile": "next build --profile"
  }
}
```

**使用 Bundle Analyzer**:

```bash
pnpm add -D @next/bundle-analyzer

# next.config.ts
const withBundleAnalyzer = require('@next/bundle-analyzer')({
  enabled: process.env.ANALYZE === 'true',
});

module.exports = withBundleAnalyzer(nextConfig);
```

---

### 9. 运行时优化

**实施虚拟滚动**:

```tsx
// 对于长列表，使用 react-window 或 @tanstack/react-virtual
import { useVirtualizer } from '@tanstack/react-virtual';

function LargeList({ items }) {
  const parentRef = React.useRef();
  
  const virtualizer = useVirtualizer({
    count: items.length,
    getScrollElement: () => parentRef.current,
    estimateSize: () => 50,
  });

  return (
    <div ref={parentRef} style={{ height: '400px', overflow: 'auto' }}>
      {virtualizer.getVirtualItems().map((virtualItem) => (
        <div key={virtualItem.index}>
          {items[virtualItem.index]}
        </div>
      ))}
    </div>
  );
}
```

---

### 10. 监控和分析

**实施 Web Vitals 监控**:

```tsx
// src/app/layout.tsx
import { SpeedInsights } from '@vercel/speed-insights/next';
import { Analytics } from '@vercel/analytics/react';

export default function RootLayout({ children }) {
  return (
    <html>
      <body>
        {children}
        <SpeedInsights />
        <Analytics />
      </body>
    </html>
  );
}
```

**添加性能监控**:

```tsx
// src/lib/performance.ts
export function reportWebVitals(metric: any) {
  // 发送到分析服务
  if (process.env.NODE_ENV === 'production') {
    // 示例：发送到 Google Analytics
    window.gtag?.('event', metric.name, {
      value: Math.round(metric.value),
      metric_id: metric.id,
      metric_value: metric.value,
      metric_delta: metric.delta,
    });
  }
}
```

---

## 📈 预期性能指标

### 优化前
- **FCP (First Contentful Paint)**: ~2.5s
- **LCP (Largest Contentful Paint)**: ~3.5s
- **TTI (Time to Interactive)**: ~4.0s
- **Bundle Size**: ~800KB (gzipped)

### 优化后目标
- **FCP**: <1.5s ⬇️ 40%
- **LCP**: <2.0s ⬇️ 43%
- **TTI**: <2.5s ⬇️ 38%
- **Bundle Size**: ~400KB ⬇️ 50%

---

## 🔧 实施检查清单

### Phase 1 (1-2天)
- [ ] 启用 Next.js 配置优化 (已完成 ✅)
- [ ] 实施组件懒加载
- [ ] 优化图标导入
- [ ] 添加 loading states

### Phase 2 (3-5天)
- [ ] 实施 ISR
- [ ] 添加 Suspense boundaries
- [ ] 优化图片加载
- [ ] 清理未使用依赖

### Phase 3 (1-2周)
- [ ] 添加 Bundle Analyzer
- [ ] 实施虚拟滚动
- [ ] 添加性能监控
- [ ] 进行性能审计

---

## 📚 参考资源

- [Next.js Performance Best Practices](https://nextjs.org/docs/app/building-your-application/optimizing)
- [Web Vitals](https://web.dev/vitals/)
- [React Performance Optimization](https://react.dev/learn/render-and-commit)
- [Webpack Bundle Analyzer](https://github.com/webpack-contrib/webpack-bundle-analyzer)

---

**最后更新**: 2025-11-15
**维护者**: Llamacto Team
