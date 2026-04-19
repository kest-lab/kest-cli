import { env } from '@/config/env';

interface ErrorMetadata {
  userId?: string;
  url?: string;
  [key: string]: unknown;
}

type ErrorHandler = (error: Error, metadata?: ErrorMetadata) => void;

class ErrorTrackingService {
  private static instance: ErrorTrackingService;
  private handlers: ErrorHandler[] = [];
  private defaultMetadata: ErrorMetadata = {};

  public static getInstance(): ErrorTrackingService {
    if (!ErrorTrackingService.instance) {
      ErrorTrackingService.instance = new ErrorTrackingService();
    }
    return ErrorTrackingService.instance;
  }

  public init(defaultMetadata: Partial<ErrorMetadata> = {}): void {
    this.defaultMetadata = { ...defaultMetadata };

    if (typeof window !== 'undefined') {
      window.addEventListener('error', (e) => this.captureError(e.error));
      window.addEventListener('unhandledrejection', (e) => this.captureError(e.reason));
    }
  }

  public addHandler(handler: ErrorHandler): void {
    this.handlers.push(handler);
  }

  public captureError(error: unknown, metadata?: ErrorMetadata): void {
    const normalizedError =
      error instanceof Error
        ? error
        : new Error(typeof error === 'string' ? error : 'Unknown error');
    const mergedMetadata = { ...this.defaultMetadata, ...metadata };

    if (env.NODE_ENV === 'development') {
      console.error('[ErrorTracking]', normalizedError, mergedMetadata);
    }

    this.handlers.forEach((handler) => handler(normalizedError, mergedMetadata));
  }
}

export const errorTracking = ErrorTrackingService.getInstance();

export function setupErrorTracking(userId?: string): void {
  errorTracking.init({ userId });
}
