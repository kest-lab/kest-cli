'use client';

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { toast } from 'sonner';
import { projectService } from '@/services/project';
import type {
  CreateProjectRequest,
  ProjectListParams,
  ProjectStats,
  UpdateProjectRequest,
} from '@/types/project';

// 项目域的 React Query key。
// 作用：统一项目列表、详情、统计的缓存命名，方便后续失效与刷新。
export const projectKeys = {
  all: ['projects'] as const,
  lists: () => [...projectKeys.all, 'list'] as const,
  list: (params: ProjectListParams) => [...projectKeys.lists(), params] as const,
  details: () => [...projectKeys.all, 'detail'] as const,
  detail: (id: number | string) => [...projectKeys.details(), id] as const,
  stats: () => [...projectKeys.all, 'stats'] as const,
  projectStats: (id: number | string) => [...projectKeys.stats(), id] as const,
};

// 项目列表查询。
// 作用：拉取当前登录用户可见的项目分页列表，并在翻页时保留上一页数据减少闪烁。
export function useProjects(params: ProjectListParams = {}) {
  return useQuery({
    queryKey: projectKeys.list(params),
    queryFn: () => projectService.list(params),
    placeholderData: (previousData) => previousData,
  });
}

// 项目详情查询。
// 作用：按项目 ID 获取详情数据，供右侧详情面板或其他页面复用。
export function useProject(id?: number | string) {
  return useQuery({
    queryKey: projectKeys.detail(id ?? 'unknown'),
    queryFn: () => projectService.getById(id as number | string),
    enabled: id !== undefined && id !== null && id !== '',
  });
}

// 项目统计查询。
// 作用：读取 `/projects/:id/stats`，展示 API specs、flows、members 等聚合信息。
export function useProjectStats(id?: number | string) {
  return useQuery<ProjectStats>({
    queryKey: projectKeys.projectStats(id ?? 'unknown'),
    queryFn: () => projectService.getStats(id as number | string),
    enabled: id !== undefined && id !== null && id !== '',
  });
}

// 创建项目 mutation。
// 作用：调用创建接口后刷新列表，并把新项目详情提前写入缓存。
export function useCreateProject() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateProjectRequest) => projectService.create(data),
    onSuccess: (project) => {
      queryClient.invalidateQueries({ queryKey: projectKeys.lists() });
      queryClient.setQueryData(projectKeys.detail(project.id), project);
      toast.success(`Created project "${project.name}"`);
    },
  });
}

// 更新项目 mutation。
// 作用：提交项目编辑后同步刷新列表、详情和统计缓存。
export function useUpdateProject() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: number | string; data: UpdateProjectRequest }) =>
      projectService.update(id, data),
    onSuccess: (project) => {
      queryClient.invalidateQueries({ queryKey: projectKeys.lists() });
      queryClient.setQueryData(projectKeys.detail(project.id), project);
      queryClient.invalidateQueries({ queryKey: projectKeys.projectStats(project.id) });
      toast.success(`Updated project "${project.name}"`);
    },
  });
}

// 删除项目 mutation。
// 作用：项目删除成功后移除对应详情/统计缓存，并触发列表刷新。
export function useDeleteProject() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: number | string) => projectService.delete(id),
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: projectKeys.lists() });
      queryClient.removeQueries({ queryKey: projectKeys.detail(id) });
      queryClient.removeQueries({ queryKey: projectKeys.projectStats(id) });
      toast.success('Project deleted');
    },
  });
}
