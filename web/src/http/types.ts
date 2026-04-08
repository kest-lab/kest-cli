// HTTP related types
export interface ApiResponse<T = unknown> {
  data: T;
  message?: string;
  success: boolean;
  code?: number;
}

export interface ApiError {
  message: string;
  code?: number;
  details?: unknown;
}

export interface RequestHeaders {
  [key: string]: string;
}

export interface RequestConfig {
  baseURL?: string;
  timeout?: number;
  headers?: RequestHeaders;
  skipAuth?: boolean;
  unwrapData?: boolean;
  retry?: number;
}
