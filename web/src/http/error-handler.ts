import { toast } from 'sonner';
import { ApiError } from './request';
import { env } from '@/config/env';
import { ErrorCode } from './codes';

/**
 * Global Error Handler Configuration
 */
export interface ErrorHandlerConfig {
  silent?: boolean;
  notify?: boolean;
  fallbackMessage?: string;
}

const DEFAULT_CONFIG: ErrorHandlerConfig = {
  silent: false,
  notify: true,
  fallbackMessage: 'An unexpected error occurred',
};

function getDebugErrorPayload(error: unknown) {
  if (error instanceof ApiError) {
    return {
      name: error.name,
      message: error.message,
      code: error.code,
      status: error.status,
      stack: error.stack,
    };
  }

  if (error instanceof Error) {
    return {
      name: error.name,
      message: error.message,
      stack: error.stack,
    };
  }

  return error;
}

/**
 * Centralized error handler for API and application errors.
 */
export function handleError(error: unknown, config: ErrorHandlerConfig = {}): void {
  const mergedConfig = { ...DEFAULT_CONFIG, ...config };
  
  if (mergedConfig.silent) return;

  // 避免同一个错误被 axios 拦截器和 React Query 全局 onError 重复处理。
  if (error instanceof ApiError && error.isHandled) {
    return;
  }

  let message = mergedConfig.fallbackMessage || 'Error';
  let errorCode: string | number | undefined;

  if (error instanceof ApiError) {
    errorCode = error.code;
    message = error.message;
    
    // Automatically map status-based messages if no specific message is provided
    // or if the message is too generic
    if (error.status && !error.message) {
      switch (error.status) {
        case 401: message = 'Session expired, please login again'; break;
        case 403: message = 'You do not have permission to perform this action'; break;
        case 404: message = 'The requested resource was not found'; break;
        case 500: message = 'Something went wrong on our server. Please try again later'; break;
      }
    }

    // You can also add mapping based on ErrorCode here if needed
    if (errorCode === ErrorCode.TOKEN_EXPIRED) {
      message = 'Your session has expired. Please log in again.';
    }
  } else if (error instanceof Error) {
    message = error.message;
  }

  if (mergedConfig.notify) {
    toast.error(message, {
      description: errorCode ? `Error Code: ${errorCode}` : undefined,
    });
  }

  if (error instanceof ApiError) {
    error.isHandled = true;
  }

  if (env.NODE_ENV === 'development' && env.NEXT_PUBLIC_DEBUG_API_ERRORS) {
    // 这里输出的是“已处理错误”的调试信息，不应该再触发开发环境红框。
    console.warn('[GlobalErrorHandler]', {
      message,
      errorCode,
      originalError: getDebugErrorPayload(error),
    });
  }
}
