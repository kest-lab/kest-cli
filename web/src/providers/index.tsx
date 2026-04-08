'use client';

import { PropsWithChildren } from 'react';
import { QueryProvider } from './query-provider';

/**
 * Root providers component
 * Combines all providers for the application
 */
export function Providers({ children }: PropsWithChildren) {
  return (
    <QueryProvider>
      {children}
    </QueryProvider>
  );
}

// Export individual providers for selective usage
export { QueryProvider } from './query-provider';
