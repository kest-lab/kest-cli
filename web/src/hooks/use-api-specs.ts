'use client';

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { toast } from 'sonner';
import { apiSpecService } from '@/services/api-spec';
import type {
  ApiSpecLanguage,
  ApiSpecListParams,
  BatchGenDocRequest,
  CreateApiExampleRequest,
  CreateApiSpecRequest,
  ImportApiSpecsRequest,
  UpdateApiSpecRequest,
} from '@/types/api-spec';

// API Specifications 域的 React Query key。
// 作用：统一管理规格列表、详情、示例、分类、成员角色和 AI 结果缓存。
export const apiSpecKeys = {
  all: ['api-specs'] as const,
  project: (projectId: number | string) => [...apiSpecKeys.all, 'project', projectId] as const,
  lists: (projectId: number | string) => [...apiSpecKeys.project(projectId), 'lists'] as const,
  list: (params: ApiSpecListParams) => [...apiSpecKeys.lists(params.projectId), params] as const,
  spec: (projectId: number | string, specId: number | string) =>
    [...apiSpecKeys.project(projectId), 'spec', specId] as const,
  detail: (projectId: number | string, specId: number | string) =>
    [...apiSpecKeys.spec(projectId, specId), 'detail'] as const,
  full: (projectId: number | string, specId: number | string) =>
    [...apiSpecKeys.spec(projectId, specId), 'full'] as const,
  examples: (projectId: number | string, specId: number | string) =>
    [...apiSpecKeys.spec(projectId, specId), 'examples'] as const,
  generatedTests: (projectId: number | string, specId: number | string) =>
    [...apiSpecKeys.spec(projectId, specId), 'generated-tests'] as const,
  generatedTest: (projectId: number | string, specId: number | string, lang: ApiSpecLanguage) =>
    [...apiSpecKeys.generatedTests(projectId, specId), lang] as const,
  categories: (projectId: number | string) =>
    [...apiSpecKeys.project(projectId), 'categories'] as const,
  memberRole: (projectId: number | string) =>
    [...apiSpecKeys.project(projectId), 'member-role'] as const,
};

// 规格列表查询。
// 作用：按项目维度拉取带筛选和分页参数的 API 规格列表。
export function useApiSpecs(params: ApiSpecListParams) {
  return useQuery({
    queryKey: apiSpecKeys.list(params),
    queryFn: () => apiSpecService.list(params),
    placeholderData: (previousData) => previousData,
  });
}

// 规格轻量详情查询。
// 作用：获取单条规格基础详情，适合编辑前预取或局部刷新。
export function useApiSpec(projectId?: number | string, specId?: number | string) {
  return useQuery({
    queryKey: apiSpecKeys.detail(projectId ?? 'unknown', specId ?? 'unknown'),
    queryFn: () => apiSpecService.getById(projectId as number | string, specId as number | string),
    enabled: projectId !== undefined && projectId !== null && specId !== undefined && specId !== null,
  });
}

// 含 examples 的完整规格详情查询。
// 作用：驱动右侧详情面板，统一承载 overview、docs 和 examples 主内容。
export function useApiSpecFull(projectId?: number | string, specId?: number | string) {
  return useQuery({
    queryKey: apiSpecKeys.full(projectId ?? 'unknown', specId ?? 'unknown'),
    queryFn: () => apiSpecService.getFullById(projectId as number | string, specId as number | string),
    enabled: projectId !== undefined && projectId !== null && specId !== undefined && specId !== null,
  });
}

// 规格 examples 列表查询。
// 作用：独立维护 examples 缓存，便于创建示例后精确刷新。
export function useApiSpecExamples(projectId?: number | string, specId?: number | string) {
  return useQuery({
    queryKey: apiSpecKeys.examples(projectId ?? 'unknown', specId ?? 'unknown'),
    queryFn: () => apiSpecService.listExamples(projectId as number | string, specId as number | string),
    enabled: projectId !== undefined && projectId !== null && specId !== undefined && specId !== null,
  });
}

// 项目分类查询。
// 作用：为 category_id 选择器提供树形分类列表。
export function useProjectApiCategories(projectId?: number | string) {
  return useQuery({
    queryKey: apiSpecKeys.categories(projectId ?? 'unknown'),
    queryFn: () => apiSpecService.listCategories(projectId as number | string),
    enabled: projectId !== undefined && projectId !== null,
    staleTime: 60_000,
  });
}

// 当前用户在项目中的角色查询。
// 作用：根据 read/write/admin/owner 控制页面上的可写操作按钮。
export function useProjectMemberRole(projectId?: number | string) {
  return useQuery({
    queryKey: apiSpecKeys.memberRole(projectId ?? 'unknown'),
    queryFn: () => apiSpecService.getMyRole(projectId as number | string),
    enabled: projectId !== undefined && projectId !== null,
    staleTime: 60_000,
  });
}

// AI 生成测试内容缓存查询。
// 作用：不触发网络请求，只订阅 mutation 写入的 flow_content 缓存。
export function useGeneratedApiTest(projectId?: number | string, specId?: number | string, lang: ApiSpecLanguage = 'en') {
  return useQuery<string | null>({
    queryKey: apiSpecKeys.generatedTest(projectId ?? 'unknown', specId ?? 'unknown', lang),
    queryFn: async () => null,
    enabled: false,
    initialData: null,
    staleTime: Number.POSITIVE_INFINITY,
  });
}

