import { notFound, redirect } from 'next/navigation';
import { buildProjectApiSpecsRoute } from '@/constants/routes';
import { ProjectDetailPage } from '@/components/features/project/project-detail-page';

interface ProjectDetailRoutePageProps {
  params: Promise<{
    projectId: string;
  }>;
  searchParams: Promise<{
    mode?: string;
  }>;
}

// 项目工作区入口。
// 作用：
// 1. 默认把 `/project/:projectId` 收敛到工作区默认模块 `/api-specs`
// 2. 通过 `?mode=details` 保留旧的项目详情页出口，避免已接好的页面能力直接丢失
export default async function ProjectDetailRoutePage({
  params,
  searchParams,
}: ProjectDetailRoutePageProps) {
  const { projectId } = await params;
  const { mode } = await searchParams;
  const numericProjectId = Number(projectId);

  if (!Number.isInteger(numericProjectId) || numericProjectId <= 0) {
    notFound();
  }

  if (mode === 'details') {
    return <ProjectDetailPage projectId={numericProjectId} />;
  }

  redirect(buildProjectApiSpecsRoute(numericProjectId));
}
