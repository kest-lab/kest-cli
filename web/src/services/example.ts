import request from '@/http';
import type { ExampleQuerySchema, CreateExampleRequest, UpdateExampleRequest } from '@/types/example';

/**
 * Example Service
 */
export const exampleService = {
    getList: (params?: ExampleQuerySchema) =>
        request.get('/example', { params }),

    getDetail: (id: string) =>
        request.get(`/example/${id}`),

    create: (data: CreateExampleRequest) =>
        request.post('/example', data),

    update: (id: string, data: UpdateExampleRequest) =>
        request.patch(`/example/${id}`, data),

    delete: (id: string) =>
        request.delete(`/example/${id}`),
};

export type ExampleService = typeof exampleService;
