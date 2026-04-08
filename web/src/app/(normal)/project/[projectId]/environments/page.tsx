import { notFound } from 'next/navigation';
import { EnvironmentManagementPage } from '@/components/features/project/environment-management-page';
import { ProjectWorkspacePage } from '@/components/features/project/project-workspace-page';

interface ProjectEnvironmentsPageProps {
  params: Promise<{
    projectId: string;
  }>;
  searchParams: Promise<{
    item?: string;
    mode?: string;
  }>;
}

// 项目环境管理页面入口。
// 作用：默认挂载新的 environments 工作区，并通过 `?mode=manage` 保留旧 CRUD 入口。
export default async function ProjectEnvironmentsPage({
  params,
  searchParams,
}: ProjectEnvironmentsPageProps) {
  const { projectId } = await params;
  const { item, mode } = await searchParams;
  const numericProjectId = Number(projectId);

  // 非法项目 ID 直接返回 404，避免把错误参数继续传进受保护页面。
  if (!Number.isInteger(numericProjectId) || numericProjectId <= 0) {
    notFound();
  }

  if (mode === 'manage') {
    return <EnvironmentManagementPage projectId={numericProjectId} />;
  }

  const selectedItemId = Number(item);

  return (
    <ProjectWorkspacePage
      projectId={numericProjectId}
      module="environments"
      selectedItemId={Number.isInteger(selectedItemId) && selectedItemId > 0 ? selectedItemId : null}
    />
  );
}
