import { buildApiUrl } from '@/config/api';
import request from '@/http';
import { getAuthTokens } from '@/store/auth-store';
import type {
  CreateFlowRequest,
  FlowDetail,
  FlowListResponse,
  FlowRun,
  FlowRunListResponse,
  FlowStreamStepEvent,
  ProjectFlow,
  SaveFlowRequest,
  StreamFlowRunOptions,
  UpdateFlowRequest,
} from '@/types/flow';

const normalizePayload = <T extends object>(payload: T) =>
  Object.fromEntries(
    Object.entries(payload as Record<string, unknown>).filter(([, value]) => value !== undefined)
  ) as T;

const readSSEEvent = (
  chunk: string,
  handlers: Required<Pick<StreamFlowRunOptions, 'onDone' | 'onStep'>>
) => {
  const lines = chunk
    .split('\n')
    .map((line) => line.trim())
    .filter(Boolean);

  let eventName = 'message';
  const dataParts: string[] = [];

  for (const line of lines) {
    if (line.startsWith('event:')) {
      eventName = line.slice(6).trim();
      continue;
    }

    if (line.startsWith('data:')) {
      dataParts.push(line.slice(5).trim());
    }
  }

  if (eventName === 'done') {
    handlers.onDone();
    return;
  }

  if (eventName !== 'step' || dataParts.length === 0) {
    return;
  }

  try {
    handlers.onStep(JSON.parse(dataParts.join('\n')) as FlowStreamStepEvent);
  } catch {
    // Ignore malformed SSE payloads and keep the stream alive.
  }
};

export const flowService = {
  list: (projectId: number | string) =>
    request.get<FlowListResponse>(`/projects/${projectId}/flows`),

  getById: (projectId: number | string, flowId: number | string) =>
    request.get<FlowDetail>(`/projects/${projectId}/flows/${flowId}`),

  create: (projectId: number | string, data: CreateFlowRequest) =>
    request.post<ProjectFlow>(`/projects/${projectId}/flows`, normalizePayload(data)),

  update: (projectId: number | string, flowId: number | string, data: UpdateFlowRequest) =>
    request.patch<ProjectFlow>(`/projects/${projectId}/flows/${flowId}`, normalizePayload(data)),

  delete: (projectId: number | string, flowId: number | string) =>
    request.delete<void>(`/projects/${projectId}/flows/${flowId}`),

  save: (projectId: number | string, flowId: number | string, data: SaveFlowRequest) =>
    request.put<FlowDetail>(`/projects/${projectId}/flows/${flowId}`, normalizePayload(data)),

  run: (projectId: number | string, flowId: number | string) =>
    request.post<FlowRun>(`/projects/${projectId}/flows/${flowId}/run`),

  listRuns: (projectId: number | string, flowId: number | string) =>
    request.get<FlowRunListResponse>(`/projects/${projectId}/flows/${flowId}/runs`),

  getRun: (projectId: number | string, flowId: number | string, runId: number | string) =>
    request.get<FlowRun>(`/projects/${projectId}/flows/${flowId}/runs/${runId}`),

  streamRun: async (
    projectId: number | string,
    flowId: number | string,
    runId: number | string,
    options: StreamFlowRunOptions = {}
  ) => {
    const { accessToken } = getAuthTokens();
    const response = await fetch(
      buildApiUrl(`/projects/${projectId}/flows/${flowId}/runs/${runId}/events`),
      {
        headers: accessToken
          ? {
              Authorization: `Bearer ${accessToken}`,
            }
          : undefined,
        signal: options.signal,
      }
    );

    if (!response.ok) {
      throw new Error(`Failed to stream flow run: ${response.status}`);
    }

    if (!response.body) {
      throw new Error('Flow run stream is not available');
    }

    const reader = response.body.getReader();
    const decoder = new TextDecoder();
    let buffer = '';
    const handlers = {
      onDone: options.onDone ?? (() => {}),
      onStep: options.onStep ?? (() => {}),
    };

    while (true) {
      const { done, value } = await reader.read();
      if (done) {
        break;
      }

      buffer += decoder.decode(value, { stream: true });
      const chunks = buffer.split('\n\n');
      buffer = chunks.pop() ?? '';

      for (const chunk of chunks) {
        readSSEEvent(chunk, handlers);
      }
    }

    if (buffer.trim()) {
      readSSEEvent(buffer, handlers);
    }
  },
};

export type FlowService = typeof flowService;
