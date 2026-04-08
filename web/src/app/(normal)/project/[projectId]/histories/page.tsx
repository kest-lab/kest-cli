import { notFound } from 'next/navigation';
import { ProjectWorkspacePage } from '@/components/features/project/project-workspace-page';

interface ProjectHistoriesPageProps {
  params: Promise<{
    projectId: string;
  }>;
}

// 项目 histories 工作区入口。
// 作用：承载新的 histories 模块占位页，确保双层导航结构完整可用。
export default async function ProjectHistoriesPage({
  params,
}: ProjectHistoriesPageProps) {
  const { projectId } = await params;
  const numericProjectId = Number(projectId);

  if (!Number.isInteger(numericProjectId) || numericProjectId <= 0) {
    notFound();
  }

  return <ProjectWorkspacePage projectId={numericProjectId} module="histories" />;
}
