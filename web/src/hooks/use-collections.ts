'use client';

import { useMutation, useQueryClient } from '@tanstack/react-query';
import { toast } from 'sonner';
import { collectionService } from '@/services/collection';
import type { UpdateCollectionRequest } from '@/services/collection';

// Collections 域的 React Query key。
// 作用：统一 collection 列表和详情的缓存命名，便于后续继续接列表/详情接口。
export const collectionKeys = {
  all: ['collections'] as const,
  project: (projectId: number | string) => [...collectionKeys.all, 'project', projectId] as const,
  list: (projectId: number | string) => [...collectionKeys.project(projectId), 'list'] as const,
  detail: (projectId: number | string, collectionId: number | string) =>
    [...collectionKeys.project(projectId), 'detail', collectionId] as const,
};

// 删除 collection mutation。
// 作用：调用后端删除接口，并清理当前项目下 collection 相关缓存。
export function useDeleteCollection(projectId: number | string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (collectionId: number | string) =>
      collectionService.delete(projectId, collectionId),
    onSuccess: (_, collectionId) => {
      queryClient.invalidateQueries({ queryKey: collectionKeys.project(projectId) });
      queryClient.removeQueries({
        queryKey: collectionKeys.detail(projectId, collectionId),
      });
      toast.success('Collection deleted');
    },
  });
}

// 更新 collection mutation。
// 作用：提交名称等字段的更新，并刷新当前项目下的 collection 缓存。
export function useUpdateCollection(projectId: number | string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      collectionId,
      data,
    }: {
      collectionId: number | string;
      data: UpdateCollectionRequest;
    }) => collectionService.update(projectId, collectionId, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: collectionKeys.project(projectId) });
      toast.success(`Renamed collection to "${variables.data.name ?? 'Untitled'}"`);
    },
  });
}
