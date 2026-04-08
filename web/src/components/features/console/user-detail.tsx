'use client';

import Link from 'next/link';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { buildApiPath } from '@/config/api';
import { ROUTES } from '@/constants/routes';
import { useUser, useUserInfo } from '@/hooks/use-users';
import { formatDate } from '@/utils';

const resolveStatusLabel = (status: number | string) => {
  if (status === 1 || status === 'active') {
    return 'Active';
  }

  if (status === 0 || status === 'inactive') {
    return 'Inactive';
  }

  return String(status);
};

export function UserDetail({ userId }: { userId: number }) {
  const { data: user, isLoading: isUserLoading } = useUser(userId);
  const { data: userInfo, isLoading: isUserInfoLoading } = useUserInfo(userId);
  const userPath = `${buildApiPath('/users')}/${userId}`;
  const userInfoPath = `${userPath}/info`;

  if (isUserLoading || isUserInfoLoading || !user || !userInfo) {
    return (
      <div className="flex-1 space-y-4 p-6 pt-6">
        <div className="h-9 w-48 animate-pulse rounded bg-muted" />
        <div className="grid gap-4 xl:grid-cols-2">
          <div className="h-80 animate-pulse rounded-xl bg-muted" />
          <div className="h-80 animate-pulse rounded-xl bg-muted" />
        </div>
      </div>
    );
  }

  return (
    <div className="flex-1 space-y-6 p-6 pt-6">
      <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div className="space-y-1">
          <h1 className="text-3xl font-bold tracking-tight">{user.nickname || user.username}</h1>
          <p className="text-sm text-muted-foreground">
            Connected to <code>GET {userPath}</code> and <code>GET {userInfoPath}</code>.
          </p>
        </div>
        <Button asChild variant="outline">
          <Link href={ROUTES.CONSOLE.USERS}>Back to users</Link>
        </Button>
      </div>

      <div className="grid gap-4 xl:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Full User Record</CardTitle>
            <CardDescription>
              Response from <code>GET {userPath}</code>.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4 text-sm">
            <div className="flex items-center justify-between rounded-lg border p-3">
              <span className="text-muted-foreground">Status</span>
              <Badge variant="outline">{resolveStatusLabel(user.status)}</Badge>
            </div>
            <div className="grid gap-3 md:grid-cols-2">
              <div className="rounded-lg border p-3">
                <div className="text-muted-foreground">Username</div>
                <div className="mt-1 font-medium">{user.username}</div>
              </div>
              <div className="rounded-lg border p-3">
                <div className="text-muted-foreground">Email</div>
                <div className="mt-1 font-medium">{user.email}</div>
              </div>
              <div className="rounded-lg border p-3">
                <div className="text-muted-foreground">Nickname</div>
                <div className="mt-1 font-medium">{user.nickname || '—'}</div>
              </div>
              <div className="rounded-lg border p-3">
                <div className="text-muted-foreground">Phone</div>
                <div className="mt-1 font-medium">{user.phone || '—'}</div>
              </div>
              <div className="rounded-lg border p-3">
                <div className="text-muted-foreground">Created</div>
                <div className="mt-1 font-medium">{formatDate(user.created_at, 'YYYY-MM-DD HH:mm')}</div>
              </div>
              <div className="rounded-lg border p-3">
                <div className="text-muted-foreground">Updated</div>
                <div className="mt-1 font-medium">{formatDate(user.updated_at, 'YYYY-MM-DD HH:mm')}</div>
              </div>
            </div>
            <div className="rounded-lg border p-3">
              <div className="text-muted-foreground">Bio</div>
              <div className="mt-1 font-medium">{user.bio || '—'}</div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Minimal User Info</CardTitle>
            <CardDescription>
              Response from <code>GET {userInfoPath}</code>.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4 text-sm">
            <div className="rounded-lg border p-3">
              <div className="text-muted-foreground">ID</div>
              <div className="mt-1 font-medium">{userInfo.id}</div>
            </div>
            <div className="rounded-lg border p-3">
              <div className="text-muted-foreground">Username</div>
              <div className="mt-1 font-medium">{userInfo.username}</div>
            </div>
            <div className="rounded-lg border p-3">
              <div className="text-muted-foreground">Nickname</div>
              <div className="mt-1 font-medium">{userInfo.nickname || '—'}</div>
            </div>
            <div className="rounded-lg border p-3">
              <div className="text-muted-foreground">Avatar</div>
              <div className="mt-1 break-all font-medium">{userInfo.avatar || '—'}</div>
            </div>
            <div className="rounded-lg border p-3">
              <div className="text-muted-foreground">Bio</div>
              <div className="mt-1 font-medium">{userInfo.bio || '—'}</div>
            </div>
            <p className="text-xs text-muted-foreground">
              {/* 当前后端里 info 接口仍返回较完整的用户结构，这里保留独立展示，方便后续和正式公共信息接口对齐。 */}
              The current backend implementation returns the same authenticated user payload for both routes, so this card
              is wired to the actual `info` endpoint and shows what it currently exposes.
            </p>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
