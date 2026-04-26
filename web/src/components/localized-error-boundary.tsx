import { ErrorBoundary } from '@/components/error-boundary';
export function LocalizedErrorBoundary({
  children,
  fallbackTitle,
  fallbackDescription,
  retryLabel,
}: {
  children: React.ReactNode;
  fallbackTitle?: string;
  fallbackDescription?: string;
  retryLabel?: string;
}) {

  return (
    <ErrorBoundary
      fallbackTitle={fallbackTitle}
      fallbackDescription={fallbackDescription}
      retryLabel={retryLabel}
    >
      {children}
    </ErrorBoundary>
  );
}
