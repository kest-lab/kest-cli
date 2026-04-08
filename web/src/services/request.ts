import request from '@/http';
import type {
  CreateRequestRequest,
  ProjectRequest,
  RequestListResponse,
  UpdateRequestRequest,
} from '@/types/request';

const normalizePayload = <T extends object>(payload: T) =>
  Object.fromEntries(
    Object.entries(payload as Record<string, unknown>).filter(([, value]) => value !== undefined)
  ) as T;

export const requestService = {
  list: (projectId: number | string, collectionId: number | string) =>
    request.get<RequestListResponse>(`/projects/${projectId}/collections/${collectionId}/requests`),

  getById: (
    projectId: number | string,
    collectionId: number | string,
    requestId: number | string
  ) =>
    request.get<ProjectRequest>(
      `/projects/${projectId}/collections/${collectionId}/requests/${requestId}`
    ),

  create: (
    projectId: number | string,
    collectionId: number | string,
    data: CreateRequestRequest
  ) =>
    request.post<ProjectRequest>(
      `/projects/${projectId}/collections/${collectionId}/requests`,
      normalizePayload(data)
    ),

  update: (
    projectId: number | string,
    collectionId: number | string,
    requestId: number | string,
    data: UpdateRequestRequest
  ) =>
    request.put<ProjectRequest>(
      `/projects/${projectId}/collections/${collectionId}/requests/${requestId}`,
      normalizePayload(data)
    ),

  delete: (
    projectId: number | string,
    collectionId: number | string,
    requestId: number | string
  ) =>
    request.delete<void>(
      `/projects/${projectId}/collections/${collectionId}/requests/${requestId}`
    ),
};

export type RequestService = typeof requestService;
