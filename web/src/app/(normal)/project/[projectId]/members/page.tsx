import { notFound } from 'next/navigation';
import { ProjectMemberManagementPage } from '@/components/features/project/project-member-management-page';

interface ProjectMembersPageProps {
  params: Promise<{
    projectId: string;
  }>;
}

// 项目成员管理页面入口。
// 作用：读取动态项目 ID，并挂载项目成员管理界面。
export default async function ProjectMembersPage({
  params,
}: ProjectMembersPageProps) {
  const { projectId } = await params;
  const numericProjectId = Number(projectId);

  if (!Number.isInteger(numericProjectId) || numericProjectId <= 0) {
    notFound();
  }

  return <ProjectMemberManagementPage projectId={numericProjectId} />;
}
