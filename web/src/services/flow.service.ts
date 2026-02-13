import { request } from '@/http';
import type {
    Flow,
    FlowDetail,
    FlowRun,
    CreateFlowRequest,
    SaveFlowRequest,
    CreateStepRequest,
    CreateEdgeRequest,
    FlowStep,
    FlowEdge,
} from '@/types/kest-api';

export const flowService = {
    // Flow CRUD
    list: (projectId: number) =>
        request.get<{ items: Flow[]; total: number }>(`/v1/projects/${projectId}/flows`),

    get: (projectId: number, flowId: number) =>
        request.get<FlowDetail>(`/v1/projects/${projectId}/flows/${flowId}`),

    create: (projectId: number, data: CreateFlowRequest) =>
        request.post<Flow>(`/v1/projects/${projectId}/flows`, data),

    update: (projectId: number, flowId: number, data: Partial<CreateFlowRequest>) =>
        request.patch<Flow>(`/v1/projects/${projectId}/flows/${flowId}`, data),

    save: (projectId: number, flowId: number, data: SaveFlowRequest) =>
        request.put<FlowDetail>(`/v1/projects/${projectId}/flows/${flowId}`, data),

    delete: (projectId: number, flowId: number) =>
        request.delete(`/v1/projects/${projectId}/flows/${flowId}`),

    // Steps
    createStep: (projectId: number, flowId: number, data: CreateStepRequest) =>
        request.post<FlowStep>(`/v1/projects/${projectId}/flows/${flowId}/steps`, data),

    deleteStep: (projectId: number, flowId: number, stepId: number) =>
        request.delete(`/v1/projects/${projectId}/flows/${flowId}/steps/${stepId}`),

    // Edges
    createEdge: (projectId: number, flowId: number, data: CreateEdgeRequest) =>
        request.post<FlowEdge>(`/v1/projects/${projectId}/flows/${flowId}/edges`, data),

    deleteEdge: (projectId: number, flowId: number, edgeId: number) =>
        request.delete(`/v1/projects/${projectId}/flows/${flowId}/edges/${edgeId}`),

    // Runs
    run: (projectId: number, flowId: number) =>
        request.post<FlowRun>(`/v1/projects/${projectId}/flows/${flowId}/run`),

    getRun: (projectId: number, flowId: number, runId: number) =>
        request.get<FlowRun>(`/v1/projects/${projectId}/flows/${flowId}/runs/${runId}`),

    listRuns: (projectId: number, flowId: number) =>
        request.get<{ items: FlowRun[]; total: number }>(`/v1/projects/${projectId}/flows/${flowId}/runs`),
};

export default flowService;
