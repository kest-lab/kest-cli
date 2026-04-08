import { env } from '@/config/env';
import { useAuthStore } from '@/store/auth-store';

interface ErrorMetadata {
  userId?: string;
  url?: string;
  [key: string]: unknown;
}

type ErrorHandler = (error: Error, metadata?: ErrorMetadata) => void;

class ErrorTrackingService {
  private static instance: ErrorTrackingService;
  private handlers: ErrorHandler[] = [];

  public static getInstance(): ErrorTrackingService {
    if (!ErrorTrackingService.instance) {
      ErrorTrackingService.instance = new ErrorTrackingService();
    }
    return ErrorTrackingService.instance;
  }

  public init(defaultMetadata: Partial<ErrorMetadata> = {}): void {
    if (typeof window !== 'undefined') {
      window.addEventListener('error', (e) => this.captureError(e.error));
      window.addEventListener('unhandledrejection', (e) => this.captureError(e.reason));
    }
  }

  public addHandler(handler: ErrorHandler): void {
    this.handlers.push(handler);
  }

  public captureError(error: any, metadata?: ErrorMetadata): void {
    if (env.NODE_ENV === 'development') {
      console.error('[ErrorTracking]', error, metadata);
    }
    this.handlers.forEach(h => h(error, metadata));
  }
}

export const errorTracking = ErrorTrackingService.getInstance();

export function setupErrorTracking(userId?: string): void {
  errorTracking.init({ userId });
}
