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

export function AccountSettings() {
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
      toast.success('Profile updated');
    } catch {
      // Error toast is handled by the global HTTP error handler.
    }
  };

  const handlePasswordSubmit = async (event: React.FormEvent) => {
    event.preventDefault();

    if (passwordForm.newPassword !== passwordForm.confirmPassword) {
      toast.error('New passwords do not match');
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

      toast.success(result.message || 'Password changed successfully');
    } catch {
      // Error toast is handled by the global HTTP error handler.
    }
  };

  const handleDeleteAccount = async () => {
    const confirmed = window.confirm('This will permanently delete your account. Continue?');
    if (!confirmed) {
      return;
    }

    try {
      await deleteAccountMutation.mutateAsync();
      // 删除成功后主动退出并回到登录页，避免保留已删除账号的上下文。
      logout();
      toast.success('Account deleted');
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
        <h1 className="text-3xl font-bold tracking-tight">Account Settings</h1>
        <p className="text-sm text-muted-foreground">
          Connected to <code>{profilePath}</code>, <code>{passwordPath}</code>, and <code>{accountPath}</code>.
        </p>
      </div>

      <div className="grid gap-4 xl:grid-cols-3">
        <Card className="xl:col-span-2">
          <CardHeader>
            <CardTitle>Profile</CardTitle>
            <CardDescription>
              Review the current account state and update editable profile fields.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleProfileSubmit} className="space-y-4">
              <div className="grid gap-4 md:grid-cols-2">
                <div className="space-y-2">
                  <Label htmlFor="username">Username</Label>
                  <Input id="username" value={profile.username} readOnly />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="email">Email</Label>
                  <Input id="email" value={profile.email} readOnly />
                </div>
              </div>

              <div className="grid gap-4 md:grid-cols-2">
                <div className="space-y-2">
                  <Label htmlFor="nickname">Nickname</Label>
                  <Input
                    id="nickname"
                    value={getProfileValue('nickname')}
                    onChange={(event) => setProfileDraft((current) => ({ ...current, nickname: event.target.value }))}
                    placeholder="Display name"
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="phone">Phone</Label>
                  <Input
                    id="phone"
                    value={getProfileValue('phone')}
                    onChange={(event) => setProfileDraft((current) => ({ ...current, phone: event.target.value }))}
                    placeholder="+353..."
                  />
                </div>
              </div>

              <div className="space-y-2">
                <Label htmlFor="avatar">Avatar URL</Label>
                <Input
                  id="avatar"
                  value={getProfileValue('avatar')}
                  onChange={(event) => setProfileDraft((current) => ({ ...current, avatar: event.target.value }))}
                  placeholder="https://example.com/avatar.png"
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="bio">Bio</Label>
                <Textarea
                  id="bio"
                  value={getProfileValue('bio')}
                  onChange={(event) => setProfileDraft((current) => ({ ...current, bio: event.target.value }))}
                  placeholder="Tell your team what this account is for."
                  rows={5}
                />
              </div>

              <Button type="submit" disabled={updateProfileMutation.isPending}>
                Save profile
              </Button>
            </form>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Account Snapshot</CardTitle>
            <CardDescription>
              Current values returned by <code>GET {profilePath}</code>.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4 text-sm">
            <div className="flex items-center justify-between rounded-lg border p-3">
              <span className="text-muted-foreground">Status</span>
              <Badge variant="outline">{resolveStatusLabel(profile.status)}</Badge>
            </div>
            <div className="flex items-center justify-between rounded-lg border p-3">
              <span className="text-muted-foreground">User ID</span>
              <span className="font-medium">{profile.id}</span>
            </div>
            <div className="flex items-center justify-between rounded-lg border p-3">
              <span className="text-muted-foreground">Created</span>
              <span className="font-medium">{formatDate(profile.created_at, 'YYYY-MM-DD HH:mm')}</span>
            </div>
            <div className="flex items-center justify-between rounded-lg border p-3">
              <span className="text-muted-foreground">Updated</span>
              <span className="font-medium">{formatDate(profile.updated_at, 'YYYY-MM-DD HH:mm')}</span>
            </div>
            <div className="flex items-center justify-between rounded-lg border p-3">
              <span className="text-muted-foreground">Last login</span>
              <span className="font-medium">
                {profile.last_login ? formatDate(profile.last_login, 'YYYY-MM-DD HH:mm') : 'Never'}
              </span>
            </div>
            <div className="flex items-center justify-between rounded-lg border p-3">
              <span className="text-muted-foreground">Super admin</span>
              <span className="font-medium">{profile.is_super_admin ? 'Yes' : 'No'}</span>
            </div>
          </CardContent>
        </Card>
      </div>

      <div className="grid gap-4 xl:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Change Password</CardTitle>
            <CardDescription>
              Calls <code>PUT {passwordPath}</code> with the current and new password.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handlePasswordSubmit} className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="old-password">Current password</Label>
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
                <Label htmlFor="new-password">New password</Label>
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
                <Label htmlFor="confirm-password">Confirm new password</Label>
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
                Update password
              </Button>
            </form>
          </CardContent>
        </Card>

        <Card className="border-destructive/40">
          <CardHeader>
            <CardTitle>Danger Zone</CardTitle>
            <CardDescription>
              Deletes the current account through <code>DELETE {accountPath}</code>.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-3 text-sm text-muted-foreground">
            <p>
              This action is permanent. The current authenticated user will be removed from the backend database.
            </p>
          </CardContent>
          <CardFooter>
            <Button
              variant="destructive"
              onClick={handleDeleteAccount}
              disabled={deleteAccountMutation.isPending}
            >
              Delete account
            </Button>
          </CardFooter>
        </Card>
      </div>
    </div>
  );
}
