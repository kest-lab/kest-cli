'use client';

import Link from 'next/link';
import { useState } from 'react';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { buildApiPath } from '@/config/api';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { ROUTES } from '@/constants/routes';
import { useUserSearch, useUsers } from '@/hooks/use-users';
import { useT } from '@/i18n/client';
import { formatDate } from '@/utils';

const PAGE_SIZE = 10;

const resolveStatusLabel = (status: number | string, activeLabel: string, inactiveLabel: string) => {
  if (status === 1 || status === 'active') {
    return activeLabel;
  }

  if (status === 0 || status === 'inactive') {
    return inactiveLabel;
  }

  return String(status);
};

export function UsersList() {
  const t = useT();
  const [page, setPage] = useState(1);
  const [searchInput, setSearchInput] = useState('');
  const [searchQuery, setSearchQuery] = useState('');
  const { data, isLoading, isFetching } = useUsers({ page, perPage: PAGE_SIZE });
  const { data: searchResults, isFetching: isSearching } = useUserSearch(searchQuery);
  const usersPath = buildApiPath('/users');

  const canGoPrev = page > 1;
  const canGoNext = Boolean(data?.meta && data.meta.current_page < data.meta.last_page);
  const isSearchMode = searchQuery.trim().length > 0;
  const users = isSearchMode ? (searchResults || []) : (data?.items || []);

  const handleSearch = (event: React.FormEvent) => {
    event.preventDefault();
    // 搜索模式和分页列表共存：有关键词时走搜索接口，没有关键词时回退到分页列表。
    setSearchQuery(searchInput.trim());
    setPage(1);
  };

  return (
    <div className="flex-1 space-y-6 p-6 pt-6">
      <div className="space-y-1">
        <h1 className="text-3xl font-bold tracking-tight">{t.console('users.title')}</h1>
        <p className="text-sm text-muted-foreground">
          {t.console('users.listConnectedDescription', { path: usersPath })}
        </p>
      </div>

      <Card>
        <CardHeader className="flex flex-col gap-1 md:flex-row md:items-center md:justify-between">
          <div>
            <CardTitle>{t.console('users.directory')}</CardTitle>
            <CardDescription>
              {isSearchMode
                ? t.console('users.searchResults', { query: searchQuery })
                : data?.meta
                ? t.console('users.pageSummary', {
                    page: data.meta.current_page,
                    pages: data.meta.last_page,
                    total: data.meta.total,
                  })
                : t.console('users.loadingList')}
            </CardDescription>
          </div>
          <div className="flex items-center gap-2 text-sm text-muted-foreground">
            {isSearching ? <span>{t.console('users.searching')}</span> : null}
            {isFetching && !isLoading && !isSearchMode ? <span>{t.console('users.refreshing')}</span> : null}
          </div>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSearch} className="mb-4 flex flex-col gap-3 md:flex-row">
            <Input
              value={searchInput}
              onChange={(event) => setSearchInput(event.target.value)}
              placeholder={t.console('users.searchPlaceholder')}
            />
            <Button type="submit" variant="outline">{t.console('users.search')}</Button>
            {isSearchMode ? (
              <Button type="button" variant="ghost" onClick={() => {
                setSearchInput('');
                setSearchQuery('');
              }}>
                {t.console('users.clear')}
              </Button>
            ) : null}
          </form>

          {isLoading && !isSearchMode ? (
            <div className="space-y-3">
              <div className="h-12 animate-pulse rounded bg-muted" />
              <div className="h-12 animate-pulse rounded bg-muted" />
              <div className="h-12 animate-pulse rounded bg-muted" />
            </div>
          ) : (
            <>
              <div className="overflow-hidden rounded-xl border">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>{t.console('users.id')}</TableHead>
                      <TableHead>{t.console('users.username')}</TableHead>
                      <TableHead>{t.console('users.email')}</TableHead>
                      <TableHead>{t.console('users.nickname')}</TableHead>
                      <TableHead>{t.console('users.status')}</TableHead>
                      <TableHead>{t.console('users.created')}</TableHead>
                      <TableHead className="text-right">{t.console('users.action')}</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {users.map((user) => (
                      <TableRow key={user.id}>
                        <TableCell className="font-mono text-xs">{user.id}</TableCell>
                        <TableCell className="font-medium">{user.username}</TableCell>
                        <TableCell>{user.email}</TableCell>
                        <TableCell>{user.nickname || '—'}</TableCell>
                        <TableCell>
                          <Badge variant="outline">
                            {resolveStatusLabel(user.status, t('common.active'), t('common.inactive'))}
                          </Badge>
                        </TableCell>
                        <TableCell>{formatDate(user.created_at, 'YYYY-MM-DD')}</TableCell>
                        <TableCell className="text-right">
                          <Button asChild size="sm" variant="outline">
                            <Link href={`${ROUTES.CONSOLE.USERS}/${user.id}`}>{t.console('users.open')}</Link>
                          </Button>
                        </TableCell>
                      </TableRow>
                    ))}
                    {users.length === 0 ? (
                      <TableRow>
                        <TableCell colSpan={7} className="py-10 text-center text-muted-foreground">
                          {t.console('users.noUsersFound')}
                        </TableCell>
                      </TableRow>
                    ) : null}
                  </TableBody>
                </Table>
              </div>

              {!isSearchMode ? (
                <div className="mt-4 flex items-center justify-between">
                  <Button variant="outline" onClick={() => setPage((current) => current - 1)} disabled={!canGoPrev}>
                    {t('common.previous')}
                  </Button>
                  <span className="text-sm text-muted-foreground">{t.console('users.page', { page })}</span>
                  <Button variant="outline" onClick={() => setPage((current) => current + 1)} disabled={!canGoNext}>
                    {t('common.next')}
                  </Button>
                </div>
              ) : null}
            </>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