// 创建规格 mutation。
// 作用：创建成功后刷新列表，并预写入新规格详情缓存。
export function useCreateApiSpec(projectId: number | string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateApiSpecRequest) => apiSpecService.create(projectId, data),
    onSuccess: (spec) => {
      queryClient.invalidateQueries({ queryKey: apiSpecKeys.lists(projectId) });
      queryClient.setQueryData(apiSpecKeys.detail(projectId, spec.id), spec);
      queryClient.invalidateQueries({ queryKey: apiSpecKeys.full(projectId, spec.id) });
      toast.success(`Created API spec ${spec.method} ${spec.path}`);
    },
  });
}

// 更新规格 mutation。
// 作用：更新成功后同步刷新列表、详情和完整详情缓存。
export function useUpdateApiSpec(projectId: number | string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ specId, data }: { specId: number | string; data: UpdateApiSpecRequest }) =>
      apiSpecService.update(projectId, specId, data),
    onSuccess: (spec) => {
      queryClient.invalidateQueries({ queryKey: apiSpecKeys.lists(projectId) });
      queryClient.setQueryData(apiSpecKeys.detail(projectId, spec.id), spec);
      queryClient.invalidateQueries({ queryKey: apiSpecKeys.full(projectId, spec.id) });
      queryClient.invalidateQueries({ queryKey: apiSpecKeys.examples(projectId, spec.id) });
      toast.success(`Updated API spec ${spec.method} ${spec.path}`);
    },
  });
}

// 删除规格 mutation。
// 作用：删除成功后清理该规格相关缓存并刷新列表。
export function useDeleteApiSpec(projectId: number | string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (specId: number | string) => apiSpecService.delete(projectId, specId),
    onSuccess: (_, specId) => {
      queryClient.invalidateQueries({ queryKey: apiSpecKeys.lists(projectId) });
      queryClient.removeQueries({ queryKey: apiSpecKeys.spec(projectId, specId) });
      toast.success('API spec deleted');
    },
  });
}

// 批量导入 mutation。
// 作用：导入成功后刷新列表并提示 upsert 已完成。
export function useImportApiSpecs(projectId: number | string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: ImportApiSpecsRequest) => apiSpecService.import(projectId, data),
    onSuccess: (result) => {
      queryClient.invalidateQueries({ queryKey: apiSpecKeys.lists(projectId) });
      toast.success(result.message || 'Specs imported successfully');
    },
  });
}

// 导出 mutation。
// 作用：触发导出接口，页面层拿到结果后负责下载文件。
export function useExportApiSpecs(projectId: number | string) {
  return useMutation({
    mutationFn: ({ format }: { format: 'json' | 'openapi' | 'swagger' | 'markdown' }) =>
      apiSpecService.export(projectId, format),
  });
}

// AI 生成文档 mutation。
// 作用：单条生成成功后刷新对应规格详情和列表缓存。
export function useGenApiDoc(projectId: number | string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ specId, lang }: { specId: number | string; lang: ApiSpecLanguage }) =>
      apiSpecService.genDoc(projectId, specId, lang),
    onSuccess: (spec) => {
      queryClient.invalidateQueries({ queryKey: apiSpecKeys.lists(projectId) });
      queryClient.setQueryData(apiSpecKeys.detail(projectId, spec.id), spec);
      queryClient.invalidateQueries({ queryKey: apiSpecKeys.full(projectId, spec.id) });
      toast.success(`Generated ${spec.path} documentation`);
    },
  });
}

// AI 生成测试 mutation。
// 作用：把 flow_content 缓存在 query 中，方便详情页直接展示。
export function useGenApiTest(projectId: number | string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ specId, lang }: { specId: number | string; lang: ApiSpecLanguage }) =>
      apiSpecService.genTest(projectId, specId, lang),
    onSuccess: (result, variables) => {
      queryClient.setQueryData(
        apiSpecKeys.generatedTest(projectId, variables.specId, variables.lang),
        result.flow_content
      );
      toast.success('Generated Kest flow test');
    },
  });
}

// 批量生成文档 mutation。
// 作用：批量任务提交成功后统一失效当前项目下的规格缓存。
export function useBatchGenApiDocs(projectId: number | string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: BatchGenDocRequest) => apiSpecService.batchGenDoc(projectId, data),
    onSuccess: (result) => {
      queryClient.invalidateQueries({ queryKey: apiSpecKeys.project(projectId) });
      toast.success(`Queued ${result.queued} specs, skipped ${result.skipped}`);
    },
  });
}

// 创建示例 mutation。
// 作用：创建成功后刷新 examples 列表和完整规格详情。
export function useCreateApiExample(projectId: number | string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ specId, data }: { specId: number | string; data: CreateApiExampleRequest }) =>
      apiSpecService.createExample(projectId, specId, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: apiSpecKeys.examples(projectId, variables.specId) });
      queryClient.invalidateQueries({ queryKey: apiSpecKeys.full(projectId, variables.specId) });
      toast.success('API example created');
    },
  });
}
