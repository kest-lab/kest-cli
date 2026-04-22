'use client';

import { useMutation, useQueryClient } from '@tanstack/react-query';
import { toast } from 'sonner';
import { collectionKeys } from '@/hooks/use-collections';
import { requestKeys } from '@/hooks/use-requests';
import { importerService } from '@/services/importer';
import type {
  ImportMarkdownCollectionRequest,
  ImportPostmanCollectionRequest,
} from '@/types/importer';

export function useImportPostmanCollection(projectId: number | string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: ImportPostmanCollectionRequest) =>
      importerService.importPostman(projectId, data),
    onSuccess: (result) => {
      queryClient.invalidateQueries({ queryKey: collectionKeys.project(projectId) });
      queryClient.invalidateQueries({ queryKey: requestKeys.project(projectId) });
      toast.success(result.message || 'Postman collection imported');
    },
  });
}

export function useImportMarkdownCollection(projectId: number | string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: ImportMarkdownCollectionRequest) =>
      importerService.importMarkdown(projectId, data),
    onSuccess: (result) => {
      queryClient.invalidateQueries({ queryKey: collectionKeys.project(projectId) });
      queryClient.invalidateQueries({ queryKey: requestKeys.project(projectId) });
      toast.success(
        `Imported ${result.requests_created} request${result.requests_created === 1 ? '' : 's'} from API Markdown`
      );
    },
  });
}
