'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { toast } from 'sonner';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { buildApiPath } from '@/config/api';
import { ROUTES } from '@/constants/routes';
import { useChangePassword, useDeleteAccount, useLogout, useProfile, useUpdateProfile } from '@/hooks/use-auth';
import { useT } from '@/i18n/client';
import { formatDate } from '@/utils';

const resolveStatusLabel = (status: number | string, active: string, inactive: string) => {
  if (status === 1 || status === 'active') {
    return active;
  }

  if (status === 0 || status === 'inactive') {
    return inactive;
  }

  return String(status);
};

export function AccountSettings() {
  const t = useT();
  const router = useRouter();
  const logout = useLogout();
  const { data: profile, isLoading } = useProfile();
  const updateProfileMutation = useUpdateProfile();
  const changePasswordMutation = useChangePassword();
  const deleteAccountMutation = useDeleteAccount();

  const [profileDraft, setProfileDraft] = useState<Partial<Record<'nickname' | 'avatar' | 'phone' | 'bio', string>>>({});
  const [passwordForm, setPasswordForm] = useState({
    oldPassword: '',
    newPassword: '',
    confirmPassword: '',
  });
  const profilePath = buildApiPath('/users/profile');
  const passwordPath = buildApiPath('/users/password');
  const accountPath = buildApiPath('/users/account');

  const getProfileValue = (field: 'nickname' | 'avatar' | 'phone' | 'bio') => {
    if (!profile) {
      return '';
    }

    // 用户未编辑时显示后端返回值；一旦本地改动，优先展示草稿内容。
    return profileDraft[field] ?? profile[field] ?? '';
  };

  const handleProfileSubmit = async (event: React.FormEvent) => {
    event.preventDefault();

    try {
      await updateProfileMutation.mutateAsync({
        nickname: getProfileValue('nickname').trim() || undefined,
        avatar: getProfileValue('avatar').trim() || undefined,
        phone: getProfileValue('phone').trim() || undefined,
        bio: getProfileValue('bio').trim() || undefined,
      });

      setProfileDraft({});
      toast.success(t.console('account.profileUpdated'));
    } catch {
      // Error toast is handled by the global HTTP error handler.
    }
  };

  const handlePasswordSubmit = async (event: React.FormEvent) => {
    event.preventDefault();

    if (passwordForm.newPassword !== passwordForm.confirmPassword) {
      toast.error(t.console('account.passwordMismatch'));
      return;
    }

    try {
      const result = await changePasswordMutation.mutateAsync({
        old_password: passwordForm.oldPassword,
        new_password: passwordForm.newPassword,
      });

      setPasswordForm({
        oldPassword: '',
        newPassword: '',
        confirmPassword: '',
      });

      toast.success(result.message || t.console('account.passwordChanged'));
    } catch {
      // Error toast is handled by the global HTTP error handler.
    }
  };

  const handleDeleteAccount = async () => {
    const confirmed = window.confirm(t.console('account.deleteConfirm'));
    if (!confirmed) {
      return;
    }

    try {
      await deleteAccountMutation.mutateAsync();
      // 删除成功后主动退出并回到登录页，避免保留已删除账号的上下文。
      logout();
      toast.success(t.console('account.accountDeleted'));
      router.replace(ROUTES.AUTH.LOGIN);
    } catch {
      // Error toast is handled by the global HTTP error handler.
    }
  };

  if (isLoading || !profile) {
    return (
      <div className="flex-1 space-y-4 p-6 pt-6">
        <div className="h-9 w-48 animate-pulse rounded bg-muted" />
        <div className="grid gap-4 xl:grid-cols-3">
          <div className="h-64 animate-pulse rounded-xl bg-muted xl:col-span-2" />
          <div className="h-64 animate-pulse rounded-xl bg-muted" />
        </div>
        <div className="grid gap-4 xl:grid-cols-2">
          <div className="h-72 animate-pulse rounded-xl bg-muted" />
          <div className="h-72 animate-pulse rounded-xl bg-muted" />
        </div>
      </div>
    );
  }

  return (
    <div className="flex-1 space-y-6 p-6 pt-6">
      <div className="space-y-1">
        <h1 className="text-3xl font-bold tracking-tight">{t.console('account.title')}</h1>
        <p className="text-sm text-muted-foreground">
          {t.console('account.connectedTo')} <code>{profilePath}</code>, <code>{passwordPath}</code>, and <code>{accountPath}</code>.
        </p>
      </div>

      <div className="grid gap-4 xl:grid-cols-3">
        <Card className="xl:col-span-2">
          <CardHeader>
            <CardTitle>{t.console('account.profile')}</CardTitle>
            <CardDescription>
              {t.console('account.profileReviewDescription')}
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleProfileSubmit} className="space-y-4">
              <div className="grid gap-4 md:grid-cols-2">
                <div className="space-y-2">
                  <Label htmlFor="username">{t.console('account.username')}</Label>
                  <Input id="username" value={profile.username} readOnly />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="email">{t.console('account.email')}</Label>
                  <Input id="email" value={profile.email} readOnly />
                </div>
              </div>

              <div className="grid gap-4 md:grid-cols-2">
                <div className="space-y-2">
                  <Label htmlFor="nickname">{t.console('account.nickname')}</Label>
                  <Input
                    id="nickname"
                    value={getProfileValue('nickname')}
                    onChange={(event) => setProfileDraft((current) => ({ ...current, nickname: event.target.value }))}
                    placeholder={t.console('account.displayNamePlaceholder')}
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="phone">{t.console('account.phone')}</Label>
                  <Input
                    id="phone"
                    value={getProfileValue('phone')}
                    onChange={(event) => setProfileDraft((current) => ({ ...current, phone: event.target.value }))}
                    placeholder="+353..."
                  />
                </div>
              </div>

              <div className="space-y-2">
                <Label htmlFor="avatar">{t.console('account.avatarUrl')}</Label>
                <Input
                  id="avatar"
                  value={getProfileValue('avatar')}
                  onChange={(event) => setProfileDraft((current) => ({ ...current, avatar: event.target.value }))}
                  placeholder="https://example.com/avatar.png"
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="bio">{t.console('account.bio')}</Label>
                <Textarea
                  id="bio"
                  value={getProfileValue('bio')}
                  onChange={(event) => setProfileDraft((current) => ({ ...current, bio: event.target.value }))}
                  placeholder={t.console('account.bioPlaceholder')}
                  rows={5}
                />
              </div>

              <Button type="submit" disabled={updateProfileMutation.isPending}>
                {t.console('account.saveProfile')}
              </Button>
            </form>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>{t.console('account.accountSnapshot')}</CardTitle>
            <CardDescription>
              {t.console('account.accountSnapshotDescription', { path: profilePath })}
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4 text-sm">
            <div className="flex items-center justify-between rounded-lg border p-3">
              <span className="text-muted-foreground">{t.console('account.status')}</span>
              <Badge variant="outline">{resolveStatusLabel(profile.status, t.common('active'), t.common('inactive'))}</Badge>
            </div>
            <div className="flex items-center justify-between rounded-lg border p-3">
              <span className="text-muted-foreground">{t.console('account.userId')}</span>
              <span className="font-medium">{profile.id}</span>
            </div>
            <div className="flex items-center justify-between rounded-lg border p-3">
              <span className="text-muted-foreground">{t.console('account.created')}</span>
              <span className="font-medium">{formatDate(profile.created_at, 'YYYY-MM-DD HH:mm')}</span>
            </div>
            <div className="flex items-center justify-between rounded-lg border p-3">
              <span className="text-muted-foreground">{t.console('account.updated')}</span>
              <span className="font-medium">{formatDate(profile.updated_at, 'YYYY-MM-DD HH:mm')}</span>
            </div>
            <div className="flex items-center justify-between rounded-lg border p-3">
              <span className="text-muted-foreground">{t.console('account.lastLogin')}</span>
              <span className="font-medium">
                {profile.last_login ? formatDate(profile.last_login, 'YYYY-MM-DD HH:mm') : t.console('account.never')}
              </span>
            </div>
            <div className="flex items-center justify-between rounded-lg border p-3">
              <span className="text-muted-foreground">{t.console('account.superAdmin')}</span>
              <span className="font-medium">{profile.is_super_admin ? t.console('account.yes') : t.console('account.no')}</span>
            </div>
          </CardContent>
        </Card>
      </div>

      <div className="grid gap-4 xl:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>{t.console('account.changePassword')}</CardTitle>
            <CardDescription>
              {t.console('account.changePasswordDescription', { path: passwordPath })}
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handlePasswordSubmit} className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="old-password">{t.console('account.currentPassword')}</Label>
                <Input
                  id="old-password"
                  type="password"
                  value={passwordForm.oldPassword}
                  onChange={(event) => setPasswordForm((current) => ({ ...current, oldPassword: event.target.value }))}
                  autoComplete="current-password"
                  required
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="new-password">{t.console('account.newPassword')}</Label>
                <Input
                  id="new-password"
                  type="password"
                  value={passwordForm.newPassword}
                  onChange={(event) => setPasswordForm((current) => ({ ...current, newPassword: event.target.value }))}
                  autoComplete="new-password"
                  required
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="confirm-password">{t.console('account.confirmNewPassword')}</Label>
                <Input
                  id="confirm-password"
                  type="password"
                  value={passwordForm.confirmPassword}
                  onChange={(event) => setPasswordForm((current) => ({ ...current, confirmPassword: event.target.value }))}
                  autoComplete="new-password"
                  required
                />
              </div>
              <Button type="submit" disabled={changePasswordMutation.isPending}>
                {t.console('account.updatePassword')}
              </Button>
            </form>
          </CardContent>
        </Card>

        <Card className="border-destructive/40">
          <CardHeader>
            <CardTitle>{t.console('account.dangerZone')}</CardTitle>
            <CardDescription>
              {t.console('account.dangerDescription', { path: accountPath })}
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-3 text-sm text-muted-foreground">
            <p>
              {t.console('account.dangerBody')}
            </p>
          </CardContent>
          <CardFooter>
            <Button
              variant="destructive"
              onClick={handleDeleteAccount}
              disabled={deleteAccountMutation.isPending}
            >
              {t.console('account.deleteAccount')}
            </Button>
          </CardFooter>
        </Card>
      </div>
    </div>
  );
}
