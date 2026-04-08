import request from '@/http';
import type { ExampleQuerySchema, CreateExampleRequest, UpdateExampleRequest } from '@/types/example';

/**
 * Example Service (Template)
 * 
 * Demonstrates the standard pattern for API service definitions.
 * All methods are stateless pure functions using the centralized request utility.
 */
export const exampleService = {
  /**
   * Fetch a paginated list of items
   */
  getList: (params?: ExampleQuerySchema) => 
    request.get('/example', { params }),
    
  /**
   * Fetch a single item by ID
   */
  getDetail: (id: string) => 
    request.get(`/example/${id}`),
    
  /**
   * Create a new item
   */
  create: (data: CreateExampleRequest) => 
    request.post('/example', data),
    
  /**
   * Update an existing item (Partial update)
   */
  update: (id: string, data: UpdateExampleRequest) => 
    request.patch(`/example/${id}`, data),
    
  /**
   * Delete an item
   */
  delete: (id: string) => 
    request.delete(`/example/${id}`),
};

export type ExampleService = typeof exampleService;
