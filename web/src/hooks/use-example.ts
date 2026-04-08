/**
 * Example React Query hooks
 * Encapsulating state, side effects, and cache management.
 */

'use client';

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { toast } from 'sonner';
import { exampleService } from '@/services/example';
import { useT } from '@/i18n/client';
import type { ExampleQuerySchema, CreateExampleRequest, UpdateExampleRequest } from '@/types/example';

// ============================================================================
// Query Keys
// ============================================================================

/**
 * Centralized query keys for the example domain
 */
export const exampleKeys = {
  all: ['examples'] as const,
  lists: () => [...exampleKeys.all, 'list'] as const,
  list: (params: ExampleQuerySchema) => [...exampleKeys.lists(), params] as const,
  details: () => [...exampleKeys.all, 'detail'] as const,
  detail: (id: string) => [...exampleKeys.details(), id] as const,
};

// ============================================================================
// Hooks
// ============================================================================

/**
 * Hook for fetching example items list
 */
export function useExamples(params?: ExampleQuerySchema) {
  return useQuery({
    queryKey: exampleKeys.list(params || {}),
    queryFn: () => exampleService.getList(params),
  });
}

/**
 * Hook for fetching example item detail
 */
export function useExample(id: string) {
  return useQuery({
    queryKey: exampleKeys.detail(id),
    queryFn: () => exampleService.getDetail(id),
    enabled: !!id,
  });
}

/**
 * Hook for creating a new example item
 * Implements Optimistic Updates (List appending)
 */
export function useCreateExample() {
  const queryClient = useQueryClient();
  const t = useT();

  return useMutation({
    mutationFn: (data: CreateExampleRequest) => exampleService.create(data),

    onMutate: async (newItem) => {
      // 1. Cancel outgoing list fetches
      await queryClient.cancelQueries({ queryKey: exampleKeys.lists() });

      // 2. Snapshot previous lists
      const previousLists = queryClient.getQueryData(exampleKeys.lists());

      // 3. Optimistically add to the cache
      // Note: In a real app, you'd match the specific list query key (e.g., page 1)
      queryClient.setQueriesData({ queryKey: exampleKeys.lists() }, (old: any) => {
        if (!old) return [newItem];
        if (Array.isArray(old)) return [newItem, ...old];
        if (old.data) return { ...old, data: [newItem, ...old.data] };
        return old;
      });

      return { previousLists };
    },

    onError: (err: any) => {
      // Rollback: Invalidate lists to get fresh data from server
      queryClient.invalidateQueries({ queryKey: exampleKeys.lists() });
      toast.error(err.message || 'Failed to create');
    },

    onSettled: (data, error) => {
      // Final synchronization
      queryClient.invalidateQueries({ queryKey: exampleKeys.lists() });
      if (data && !error) {
        toast.success(t.common?.('success') || 'Created successfully');
      }
    },
  });
}

/**
 * Hook for updating an existing example item
 * Implements Optimistic Updates with "Refetch-on-Failure" rollback strategy.
 */
export function useUpdateExample() {
  const queryClient = useQueryClient();
  const t = useT();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateExampleRequest }) =>
      exampleService.update(id, data),

    onMutate: async ({ id, data }) => {
      await queryClient.cancelQueries({ queryKey: exampleKeys.detail(id) });
      const previousItem = queryClient.getQueryData(exampleKeys.detail(id));

      if (previousItem) {
        queryClient.setQueryData(exampleKeys.detail(id), {
          ...(previousItem as any),
          ...data,
        });
      }

      return { previousItem };
    },

    onError: (err: any, { id }) => {
      queryClient.invalidateQueries({ queryKey: exampleKeys.detail(id) });
      toast.error(err.message || 'Failed to update');
    },

    onSettled: (data, error, { id }) => {
      queryClient.invalidateQueries({ queryKey: exampleKeys.detail(id) });
      queryClient.invalidateQueries({ queryKey: exampleKeys.lists() });
      
      if (data && !error) {
        toast.success(t.common?.('success') || 'Updated successfully');
      }
    },
  });
}

/**
 * Hook for deleting an example item
 * Implements Optimistic Updates (List filtering)
 */
export function useDeleteExample() {
  const queryClient = useQueryClient();
  const t = useT();

  return useMutation({
    mutationFn: (id: string) => exampleService.delete(id),

    onMutate: async (id) => {
      // 1. Cancel outgoing fetches
      await queryClient.cancelQueries({ queryKey: exampleKeys.lists() });
      await queryClient.cancelQueries({ queryKey: exampleKeys.detail(id) });

      // 2. Optimistically remove from all lists
      queryClient.setQueriesData({ queryKey: exampleKeys.lists() }, (old: any) => {
        if (!old) return [];
        if (Array.isArray(old)) return old.filter((item: any) => item.id !== id);
        if (old.data) return { ...old, data: old.data.filter((item: any) => item.id !== id) };
        return old;
      });

      return { id };
    },

    onError: (err: any, id) => {
      queryClient.invalidateQueries({ queryKey: exampleKeys.lists() });
      queryClient.invalidateQueries({ queryKey: exampleKeys.detail(id) });
      toast.error(err.message || 'Failed to delete');
    },

    onSettled: (data, error, id) => {
      queryClient.invalidateQueries({ queryKey: exampleKeys.lists() });
      queryClient.invalidateQueries({ queryKey: exampleKeys.detail(id) });
      
      if (!error) {
        toast.success(t.common?.('success') || 'Deleted successfully');
      }
    },
  });
}
