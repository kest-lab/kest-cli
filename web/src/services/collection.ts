import request from '@/http';

export interface UpdateCollectionRequest {
  name?: string;
}

// 请求体清理器。
// 作用：过滤 `undefined` 字段，避免把无意义空值提交给后端。
const normalizePayload = <T extends object>(payload: T) =>
  Object.fromEntries(
    Object.entries(payload as Record<string, unknown>).filter(([, value]) => value !== undefined)
  ) as T;

// Collections 服务层。
// 作用：集中封装项目 collections 的 HTTP 请求，先提供更新和删除接口给工作台复用。
export const collectionService = {
  update: (
    projectId: number | string,
    collectionId: number | string,
    data: UpdateCollectionRequest
  ) =>
    request.put<unknown>(
      `/projects/${projectId}/collections/${collectionId}`,
      normalizePayload(data)
    ),

  delete: (projectId: number | string, collectionId: number | string) =>
    request.delete<void>(`/projects/${projectId}/collections/${collectionId}`),
};

export type CollectionService = typeof collectionService;
