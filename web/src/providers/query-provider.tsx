'use client';

import { QueryClientProvider } from '@tanstack/react-query';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import { queryClient } from '@/config';
import { PropsWithChildren, useState } from 'react';

/**
 * React Query provider.
 * Provides a QueryClient instance for the entire client application.
 */
export function QueryProvider({ children }: PropsWithChildren) {
  // Ensure a stable QueryClient instance on the client.
  const [client] = useState(() => queryClient);

  return (
    <QueryClientProvider client={client}>
      {children}
      <ReactQueryDevtools initialIsOpen={false} position="bottom" />
    </QueryClientProvider>
  );
}
