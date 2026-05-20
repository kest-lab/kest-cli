import { ProjectWorkspacePage } from '@/components/features/project/project-workspace-page';

interface WorkspaceHistoriesPageProps {
  params: Promise<{
    workspaceId: string;
  }>;
  searchParams: Promise<{
    item?: string;
    entityType?: string;
    view?: string;
    run?: string;
    sourceType?: string;
    status?: string;
  }>;
}

// 项目 histories 工作区入口。
// 作用：挂载项目历史工作区，并通过 `?item=` 支持选中具体记录。
export default async function WorkspaceHistoriesPage({
  params,
  searchParams,
}: WorkspaceHistoriesPageProps) {
  const { workspaceId } = await params;
  const { item, entityType, run, view, sourceType, status } = await searchParams;
  const normalizedView = view?.trim() === 'runs' ? 'runs' : 'activity';
  const selectedItemId = normalizedView === 'runs' ? (run?.trim() ? run : null) : item?.trim() ? item : null;
  const initialHistoryEntityType = entityType?.trim() ? entityType : null;
  const initialRunSourceType = sourceType?.trim() ? sourceType : null;
  const initialRunStatus = status?.trim() ? status : null;

  return (
    <ProjectWorkspacePage
      projectId={workspaceId}
      module="histories"
      selectedItemId={selectedItemId}
      initialHistoryEntityType={initialHistoryEntityType}
      initialHistoryView={normalizedView}
      initialRunSourceType={initialRunSourceType as
        | 'request'
        | 'collection'
        | 'test_case'
        | 'flow'
        | null}
      initialRunStatus={initialRunStatus}
    />
  );
}
