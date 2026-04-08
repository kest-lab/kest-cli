'use client';

import { useQuery } from '@tanstack/react-query';
import { userService } from '@/services/auth';
import type { UserListParams } from '@/types/auth';

export const userKeys = {
  all: ['users'] as const,
  lists: () => [...userKeys.all, 'list'] as const,
  list: (params: UserListParams) => [...userKeys.lists(), params] as const,
  searches: () => [...userKeys.all, 'search'] as const,
  search: (query: string, limit: number) => [...userKeys.searches(), query, limit] as const,
  details: () => [...userKeys.all, 'detail'] as const,
  detail: (id: number | string) => [...userKeys.details(), id] as const,
  infos: () => [...userKeys.all, 'info'] as const,
  info: (id: number | string) => [...userKeys.infos(), id] as const,
};

export function useUsers(params: UserListParams = {}) {
  return useQuery({
    queryKey: userKeys.list(params),
    queryFn: () => userService.list(params),
    // 翻页时保留上一页数据，减少表格闪烁。
    placeholderData: (previousData) => previousData,
  });
}

export function useUser(id?: number | string) {
  return useQuery({
    queryKey: userKeys.detail(id ?? 'unknown'),
    queryFn: () => userService.getById(id as number | string),
    enabled: id !== undefined && id !== null && id !== '',
  });
}

export function useUserSearch(query: string, limit = 20) {
  return useQuery({
    queryKey: userKeys.search(query, limit),
    queryFn: () => userService.search(query, limit),
    // 只有真正输入搜索词时才请求，避免空字符串触发无意义搜索。
    enabled: query.trim().length > 0,
  });
}

export function useUserInfo(id?: number | string) {
  return useQuery({
    queryKey: userKeys.info(id ?? 'unknown'),
    queryFn: () => userService.getInfo(id as number | string),
    enabled: id !== undefined && id !== null && id !== '',
  });
}
