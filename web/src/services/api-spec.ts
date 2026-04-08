import request from '@/http';
import type {
  ApiSpec,
  ApiSpecExample,
  ApiSpecExamplesResponse,
  ApiSpecExportFormat,
  ApiSpecExportPayload,
  ApiSpecLanguage,
  ApiSpecListParams,
  ApiSpecListResponse,
  BatchGenDocRequest,
  BatchGenDocResponse,
  CreateApiExampleRequest,
  CreateApiSpecRequest,
  GenApiTestResponse,
  ImportApiSpecsRequest,
  ImportApiSpecsResponse,
  ProjectCategoryListResponse,
  ProjectMemberRoleResponse,
  UpdateApiSpecRequest,
} from '@/types/api-spec';

// API Specifications 服务层。
// 作用：集中封装项目级 API 规格、示例、导入导出、AI 生成、分类和成员角色查询请求。
// 额外约束：请求体会先清理 `undefined` 字段，避免把无意义空字段发给后端。
const normalizePayload = <T extends object>(payload: T) =>
  Object.fromEntries(
    Object.entries(payload as Record<string, unknown>).filter(([, value]) => value !== undefined)
  ) as T;

export const apiSpecService = {
  list: ({
    projectId,
    page = 1,
    pageSize = 20,
    version,
    method,
    tag,
    keyword,
  }: ApiSpecListParams) =>
    request.get<ApiSpecListResponse>(`/projects/${projectId}/api-specs`, {
      params: normalizePayload({
        page,
        page_size: pageSize,
        version,
        method,
        tag,
        keyword,
      }),
    }),

  getById: (projectId: number | string, specId: number | string) =>
    request.get<ApiSpec>(`/projects/${projectId}/api-specs/${specId}`),

  getFullById: (projectId: number | string, specId: number | string) =>
    request.get<ApiSpec>(`/projects/${projectId}/api-specs/${specId}/full`),

  create: (projectId: number | string, data: CreateApiSpecRequest) =>
    request.post<ApiSpec>(`/projects/${projectId}/api-specs`, normalizePayload(data)),

  update: (projectId: number | string, specId: number | string, data: UpdateApiSpecRequest) =>
    request.patch<ApiSpec>(`/projects/${projectId}/api-specs/${specId}`, normalizePayload(data)),

  delete: (projectId: number | string, specId: number | string) =>
    request.delete<void>(`/projects/${projectId}/api-specs/${specId}`),

  import: (projectId: number | string, data: ImportApiSpecsRequest) =>
    request.post<ImportApiSpecsResponse>(`/projects/${projectId}/api-specs/import`, data),

  export: (projectId: number | string, format: ApiSpecExportFormat) =>
    request.get<ApiSpecExportPayload>(`/projects/${projectId}/api-specs/export`, {
      params: { format },
    }),

  genDoc: (projectId: number | string, specId: number | string, lang: ApiSpecLanguage) =>
    request.post<ApiSpec>(`/projects/${projectId}/api-specs/${specId}/gen-doc`, undefined, {
      params: { lang },
    }),

  genTest: (projectId: number | string, specId: number | string, lang: ApiSpecLanguage) =>
    request.post<GenApiTestResponse>(`/projects/${projectId}/api-specs/${specId}/gen-test`, undefined, {
      params: { lang },
    }),

  batchGenDoc: (projectId: number | string, data: BatchGenDocRequest) =>
    request.post<BatchGenDocResponse>(
      `/projects/${projectId}/api-specs/batch-gen-doc`,
      normalizePayload(data)
    ),

  listExamples: (projectId: number | string, specId: number | string) =>
    request.get<ApiSpecExamplesResponse>(`/projects/${projectId}/api-specs/${specId}/examples`),

  createExample: (projectId: number | string, specId: number | string, data: CreateApiExampleRequest) =>
    request.post<ApiSpecExample>(`/projects/${projectId}/api-specs/${specId}/examples`, normalizePayload(data)),

  listCategories: (projectId: number | string) =>
    request.get<ProjectCategoryListResponse>(`/projects/${projectId}/categories`, {
      params: { tree: true },
    }),

  getMyRole: (projectId: number | string) =>
    request.get<ProjectMemberRoleResponse>(`/projects/${projectId}/members/me`),
};

export type ApiSpecService = typeof apiSpecService;
