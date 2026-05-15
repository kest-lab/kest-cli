import { ProjectDetailPage } from '@/components/features/project/project-detail-page';

interface ProjectDetailRoutePageProps {
  params: Promise<{
    projectId: string;
  }>;
}

// 项目概览入口。
// 作用：展示项目详情、工作区入口和 CLI Sync 配置，让 `/project/:projectId` 保持为项目级入口。
export default async function ProjectDetailRoutePage({ params }: ProjectDetailRoutePageProps) {
  const { projectId } = await params;
  return <ProjectDetailPage projectId={projectId} />;
}
