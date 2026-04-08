import { notFound } from 'next/navigation';
import { ProjectWorkspacePage } from '@/components/features/project/project-workspace-page';

interface ProjectFlowsPageProps {
  params: Promise<{
    projectId: string;
  }>;
}

// 项目 flows 工作区入口。
// 作用：承载新的 flows 模块占位页，保证工作区一级导航包含完整的信息架构。
export default async function ProjectFlowsPage({
  params,
}: ProjectFlowsPageProps) {
  const { projectId } = await params;
  const numericProjectId = Number(projectId);

  if (!Number.isInteger(numericProjectId) || numericProjectId <= 0) {
    notFound();
  }

  return <ProjectWorkspacePage projectId={numericProjectId} module="flows" />;
}
