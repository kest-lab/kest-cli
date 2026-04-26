'use client';

import { QueryClientProvider } from '@tanstack/react-query';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import { queryClient } from '@/config';
import { setErrorMessageResolver } from '@/http/error-handler';
import { useT } from '@/i18n/client';
import { PropsWithChildren, useEffect, useState } from 'react';

function LocalizedErrorMessages() {
  const t = useT();

  useEffect(() => {
    setErrorMessageResolver({
      unexpected: t.errors('unexpected'),
      genericError: t.errors('error'),
      tokenExpired: t.errors('tokenExpired'),
      statusMessages: {
        401: t.errors('sessionExpiredLoginAgain'),
        403: t.errors('permissionDenied'),
        404: t.errors('resourceNotFound'),
        500: t.errors('serverTryLater'),
      },
      formatErrorCode: (code) => t.errors('errorCode', { code }),
    });

    return () => setErrorMessageResolver(undefined);
  }, [t]);

  return null;
}

/**
 * React Query provider.
 * Provides a QueryClient instance for the entire client application.
 */
export function QueryProvider({ children }: PropsWithChildren) {
  // Ensure a stable QueryClient instance on the client.
  const [client] = useState(() => queryClient);

  return (
    <QueryClientProvider client={client}>
      <LocalizedErrorMessages />
      {children}
      <ReactQueryDevtools initialIsOpen={false} position="bottom" />
    </QueryClientProvider>
  );
}
