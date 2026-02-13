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

/**
 * Centralized error handler for API and application errors.
 */
export function handleError(error: unknown, config: ErrorHandlerConfig = {}): void {
    const mergedConfig = { ...DEFAULT_CONFIG, ...config };

    if (mergedConfig.silent) return;

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

    if (mergedConfig.notify && typeof window !== 'undefined') {
        toast.error(message, {
            description: errorCode ? `Error Code: ${errorCode}` : undefined,
        });
    }

    if (env.MODE === 'development') {
        console.error('[GlobalErrorHandler]', {
            message,
            errorCode,
            originalError: error
        });
    }
}
