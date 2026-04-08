import { notFound } from 'next/navigation';
import { CategoryManagementPage } from '@/components/features/project/category-management-page';
import { ProjectWorkspacePage } from '@/components/features/project/project-workspace-page';

interface ProjectCategoriesPageProps {
  params: Promise<{
    projectId: string;
  }>;
  searchParams: Promise<{
    item?: string;
    mode?: string;
  }>;
}

// 项目分类管理页面入口。
// 作用：默认挂载新的 categories 工作区，并通过 `?mode=manage` 兼容旧管理页。
export default async function ProjectCategoriesPage({
  params,
  searchParams,
}: ProjectCategoriesPageProps) {
  const { projectId } = await params;
  const { item, mode } = await searchParams;
  const numericProjectId = Number(projectId);

  // 非法项目 ID 直接返回 404，避免把错误参数继续传入受保护页面。
  if (!Number.isInteger(numericProjectId) || numericProjectId <= 0) {
    notFound();
  }

  if (mode === 'manage') {
    return <CategoryManagementPage projectId={numericProjectId} />;
  }

  const selectedItemId = Number(item);

  return (
    <ProjectWorkspacePage
      projectId={numericProjectId}
      module="categories"
      selectedItemId={Number.isInteger(selectedItemId) && selectedItemId > 0 ? selectedItemId : null}
    />
  );
}
