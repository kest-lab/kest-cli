'use client';

import Link from 'next/link';
import { useEffect, useMemo, useRef, useState } from 'react';
import { useRouter } from 'next/navigation';
import {
  Bell,
  LayoutPanelTop,
  LogOut,
} from 'lucide-react';
import { LanguageSwitcher } from '@/components/common';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { ROUTES, buildProjectInviteRoute } from '@/constants/routes';
import { useLogout } from '@/hooks/use-auth';
import { usePendingProjectInvitations } from '@/hooks/use-project-invitations';
import { useT } from '@/i18n/client';
import { useAuthStore } from '@/store/auth-store';
import { formatDate } from '@/utils';

const buildInitials = (name: string) =>
  name
    .split(/\s+/)
    .filter(Boolean)
    .slice(0, 2)
    .map((part) => part[0]?.toUpperCase() || '')
    .join('') || 'U';

export function ProjectTopbar() {
  const t = useT('project');
  const router = useRouter();
  const logout = useLogout();
  const user = useAuthStore.use.user();
  const isAuthenticated = useAuthStore.use.isAuthenticated();
  const [isNotificationOpen, setIsNotificationOpen] = useState(false);
  const closeNotificationTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  const displayName = user?.nickname || user?.username || user?.email || 'User';
  const initials = useMemo(() => buildInitials(displayName), [displayName]);
  const pendingInvitationsQuery = usePendingProjectInvitations(isAuthenticated);
  const pendingInvitations = pendingInvitationsQuery.data ?? [];
  const hasPendingInvitations = pendingInvitations.length > 0;

  useEffect(() => {
    return () => {
      if (closeNotificationTimerRef.current) {
        clearTimeout(closeNotificationTimerRef.current);
      }
    };
  }, []);

  const openNotifications = () => {
    if (closeNotificationTimerRef.current) {
      clearTimeout(closeNotificationTimerRef.current);
      closeNotificationTimerRef.current = null;
    }
    setIsNotificationOpen(true);
  };

  const scheduleCloseNotifications = () => {
    if (closeNotificationTimerRef.current) {
      clearTimeout(closeNotificationTimerRef.current);
    }
    closeNotificationTimerRef.current = setTimeout(() => {
      setIsNotificationOpen(false);
      closeNotificationTimerRef.current = null;
    }, 120);
  };

  const handleLogout = () => {
    logout();
    router.replace(ROUTES.AUTH.LOGIN);
  };

  return (
    <header className="z-40 flex h-16 shrink-0 items-center justify-between border-b border-border/60 bg-bg-surface/95 px-4 backdrop-blur md:px-6">
      <div className="flex min-w-0 items-center gap-4">
        <Link href={ROUTES.CONSOLE.PROJECTS} className="group flex items-center gap-3">
          <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-primary text-primary-foreground shadow-button-primary transition-transform group-hover:scale-105">
            <LayoutPanelTop className="h-4.5 w-4.5" />
          </div>
          <div className="min-w-0">
            <p className="truncate text-sm font-semibold tracking-tight text-text-main">KEST</p>
          </div>
        </Link>
      </div>

      <div className="flex items-center gap-2">
        <LanguageSwitcher />

        <Popover open={isNotificationOpen} onOpenChange={setIsNotificationOpen}>
          <PopoverTrigger asChild>
            <Button
              variant="ghost"
              isIcon
              className="relative h-9 w-9 rounded-full"
              onMouseEnter={openNotifications}
              onMouseLeave={scheduleCloseNotifications}
              onFocus={openNotifications}
              onBlur={scheduleCloseNotifications}
            >
              <Bell className="h-4 w-4 text-text-muted" />
              {hasPendingInvitations ? (
                <span className="absolute right-1.5 top-1.5 flex h-4 min-w-4 items-center justify-center rounded-full border-2 border-bg-surface bg-primary px-1 text-[10px] font-semibold leading-none text-primary-foreground">
                  {pendingInvitations.length > 9 ? '9+' : pendingInvitations.length}
                </span>
              ) : null}
              <span className="sr-only">{t('topbar.notifications')}</span>
            </Button>
          </PopoverTrigger>
          <PopoverContent
            align="end"
            className="w-[min(calc(100vw-2rem),22rem)] rounded-xl border bg-bg-surface p-0 shadow-premium"
            onMouseEnter={openNotifications}
            onMouseLeave={scheduleCloseNotifications}
            onOpenAutoFocus={event => event.preventDefault()}
          >
            <div className="border-b border-border/60 px-4 py-3">
              <div className="text-sm font-semibold text-text-main">
                {t('topbar.notifications')}
              </div>
              <div className="mt-0.5 text-xs text-text-muted">
                {hasPendingInvitations
                  ? t('topbar.pendingInvitationsCount', { count: pendingInvitations.length })
                  : t('topbar.noPendingInvitations')}
              </div>
            </div>
            <div className="max-h-96 overflow-y-auto p-2">
              {pendingInvitationsQuery.isLoading ? (
                <div className="space-y-2 p-2">
                  {Array.from({ length: 2 }).map((_, index) => (
                    <div key={index} className="h-16 animate-pulse rounded-lg bg-muted/50" />
                  ))}
                </div>
              ) : hasPendingInvitations ? (
                <div className="space-y-1">
                  {pendingInvitations.map(invitation => (
                    <Link
                      key={invitation.id}
                      href={buildProjectInviteRoute(invitation.slug)}
                      className="block rounded-lg px-3 py-2.5 text-left transition-colors hover:bg-primary/5"
                    >
                      <div className="flex items-start gap-3">
                        <Avatar className="h-8 w-8 border border-border/60">
                          <AvatarImage src={invitation.inviter_avatar} />
                          <AvatarFallback className="bg-primary/10 text-xs text-primary">
                            {buildInitials(invitation.inviter_name || invitation.inviter_email)}
                          </AvatarFallback>
                        </Avatar>
                        <div className="min-w-0 flex-1">
                          <div className="text-sm font-medium text-text-main">
                            {t('topbar.invitedByToProject', {
                              inviter: invitation.inviter_name,
                              project: invitation.project_name,
                            })}
                          </div>
                          <div className="mt-1 flex flex-wrap items-center gap-x-2 gap-y-1 text-xs text-text-muted">
                            <span>{t(`roles.${invitation.role}`)}</span>
                            <span>{formatDate(invitation.created_at, 'YYYY-MM-DD HH:mm')}</span>
                          </div>
                        </div>
                      </div>
                    </Link>
                  ))}
                </div>
              ) : (
                <div className="px-3 py-8 text-center text-sm text-text-muted">
                  {t('topbar.noPendingInvitations')}
                </div>
              )}
            </div>
          </PopoverContent>
        </Popover>

        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button
              variant="ghost"
              isIcon
              noScale
              className="h-9 w-9 overflow-hidden rounded-full border border-border/60"
            >
              <Avatar className="h-full w-full">
                <AvatarImage src={user?.avatar} />
              <AvatarFallback className="bg-primary/10 text-primary">
                  {initials}
                </AvatarFallback>
              </Avatar>
              <span className="sr-only">{t('topbar.profile')}</span>
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="w-56 rounded-xl p-1 shadow-premium">
            <div className="px-2 py-1.5 text-xs font-medium uppercase tracking-wider text-text-muted">
              {displayName}
            </div>
            <div className="my-1 h-px bg-border/60" />
            <DropdownMenuItem
              className="cursor-pointer rounded-lg text-destructive focus:bg-destructive/10"
              onClick={handleLogout}
            >
              <LogOut className="mr-2 h-4 w-4" />
              {t('topbar.logout')}
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </header>
  );
}
