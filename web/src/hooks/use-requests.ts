'use client';

import { useMutation, useQueryClient } from '@tanstack/react-query';
import { toast } from 'sonner';
import { requestService } from '@/services/request';
import type { CreateRequestRequest } from '@/types/request';

export const requestKeys = {
  all: ['requests'] as const,
  project: (projectId: number | string) => [...requestKeys.all, 'project', projectId] as const,
  collection: (projectId: number | string, collectionId: number | string) =>
    [...requestKeys.project(projectId), 'collection', collectionId] as const,
  list: (projectId: number | string, collectionId: number | string) =>
    [...requestKeys.collection(projectId, collectionId), 'list'] as const,
  detail: (
    projectId: number | string,
    collectionId: number | string,
    requestId: number | string
  ) => [...requestKeys.collection(projectId, collectionId), 'detail', requestId] as const,
};

export function useCreateRequest(projectId: number | string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      collectionId,
      data,
    }: {
      collectionId: number | string;
      data: CreateRequestRequest;
    }) => requestService.create(projectId, collectionId, data),
    onSuccess: (createdRequest, variables) => {
      queryClient.invalidateQueries({
        queryKey: requestKeys.collection(projectId, variables.collectionId),
      });
      toast.success(`Created request "${createdRequest.name}"`);
    },
  });
}
