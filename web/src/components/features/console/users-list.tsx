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
import { formatDate } from '@/utils';

const PAGE_SIZE = 10;

const resolveStatusLabel = (status: number | string) => {
  if (status === 1 || status === 'active') {
    return 'Active';
  }

  if (status === 0 || status === 'inactive') {
    return 'Inactive';
  }

  return String(status);
};

export function UsersList() {
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
        <h1 className="text-3xl font-bold tracking-tight">Users</h1>
        <p className="text-sm text-muted-foreground">
          Connected to <code>GET {usersPath}</code> for paginated user management.
        </p>
      </div>

      <Card>
        <CardHeader className="flex flex-col gap-1 md:flex-row md:items-center md:justify-between">
          <div>
            <CardTitle>User Directory</CardTitle>
            <CardDescription>
              {isSearchMode
                ? `Search results for "${searchQuery}"`
                : data?.meta
                ? `Page ${data.meta.current_page} of ${data.meta.last_page}, ${data.meta.total} total users`
                : 'Loading user list'}
            </CardDescription>
          </div>
          <div className="flex items-center gap-2 text-sm text-muted-foreground">
            {isSearching ? <span>Searching…</span> : null}
            {isFetching && !isLoading && !isSearchMode ? <span>Refreshing…</span> : null}
          </div>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSearch} className="mb-4 flex flex-col gap-3 md:flex-row">
            <Input
              value={searchInput}
              onChange={(event) => setSearchInput(event.target.value)}
              placeholder="Search by username or email"
            />
            <Button type="submit" variant="outline">Search</Button>
            {isSearchMode ? (
              <Button type="button" variant="ghost" onClick={() => {
                setSearchInput('');
                setSearchQuery('');
              }}>
                Clear
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
                      <TableHead>ID</TableHead>
                      <TableHead>Username</TableHead>
                      <TableHead>Email</TableHead>
                      <TableHead>Nickname</TableHead>
                      <TableHead>Status</TableHead>
                      <TableHead>Created</TableHead>
                      <TableHead className="text-right">Action</TableHead>
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
                          <Badge variant="outline">{resolveStatusLabel(user.status)}</Badge>
                        </TableCell>
                        <TableCell>{formatDate(user.created_at, 'YYYY-MM-DD')}</TableCell>
                        <TableCell className="text-right">
                          <Button asChild size="sm" variant="outline">
                            <Link href={`${ROUTES.CONSOLE.USERS}/${user.id}`}>Open</Link>
                          </Button>
                        </TableCell>
                      </TableRow>
                    ))}
                    {users.length === 0 ? (
                      <TableRow>
                        <TableCell colSpan={7} className="py-10 text-center text-muted-foreground">
                          No users found.
                        </TableCell>
                      </TableRow>
                    ) : null}
                  </TableBody>
                </Table>
              </div>

              {!isSearchMode ? (
                <div className="mt-4 flex items-center justify-between">
                  <Button variant="outline" onClick={() => setPage((current) => current - 1)} disabled={!canGoPrev}>
                    Previous
                  </Button>
                  <span className="text-sm text-muted-foreground">Page {page}</span>
                  <Button variant="outline" onClick={() => setPage((current) => current + 1)} disabled={!canGoNext}>
                    Next
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
