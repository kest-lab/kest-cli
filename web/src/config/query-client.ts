import { QueryClient } from '@tanstack/react-query';

/**
 * TanStack Query Configuration
 * Optimized for performance
 */

export const queryClient = new QueryClient({
    defaultOptions: {
        queries: {
            // Stale time: data considered fresh for 1 minute
            staleTime: 1000 * 60,

            // Cache time: unused data kept for 5 minutes
            gcTime: 1000 * 60 * 5,

            // Retry failed requests
            retry: 1,
            retryDelay: (attemptIndex) => Math.min(1000 * 2 ** attemptIndex, 30000),

            // Refetch configuration
            refetchOnWindowFocus: false,
            refetchOnReconnect: true,
            refetchOnMount: true,

            // Network mode
            networkMode: 'online',
        },
        mutations: {
            retry: 0,
            networkMode: 'online',
        },
    },
});

export default queryClient;
