import { ConsoleShell } from '@/components/features/console/console-shell';

// 控制台路由组布局。
// 作用：让 `/console/*` 页面统一复用后台壳层，不在每个 page 里重复写导航和头部。
export default function ConsoleLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return <ConsoleShell>{children}</ConsoleShell>;
}
