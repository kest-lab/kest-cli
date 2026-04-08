import { QueryClient, QueryCache, MutationCache } from '@tanstack/react-query';
import { handleError } from '@/http/error-handler';

// Create a React Query client.
export const queryClient = new QueryClient({
  queryCache: new QueryCache({
    onError: (error) => {
      // Background query errors
      handleError(error);
    },
  }),
  mutationCache: new MutationCache({
    onError: (error) => {
      // Global mutation errors (as a safety net)
      handleError(error);
    },
  }),
  defaultOptions: {
    queries: {
      // Baseline defaults.
      refetchOnWindowFocus: false,
      retry: 1,
      staleTime: 5 * 60 * 1000,
    },
    mutations: {
      // Baseline defaults.
      retry: 1,
    },
  },
});

export default queryClient;
